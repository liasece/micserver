package subnet

import (
	"fmt"
	"net"

	"github.com/liasece/micserver/conf"
	"github.com/liasece/micserver/connect"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/process"
	"github.com/liasece/micserver/servercomm"
	"github.com/liasece/micserver/util/sysutil"
	"github.com/liasece/micserver/util/uid"
)

// OnServerLogin 目标服务器尝试登陆到本服务器
func (manager *Manager) OnServerLogin(conn *connect.Server,
	tarinfo *servercomm.SLoginCommand) {
	manager.connectMutex.Lock()
	defer manager.connectMutex.Unlock()

	manager.Syslog("收到登陆请求 Server:%s", tarinfo.ModuleID)

	// 来源服务器请求登陆本服务器
	myconn := manager.GetServer(fmt.Sprint(tarinfo.ModuleID))
	if myconn != nil {
		manager.Syslog("[Manager.OnServerLogin] 重复连接 %s 优先级：%d:%d", tarinfo.ModuleID, myconn.ConnectPriority, tarinfo.ConnectPriority)
		if myconn.ConnectPriority < tarinfo.ConnectPriority {
			// 对方连接的优先级比较高，删除我方连接
			myconn.IsNormalDisconnect = true
			myconn.Terminate()
			unuseid, _ := uid.GenUniqueID(0)
			manager.ChangeServerTempid(
				myconn, myconn.GetTempID()+"unuse"+unuseid)
		} else {
			// 我方优先级比较高已经连接成功过了，非法连接
			manager.Syslog("[SubNetManager.OnServerLogin] 收到了重复的Server连接请求 Msg[%s]", tarinfo.GetJSON())
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
		manager.Error("[SubNetManager.OnServerLogin] 连接分配异常 未知服务器连接 Addr[%s] Msg[%s]", remoteaddr, tarinfo.GetJSON())
		retmsg := &servercomm.SLoginRetCommand{}
		retmsg.Loginfailed = servercomm.LOGINRETCODE_IDENTITY
		conn.SendCmd(retmsg)
		conn.Terminate()
		return
	}
	manager.Syslog("[SubNetManager.OnServerLogin] 客户端连接验证成功 SerID[%s] IP[%s]", serverInfo.ModuleID, serverInfo.ModuleAddr)

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
		manager.Error("[SubNetManager.BindTCPSubnet] 服务器绑定失败 IPPort[%s] Err[%s]", addr, err.Error())
		return err
	}
	manager.Syslog("[SubNetManager.BindTCPSubnet] 服务器绑定成功 IPPort[%s]", addr)
	manager.myServerInfo.ModuleAddr = addr
	go manager.TCPServerListenerProcess(netlisten)
	return nil
}

// TCPServerListenerProcess 监听本服务器的子网端口线程
func (manager *Manager) TCPServerListenerProcess(listener net.Listener) {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if stackInfo, err := sysutil.GetPanicInfo(recover()); err != nil {
			manager.Error("[SubNetManager.TCPServerListenerProcess] Panic: Err[%v] \n Stack[%s]", err, stackInfo)
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
			manager.Error("[SubNetManager.mTCPServerListener] Panic: ErrName[%v] \n Stack[%s]", err, stackInfo)
		}
	}()

	for true {
		newconn, err := listener.Accept()
		if err != nil {
			manager.Error("[SubNetManager.mTCPServerListener] 服务器端口监听异常 Err[%s]", err.Error())
			continue
		}
		manager.Syslog("[SubNetManager.mTCPServerListener] 收到新的TCP连接 Addr[%s]", newconn.RemoteAddr().String())
		conn := manager.NewTCPServer(connect.ServerSCTypeTask, newconn, "",
			manager.onConnectRecv, manager.onConnectClose)
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
	manager.Syslog("BindChanSubnet ChanServer 注册成功")
	return nil
}

// ChanServerListenerProcess 监听本地 chan 连接的消息
func (manager *Manager) ChanServerListenerProcess(
	serverChan chan *process.ChanServerHandshake) {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if stackInfo, err := sysutil.GetPanicInfo(recover()); err != nil {
			manager.Error("[ChanServerListenerProcess] Panic: Err[%v] \n Stack[%s]", err, stackInfo)
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
func (manager *Manager) mChanServerListener(
	serverChan chan *process.ChanServerHandshake) {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if stackInfo, err := sysutil.GetPanicInfo(recover()); err != nil {
			// 这里的err其实就是panic传入的内容
			manager.Error("[SubNetManager.mChanServerListener] Panic: ErrName[%v] \n Stack[%s]", err, stackInfo)
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
func (manager *Manager) processChanServerRequest(
	newinfo *process.ChanServerHandshake) {
	remoteChan := process.GetServerChan(newinfo.ModuleInfo.ModuleID)
	if remoteChan != nil {
		if newinfo.Seq == 0 {
			manager.connectMutex.Lock()
			oldconn := manager.GetServer(newinfo.ModuleInfo.ModuleID)
			// 重复连接
			if oldconn != nil {
				manager.Syslog("[Manager.mChanServerListener] ModuleID[%s] 重复的连接", newinfo.ModuleInfo.ModuleID)
			} else {
				// 请求开始
				newMsgChan := make(chan *msg.MessageBinary, 1000)
				remoteChan <- &process.ChanServerHandshake{
					ModuleInfo:    manager.myServerInfo,
					ServerMsgChan: newMsgChan,
					ClientMsgChan: newinfo.ClientMsgChan,
					Seq:           newinfo.Seq + 1,
				}
				manager.Syslog("[SubNetManager.mChanServerListener] 收到新的 ServerChan 连接 ModuleID[%s]",
					newinfo.ModuleInfo.ModuleID)
				// 建立本地通信Server对象
				conn := manager.NewChanServer(connect.ServerSCTypeTask,
					newinfo.ClientMsgChan, newMsgChan, "",
					manager.onConnectRecv, manager.onConnectClose)
				conn.Logger = manager.Logger
				manager.OnCreateNewServer(conn)
			}
			manager.connectMutex.Unlock()
		} else if newinfo.Seq == 1 {
			manager.connectMutex.Lock()
			oldconn := manager.GetServer(newinfo.ModuleInfo.ModuleID)
			// 重复连接
			if oldconn != nil {
				manager.Syslog("[Manager.mChanServerListener] ModuleID[%s] 重复的连接", newinfo.ModuleInfo.ModuleID)
			} else {
				// 请求回复
				manager.Syslog("[SubNetManager.mChanServerListener] 收到 ServerChan 连接请求回复 ModuleID[%s]", newinfo.ModuleInfo.ModuleID)
				// 建立本地通信Server对象
				manager.doConnectChanServer(
					newinfo.ServerMsgChan, newinfo.ClientMsgChan,
					newinfo.ModuleInfo.ModuleID)
			}
			manager.connectMutex.Unlock()
		}
	}
}
