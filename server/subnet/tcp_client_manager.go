/**
 * \file GBTCPClientManager.go
 * \version
 * \author wzy
 * \date  2018年01月15日 18:22:43
 * \brief client连接管理器
 *
 */

package subnet

import (
	"base"
	"base/def"
	"base/functime"
	"base/logger"
	"base/tcpconn"
	"base/util"
	// "errors"
	// "fmt"
	"io"
	"net"
	// "runtime"
	"servercomm"
	"sync"
	"time"
)

// websocket连接管理器
type GBTCPClientManager struct {
	connPool            tcpconn.ServerConnPool
	superexitchan       chan bool
	serverexitchanmutex sync.Mutex
	serverexitchan      map[uint32]chan bool
	serverhandler       IServerHandler

	// runningMsgQueue  int
	runningMsgChan   []chan *ConnectMsgQueueStruct
	maxRunningMsgNum int

	// lastchan         int
	lastwarningtime1 uint32
	lastwarningtime2 uint32

	connectMutex sync.Mutex
}

// 初始化消息处理线程
func (this *GBTCPClientManager) InitMsgQueue(runsum int) {
	// 最大同时处理的消息数量
	this.maxRunningMsgNum = runsum // 消息队列线程数
	if this.maxRunningMsgNum < 1 {
		this.maxRunningMsgNum = 1
		logger.Error("[GBTCPClientManager.InitMsgQueue] " +
			"消息处理线程数量过小，置为1...")
	}
	this.runningMsgChan = make([]chan *ConnectMsgQueueStruct,
		this.maxRunningMsgNum)
	logger.Debug("[GBTCPClientManager.InitMsgQueue] "+
		"Client 消息处理线程数量 ThreadNum[%d]", this.maxRunningMsgNum)
	i := 0
	for i < this.maxRunningMsgNum {
		this.runningMsgChan[i] = make(chan *ConnectMsgQueueStruct,
			500000)
		go this.RecvmsgProcess(i)
		i++
	}
}

// 并行处理接收消息队列数据
func (this *GBTCPClientManager) MultiRecvmsgQueue(
	index int) (normalreturn bool) {
	if this.runningMsgChan == nil ||
		this.runningMsgChan[index] == nil {
		return true
	}
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := util.GetPanicInfo(recover()); err != nil {
			logger.Error("[GBTCPClientManager.MultiRecvmsgQueue] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
			normalreturn = false
		}
	}()
	normalreturn = true

	msgchan := this.runningMsgChan[index]
	for msgqueues := range msgchan {
		curtime := uint32(time.Now().Unix())
		if curtime > msgqueues.msg.TimeStamp+1 &&
			curtime > this.lastwarningtime2 {
			logger.Error("[GBTCPClientManager.MultiRecvmsgQueue] "+
				"[消耗时间统计]服务器TCPClient处理消息延迟很严重"+
				" TimeIns[%d] Thread[%d]",
				curtime-msgqueues.msg.TimeStamp, index)
			this.lastwarningtime2 = curtime
		}
		this.serverhandler.TCPMsgParse(msgqueues.task,
			msgqueues.msg)
	}
	return
}

// 并行处理接收消息队列数据 守护协程
func (this *GBTCPClientManager) RecvmsgProcess(index int) {
	for {
		if this.MultiRecvmsgQueue(index) {
			// 正常退出
			break
		}
	}
}

// 分配消息处理线程
func (this *GBTCPClientManager) MultiQueueControl(
	msgqueues *ConnectMsgQueueStruct) {
	if this.maxRunningMsgNum < 1 {
		this.serverhandler.TCPMsgParse(msgqueues.task, msgqueues.msg)
		return
	}
	who := this.serverhandler.GetTCPMsgParseChan(msgqueues.task,
		this.maxRunningMsgNum, msgqueues.msg)
	this.runningMsgChan[who] <- msgqueues
}

