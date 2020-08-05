package subnet

import (
	"fmt"
	"net"

	"github.com/liasece/micserver/conf"
	"github.com/liasece/micserver/connect"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/process"
	"github.com/liasece/micserver/servercomm"
	"github.com/liasece/micserver/util/sysutil"
	"github.com/liasece/micserver/util/uid"
)

// OnServerLogin 目标服务器尝试登陆到本服务器
func (manager *Manager) OnServerLogin(conn *connect.Server, tarinfo *servercomm.SLoginCommand) {
	manager.connectMutex.Lock()
	defer manager.connectMutex.Unlock()

	manager.Syslog("收到登陆请求 Server:%s", tarinfo.ModuleID)

	// 来源服务器请求登陆本服务器
	myconn := manager.GetServer(fmt.Sprint(tarinfo.ModuleID))
	if myconn != nil {
		manager.Syslog("[Manager.OnServerLogin] Duplicate connections", log.String("ModuleID", tarinfo.ModuleID),
			log.Int64("LocalPriority", myconn.ConnectPriority), log.Int64("IncomePriority", tarinfo.ConnectPriority))
		if myconn.ConnectPriority < tarinfo.ConnectPriority {
			// 对方连接的优先级比较高，删除我方连接
			myconn.IsNormalDisconnect = true
			myconn.Terminate()
			unuseid, _ := uid.GenUniqueID(0)
			manager.ChangeServerTempid(
				myconn, myconn.GetTempID()+"unuse"+unuseid)
		} else {
			// 我方优先级比较高已经连接成功过了，非法连接
			manager.Syslog("[Manager.OnServerLogin] Duplicate Server connection requests received", log.Reflect("Msg", tarinfo))
			return
		}
	}
	serverInfo := &servercomm.ModuleInfo{}

	// 来源服务器地址
	remoteaddr := conn.RemoteAddr()
	// 获取来源服务器ID在本地的配置
	serverInfo.ModuleID = tarinfo.ModuleID
	serverInfo.ModuleAddr = tarinfo.ModuleAddr
	// 检查是否获取信息成功
	if serverInfo.ModuleID == "" {
		// 如果获取信息不成功
		manager.Error("[Manager.OnServerLogin] Connection allocation abnormal, unknown server connection",
			log.String("Addr", remoteaddr), log.Reflect("Msg", tarinfo))
		retmsg := &servercomm.SLoginRetCommand{}
		retmsg.Loginfailed = servercomm.LoginRetCodeIdentity
		conn.SendCmd(retmsg)
		conn.Terminate()
		return
	}
	manager.Syslog("[Manager.OnServerLogin] Successful client connection verification",
		log.String("ModuleID", serverInfo.ModuleID), log.String("Addr", serverInfo.ModuleAddr))

	serverInfo.Version = tarinfo.Version

	// 来源服务器检查完毕
	// 完善来源服务器在本服务器的信息
	manager.ChangeServerTempid(conn, fmt.Sprint(serverInfo.ModuleID))
	conn.ModuleInfo = serverInfo
	conn.SetVertify(true)
	conn.SetTerminateTime(0) // 清除终止时间状态
	// 向来源服务器回复登陆成功消息
	retmsg := &servercomm.SLoginRetCommand{}
	retmsg.Loginfailed = 0
	retmsg.Destination = manager.myServerInfo
	conn.SendCmd(retmsg)

	// 通知其他服务器，次服务器登陆完成
	// 如果我是SuperServer
	// 向来源服务器发送本地已有的所有服务器信息
	manager.NotifyAllServerInfo(conn)
	// 把来源服务器信息广播给其它所有服务器
	notifymsg := &servercomm.SStartMyNotifyCommand{}
	notifymsg.ModuleInfo = serverInfo
	manager.BroadcastCmd(notifymsg)
	manager.subnetHook.OnServerJoinSubnet(conn)
}

// BindTCPSubnet 绑定本服务器对子网开放的端口
func (manager *Manager) BindTCPSubnet(settings *conf.ModuleConfig) error {
	if !settings.Exist(conf.SubnetTCPAddr) {
		return fmt.Errorf("subnettcpaddr hasn't set")
	}
	addr := settings.GetString(conf.SubnetTCPAddr)
	// init tcp subnet port
	netlisten, err := net.Listen("tcp", addr)
	if err != nil {
		manager.Error("[Manager.BindTCPSubnet] Server binding failure", log.String("IPPort", addr), log.ErrorField(err))
		return err
	}
	manager.Syslog("[Manager.BindTCPSubnet] Server binding success", log.String("IPPort", addr))
	manager.myServerInfo.ModuleAddr = addr
	go manager.TCPServerListenerProcess(netlisten)
	return nil
}

// TCPServerListenerProcess 监听本服务器的子网端口线程
func (manager *Manager) TCPServerListenerProcess(listener net.Listener) {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if stackInfo, err := sysutil.GetPanicInfo(recover()); err != nil {
			manager.Error("[Manager.TCPServerListenerProcess] Panic", log.ErrorField(err), log.String("Stack", stackInfo))
		}
	}()
	defer listener.Close()
	for true {
		manager.mTCPServerListener(listener)
	}
}

