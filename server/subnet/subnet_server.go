package subnet

import (
	"fmt"
	"github.com/liasece/micserver/connect"
	"github.com/liasece/micserver/servercomm"
	"github.com/liasece/micserver/util/sysutil"
	"github.com/liasece/micserver/util/uid"
	"net"
)

func (this *SubnetManager) OnServerLogin(conn *connect.Server,
	tarinfo *servercomm.SLoginCommand) {
	this.connectMutex.Lock()
	defer this.connectMutex.Unlock()

	// 来源服务器请求登陆本服务器
	myconn := this.GetServer(fmt.Sprint(tarinfo.ServerID))
	if myconn != nil {
		this.Debug("[SubnetManager.OnServerLogin] 重复连接 %s 优先级：%d:%d",
			tarinfo.ServerID,
			myconn.ConnectPriority, tarinfo.ConnectPriority)
		if myconn.ConnectPriority < tarinfo.ConnectPriority {
			// 对方连接的优先级比较高，删除我方连接
			myconn.IsNormalDisconnect = true
			myconn.Terminate()
			unuseid, _ := uid.NewUniqueID(0xff)
			this.ChangeServerTempid(
				myconn, myconn.Tempid+"unuse"+fmt.Sprint(unuseid))
		} else {
			// 我方优先级比较高已经连接成功过了，非法连接
			this.Debug("[SubNetManager.OnServerLogin] "+
				"收到了重复的Server连接请求 Msg[%s]",
				tarinfo.GetJson())
			conn.IsNormalDisconnect = true
			retmsg := &servercomm.SLoginRetCommand{}
			retmsg.Loginfailed = servercomm.LOGINRETCODE_IDENTICAL
			conn.SendCmd(retmsg)
			conn.Terminate()
			return
		}
	}
	serverInfo := &servercomm.ServerInfo{}

	// 来源服务器地址
	remoteaddr := conn.RemoteAddr()
	// 获取来源服务器ID在本地的配置
	serverInfo.ServerID = tarinfo.ServerID
	serverInfo.ServerAddr = tarinfo.ServerAddr
	// 检查是否获取信息成功
	if serverInfo.ServerID == "" {
		// 如果获取信息不成功
		this.Error("[SubNetManager.OnServerLogin] "+
			"连接分配异常 未知服务器连接 "+
			"Addr[%s] Msg[%s]",
			remoteaddr, tarinfo.GetJson())
		retmsg := &servercomm.SLoginRetCommand{}
		retmsg.Loginfailed = servercomm.LOGINRETCODE_IDENTITY
		conn.SendCmd(retmsg)
		conn.Terminate()
		return
	}

	serverInfo.Version = tarinfo.Version

	// 来源服务器检查完毕
	// 完善来源服务器在本服务器的信息
	this.ChangeServerTempid(conn, fmt.Sprint(serverInfo.ServerID))
	conn.ServerInfo = serverInfo
	conn.SetVertify(true)
	conn.SetTerminateTime(0) // 清除终止时间状态
	this.Debug("[SubNetManager.OnServerLogin] "+
		"客户端连接验证成功 "+
		" SerID[%s] IP[%s]",
		serverInfo.ServerID, serverInfo.ServerAddr)
	// 向来源服务器回复登陆成功消息
	retmsg := &servercomm.SLoginRetCommand{}
	retmsg.Loginfailed = 0
	retmsg.Destination = this.myServerInfo
	conn.SendCmd(retmsg)

	// 通知其他服务器，次服务器登陆完成
	// 如果我是SuperServer
	// 向来源服务器发送本地已有的所有服务器信息
	this.NotifyAllServerInfo(conn)
	// 把来源服务器信息广播给其它所有服务器
	notifymsg := &servercomm.SStartMyNotifyCommand{}
	notifymsg.ServerInfo = serverInfo
	this.BroadcastCmd(notifymsg)
}

func (this *SubnetManager) BindTCPSubnet(settings map[string]string) error {
	addr, hasconf := settings["subnettcpaddr"]
	if !hasconf {
		return fmt.Errorf("subnettcpaddr hasn't set.")
	}
	// init tcp subnet port
	netlisten, err := net.Listen("tcp", addr)
	if err != nil {
		this.Error("[SubNetManager.BindTCPServer] "+
			"服务器绑定失败 IPPort[%s] Err[%s]",
			addr, err.Error())
		return err
	}
	this.Debug("[SubNetManager.BindTCPServer] "+
		"服务器绑定成功 IPPort[%s]", addr)
	this.myServerInfo.ServerAddr = addr
	go this.TCPServerListenerProcess(netlisten)
	return nil
}

func (this *SubnetManager) TCPServerListenerProcess(listener net.Listener) {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := sysutil.GetPanicInfo(recover()); err != nil {
			this.Error("[TCPServerListenerProcess] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
		}
	}()
	defer listener.Close()
	for true {
		this.mTCPServerListener(listener)
	}
}

func (this *SubnetManager) mTCPServerListener(listener net.Listener) {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := sysutil.GetPanicInfo(recover()); err != nil {
			// 这里的err其实就是panic传入的内容
			this.Error("[mTCPServerListener] "+
				"Panic: ErrName[%v] \n Stack[%s]", err, stackInfo)
		}
	}()

	for true {
		newconn, err := listener.Accept()
		if err != nil {
			this.Error("[SubNetManager.TCPServerListenerProcess] "+
				"服务器端口监听异常 Err[%s]",
				err.Error())
			continue
		}
		this.Debug("[SubNetManager.BindTCPServer] "+
			"收到新的TCP连接 Addr[%s]",
			newconn.RemoteAddr().String())
		conn := this.NewServer(connect.ServerSCTypeTask, newconn, "",
			this.onConnectRecv, this.onConnectClose)
		if conn != nil {
			conn.Logger = this.Logger
			this.OnCreateTCPConnect(conn)
		}
	}
}