// 添加一个Client
func (this *GBTCPClientManager) AddTCPClient(
	conn net.Conn, serverid uint64) *tcpconn.ServerConn {
	tcptask := this.connPool.NewServerConn(conn, serverid)

	tcptask.SetAlive(true)

	logger.Debug("[GBTCPClientManager.AddTCPClient] "+
		"AddTCPClient ServerID[%d] 当前连接数量 NowSum[%d]",
		serverid, this.connPool.Len())

	return tcptask
}

// 根据tempid 获取一个client
func (this *GBTCPClientManager) GetTCPClient(
	tempid uint64) *tcpconn.ServerConn {
	return this.connPool.Get(tempid)
}

// 根据tempid移除一个client
func (this *GBTCPClientManager) RemoveTCPClient(tempid uint64) {
	this.connPool.Remove(tempid)
}

// 修改链接的 tempip
func (this *GBTCPClientManager) ChangeTCPClientTempid(
	tcptask *tcpconn.ServerConn,
	newtempid uint64) {
	this.connPool.ChangeTempid(tcptask, newtempid)
}

// 广播消息
func (this *GBTCPClientManager) BroadcastCmd(v base.MsgStruct) {
	this.connPool.BroadcastCmd(v)
}

func (this *GBTCPClientManager) BroadcastByType(servertype uint32,
	v base.MsgStruct) {
	this.connPool.BroadcastByType(servertype, v)
}

func (this *GBTCPClientManager) GetMinUserServerClient() *tcpconn.
	ServerConn {
	return this.getMinServerClient(def.TypeUserServer)
}

func (this *GBTCPClientManager) GetMinMatchServerClient() *tcpconn.
	ServerConn {
	return this.getMinServerClient(def.TypeMatchServer)
}

func (this *GBTCPClientManager) GetMinRoomServerClient() *tcpconn.
	ServerConn {
	return this.getMinServerClient(def.TypeRoomServer)
}

func (this *GBTCPClientManager) GetMinRoomServerClientLatestVersion() *tcpconn.
	ServerConn {
	return this.getMinServerClientLatestVersion(def.TypeRoomServer)
}

func (this *GBTCPClientManager) getMinServerClient(
	servertype uint32) *tcpconn.ServerConn {
	return this.connPool.GetMinClient(servertype)
}

func (this *GBTCPClientManager) getMinServerClientLatestVersion(
	servertype uint32) *tcpconn.ServerConn {
	return this.connPool.GetMinClientLatestVersion(servertype)
}

//随机获取 userserver 的连接client gateway那边使用
func (this *GBTCPClientManager) RandomGetUserServerClient() *tcpconn.
	ServerConn {
	return this.RandomGetServerClient(def.TypeUserServer)
}

//随机获取 matchserver 的连接
func (this *GBTCPClientManager) RandomGetMatchServerClient() *tcpconn.
	ServerConn {
	return this.RandomGetServerClient(def.TypeMatchServer)
}

//随机获取 roomserver 的连接
func (this *GBTCPClientManager) RandomGetRoomServerClient() *tcpconn.
	ServerConn {
	return this.RandomGetServerClient(def.TypeRoomServer)
}

func (this *GBTCPClientManager) RandomGetServerClient(
	servertype uint32) *tcpconn.ServerConn {
	return this.connPool.GetRandom(servertype)
}

// 发送数据到super
func (this *GBTCPClientManager) SendCmdToSuper(v base.MsgStruct) {
	this.connPool.BroadcastByType(def.TypeSuperServer, v)
}

// 发送数据到data
func (this *GBTCPClientManager) SendCmdToDataServer(v base.MsgStruct) {
	this.connPool.BroadcastByType(def.TypeDataServer, v)
}

// 发送数据到bridge
func (this *GBTCPClientManager) SendCmdToRandBridgeServer(
	v base.MsgStruct) {
	client := this.RandomGetServerClient(def.TypeBridgeServer)
	if client == nil {
		logger.Error("[GBTCPClientManager.SendCmdToRandBridgeServer] "+
			"找不到服务器 MsgID[%d] MsgName[%s]", v.GetMsgId(), v.GetMsgName())
		return
	}
	client.SendCmd(v)
}