// mTCPServerListener 保持监听本地服务器对子网端口
func (manager *Manager) mTCPServerListener(listener net.Listener) {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if stackInfo, err := sysutil.GetPanicInfo(recover()); err != nil {
			// 这里的err其实就是panic传入的内容
			manager.Error("[Manager.mTCPServerListener] Panic", log.ErrorField(err), log.String("Stack", stackInfo))
		}
	}()

	for true {
		newconn, err := listener.Accept()
		if err != nil {
			manager.Error("[Manager.mTCPServerListener] Server port listening error", log.ErrorField(err))
			continue
		}
		manager.Syslog("[Manager.mTCPServerListener] New TCP connection received", log.String("Addr", newconn.RemoteAddr().String()))
		conn := manager.NewTCPServer(connect.ServerSCTypeTask, newconn, "", manager.onConnectRecv, manager.onConnectClose)
		if conn != nil {
			conn.Logger = manager.Logger
			manager.OnCreateNewServer(conn)
		}
	}
}

// Local chan server init

// BindChanSubnet 绑定本地 chan 连接类型
func (manager *Manager) BindChanSubnet(settings *conf.ModuleConfig) error {
	nochan := settings.GetBool(conf.SubnetNoChan)
	if nochan {
		return nil
	}
	serverChan := make(chan *process.ChanServerHandshake, 1000)
	process.AddServerChan(manager.myServerInfo.ModuleID, serverChan)
	go manager.ChanServerListenerProcess(serverChan)
	manager.Syslog("[Manager.BindChanSubnet] ChanServer register successfully")
	return nil
}

// ChanServerListenerProcess 监听本地 chan 连接的消息
func (manager *Manager) ChanServerListenerProcess(serverChan chan *process.ChanServerHandshake) {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if stackInfo, err := sysutil.GetPanicInfo(recover()); err != nil {
			manager.Error("[ChanServerListenerProcess] Panic", log.ErrorField(err), log.String("Stack", stackInfo))
		}
	}()
	defer func() {
		process.DeleteServerChan(manager.myServerInfo.ModuleID)
		close(serverChan)
	}()
	for true {
		manager.mChanServerListener(serverChan)
	}
}

// mChanServerListener 监听本地 chan 连接握手请求
func (manager *Manager) mChanServerListener(serverChan chan *process.ChanServerHandshake) {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if stackInfo, err := sysutil.GetPanicInfo(recover()); err != nil {
			// 这里的err其实就是panic传入的内容
			manager.Error("[Manager.mChanServerListener] Panic", log.ErrorField(err), log.String("Stack", stackInfo))
		}
	}()

	for true {
		select {
		case newinfo, ok := <-serverChan:
			if !ok {
				break
			}
			manager.processChanServerRequest(newinfo)
		}
	}
}

// processChanServerRequest 处理 chan 握手请求的返回信息
func (manager *Manager) processChanServerRequest(newinfo *process.ChanServerHandshake) {
	remoteChan := process.GetServerChan(newinfo.ModuleInfo.ModuleID)
	if remoteChan != nil {
		if newinfo.Seq == 0 {
			manager.connectMutex.Lock()
			oldconn := manager.GetServer(newinfo.ModuleInfo.ModuleID)
			// 重复连接
			if oldconn != nil {
				manager.Syslog("[Manager.mChanServerListener] Duplicate Connections", log.String("ModuleID", newinfo.ModuleInfo.ModuleID))
			} else {
				// 请求开始
				newMsgChan := make(chan *msg.MessageBinary, 1000)
				remoteChan <- &process.ChanServerHandshake{
					ModuleInfo:    manager.myServerInfo,
					ServerMsgChan: newMsgChan,
					ClientMsgChan: newinfo.ClientMsgChan,
					Seq:           newinfo.Seq + 1,
				}
				manager.Syslog("[Manager.mChanServerListener] Receive a new ServerChan connection", log.String("ModuleID", newinfo.ModuleInfo.ModuleID))
				// 建立本地通信Server对象
				conn := manager.NewChanServer(connect.ServerSCTypeTask, newinfo.ClientMsgChan, newMsgChan, "", manager.onConnectRecv, manager.onConnectClose)
				conn.Logger = manager.Logger
				manager.OnCreateNewServer(conn)
			}
			manager.connectMutex.Unlock()
		} else if newinfo.Seq == 1 {
			manager.connectMutex.Lock()
			oldconn := manager.GetServer(newinfo.ModuleInfo.ModuleID)
			// 重复连接
			if oldconn != nil {
				manager.Syslog("[Manager.mChanServerListener] Duplicate Connections", log.String("ModuleID", newinfo.ModuleInfo.ModuleID))
			} else {
				// 请求回复
				manager.Syslog("[Manager.mChanServerListener] Receive a response to a ServerChan connection request", log.String("ModuleID", newinfo.ModuleInfo.ModuleID))
				// 建立本地通信Server对象
				manager.doConnectChanServer(newinfo.ServerMsgChan, newinfo.ClientMsgChan, newinfo.ModuleInfo.ModuleID)
			}
			manager.connectMutex.Unlock()
		}
	}
}