func handleClientConnection(client *tcpconn.ServerConn,
	clienter IServerHandler) {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := util.GetPanicInfo(recover()); err != nil {
			logger.Error("[handleClientConnection] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
		}
	}()
	// 消息缓冲
	netbuffer := base.NewIOBuffer(client.Conn.Conn, 6400*1024)
	msgReader := base.NewMessageBinaryReader(netbuffer)

	for !base.GetGBServerConfigM().TerminateServer {
		derr := client.Conn.Conn.SetReadDeadline(time.Now().
			Add(time.Duration(time.Millisecond * 1000)))
		if derr != nil {
			logger.Error("[handleClientConnection] SetReadDeadline Err[%s]"+
				"ServerName[%s] ServerID[%d]",
				derr.Error(), client.Serverinfo.Servername,
				client.Serverinfo.Serverid)
			onClientDisconnected(client, clienter)
			return
		}
		n, err := netbuffer.ReadFromReader()
		if err != nil {
			if err == io.EOF {
				logger.Debug("[GBTCPClientManager.handleClientConnection] "+
					"数据读写异常，断开连接"+
					" ReadLen[%d] ServerID[%d] Error[%s]",
					n, client.Tempid, err.Error())
				onClientDisconnected(client, clienter)
				return
			} else {
				continue
			}
		}
		functiontime := functime.FunctionTime{}
		functiontime.Start("handleClientConnection")

		err = msgReader.RangeMsgBinary(func(msgbinary *base.MessageBinary) {
			// 判断消息是否阻塞严重
			curtime := uint32(time.Now().Unix())
			if curtime > msgbinary.TimeStamp+1 &&
				curtime > GetGBTCPClientManager().lastwarningtime1 {
				logger.Error("[GBTCPClientManager.handleClientConnection] "+
					"[消耗时间统计] 服务器TCPClient"+
					"接收消息延迟很严重 TimeInc[%d]",
					curtime-msgbinary.TimeStamp)
				GetGBTCPClientManager().lastwarningtime1 = curtime
			}
			// 重置消息标签为接收时间，用于处理阻塞判断
			msgbinary.TimeStamp = curtime
			// 解析消息
			msgParseTCPClient(client, msgbinary, clienter)
		})
		functiontime.Stop()
		if err != nil {
			logger.Error("[GBTCPClientManager.handleClientConnection] "+
				"RangeMsgBinary读消息失败，断开连接 Err[%s]", err.Error())
			onClientDisconnected(client, clienter)
			return
		}
	}
}

func msgParseTCPClient(client *tcpconn.ServerConn,
	msgbin *base.MessageBinary, clienter IServerHandler) {
	switch msgbin.CmdID {
	case servercomm.STestCommandID:
		performance_test := base.GetGBServerConfigM().
			GetProp("performance_test")
		if performance_test == "true" {
			recvmsg := &servercomm.STestCommand{}
			recvmsg.ReadBinary([]byte(msgbin.ProtoData))
			logger.Debug("[GBTCPClientManager.msgParseTCPClient] "+
				"Server 收到测试消息 CmdLen[%d] No.[%d]",
				msgbin.CmdLen, recvmsg.Testno)
		}
		return
	case servercomm.STimeTickCommandID:
		recvmsg := &servercomm.STimeTickCommand{}
		recvmsg.ReadBinary([]byte(msgbin.ProtoData))
		return
	case servercomm.SLoginRetCommandID:
		// 收到登陆服务器返回的消息
		recvmsg := &servercomm.SLoginRetCommand{}
		recvmsg.ReadBinary([]byte(msgbin.ProtoData))
		if recvmsg.Loginfailed > 0 {
			GetGBTCPClientManager().RemoveTCPClient(client.Tempid)
			logger.Error("[GBTCPClientManager.msgParseTCPClient] " +
				"连接验证失败,断开连接")
			return
		}
		client.Serverinfo = recvmsg.Taskinfo
		logger.Debug("[GBTCPClientManager.msgParseTCPClient] "+
			"连接服务器验证成功,id:%d,ip:%s,port:%d",
			client.Serverinfo.Serverid, client.Serverinfo.Serverip,
			client.Serverinfo.Serverport)

		if client.Serverinfo.Servertype == def.TypeSuperServer &&
			recvmsg.Taskinfo.Servertype == def.TypeSuperServer &&
			base.GetGBServerConfigM().Myserverinfo.Servertype != def.TypeSuperServer {
			// 如果登陆 SuperServer 成功，且当前不是 SuperServer
			base.GetGBServerConfigM().Myserverinfo = recvmsg.Clientinfo
			logger.Debug("[GBTCPClientManager.msgParseTCPClient] "+
				"获取本机信息成功"+
				" ServerID[%d] IP[%s] Port[%d] HTTPPort[%d] Name[%s]",
				recvmsg.Clientinfo.Serverid,
				recvmsg.Clientinfo.Serverip,
				recvmsg.Clientinfo.Serverport,
				recvmsg.Clientinfo.Httpport,
				recvmsg.Clientinfo.Servername)
			base.GetGBServerConfigM().RedisConfig = recvmsg.Redisinfo
		}
		return
	case servercomm.SLogoutCommandID:
		// 服务器已主动关闭，不再尝试连接它了
		logger.Debug("[msgParseTCPClient] 服务器已主动关闭，不再尝试连接它了 "+
			"ServerInfo[%s]", client.Serverinfo.GetJson())
		GetGBTCPClientManager().serverexitchanmutex.Lock()
		defer GetGBTCPClientManager().serverexitchanmutex.Unlock()
		GetSubnetManager().ServerConfigs.
			RemoveServerConfig(client.Serverinfo.Serverid)
		return
	case servercomm.SStartMyNotifyCommandID:
		// 收到依赖我的服务器的信息
		recvmsg := &servercomm.SStartMyNotifyCommand{}
		recvmsg.ReadBinary([]byte(msgbin.ProtoData))
		serverinfo := recvmsg.Serverinfo
		logger.Debug("[GBTCPClientManager.msgParseTCPClient] "+
			"收到服务器启动信息"+
			" ServerID[%d] IP[%s] Port[%d] HTTPPort[%d] Name[%s]",
			serverinfo.Serverid, serverinfo.Serverip,
			serverinfo.Serverport, serverinfo.Httpport,
			serverinfo.Servername)
		GetGBTCPClientManager().serverexitchanmutex.Lock()
		defer GetGBTCPClientManager().serverexitchanmutex.Unlock()
		GetSubnetManager().ServerConfigs.AddServerConfig(serverinfo)
		return
	case servercomm.SNotifyAllInfoID:
		// 收到所有服务器的配置信息
		recvmsg := &servercomm.SNotifyAllInfo{}
		recvmsg.ReadBinary([]byte(msgbin.ProtoData))
		GetGBTCPClientManager().serverexitchanmutex.Lock()
		defer GetGBTCPClientManager().serverexitchanmutex.Unlock()
		if client.Serverinfo.Servertype == def.TypeSuperServer {
			GetSubnetManager().ServerConfigs.CleanServerConfig()
		}
		logger.Debug("[GBTCPClientManager.msgParseTCPClient] " +
			"收到所有服务器列表信息")
		// 所有服务器信息列表
		for i := 0; i < len(recvmsg.Serverinfos); i++ {
			serverinfo := recvmsg.Serverinfos[i]
			GetSubnetManager().ServerConfigs.AddServerConfig(serverinfo)
		}
		return
	case servercomm.SStartRelyNotifyCommandID:
		if client.Serverinfo.Servertype == def.TypeSuperServer {
			// 收到了我依赖的服务器的信息
			recvmsg := &servercomm.SStartRelyNotifyCommand{}
			recvmsg.ReadBinary([]byte(msgbin.ProtoData))
			connectRelyServers(recvmsg.Serverinfos, clienter)
		}
		return
	}
	msgqueues := &ConnectMsgQueueStruct{}
	msgqueues.task = client
	msgqueues.msg = msgbin
	GetGBTCPClientManager().MultiQueueControl(msgqueues)
}
