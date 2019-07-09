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
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/tcpconn"
	"github.com/liasece/micserver/util"
	// "errors"
	// "fmt"
	"io"
	"net"
	// "runtime"
	"github.com/liasece/micserver/comm"
	"github.com/liasece/micserver/msg"
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

	subnetManager *SubnetManager
}

// 初始化消息处理线程
func (this *GBTCPClientManager) InitMsgQueue(runsum int) {
	// 最大同时处理的消息数量
	this.maxRunningMsgNum = runsum // 消息队列线程数
	if this.maxRunningMsgNum < 1 {
		this.maxRunningMsgNum = 1
		log.Error("[GBTCPClientManager.InitMsgQueue] " +
			"消息处理线程数量过小，置为1...")
	}
	this.runningMsgChan = make([]chan *ConnectMsgQueueStruct,
		this.maxRunningMsgNum)
	log.Debug("[GBTCPClientManager.InitMsgQueue] "+
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
			log.Error("[GBTCPClientManager.MultiRecvmsgQueue] "+
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
			log.Error("[GBTCPClientManager.MultiRecvmsgQueue] "+
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

	log.Debug("[GBTCPClientManager.AddTCPClient] "+
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
func (this *GBTCPClientManager) BroadcastCmd(v msg.MsgStruct) {
	this.connPool.BroadcastCmd(v)
}

func (this *GBTCPClientManager) BroadcastByType(servertype uint32,
	v msg.MsgStruct) {
	this.connPool.BroadcastByType(servertype, v)
}

func (this *GBTCPClientManager) RandomGetServerClient(
	servertype uint32) *tcpconn.ServerConn {
	return this.connPool.GetRandom(servertype)
}

func (this *GBTCPClientManager) handleClientConnection(client *tcpconn.ServerConn,
	clienter IServerHandler) {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := util.GetPanicInfo(recover()); err != nil {
			log.Error("[handleClientConnection] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
		}
	}()
	// 消息缓冲
	netbuffer := util.NewIOBuffer(client.Conn.Conn, 6400*1024)
	msgReader := msg.NewMessageBinaryReader(netbuffer)

	for true {
		derr := client.Conn.Conn.SetReadDeadline(time.Now().
			Add(time.Duration(time.Millisecond * 1000)))
		if derr != nil {
			log.Error("[handleClientConnection] SetReadDeadline Err[%s]"+
				"ServerName[%s] ServerID[%d]",
				derr.Error(), client.Serverinfo.Servername,
				client.Serverinfo.Serverid)
			this.onClientDisconnected(client, clienter)
			return
		}
		n, err := netbuffer.ReadFromReader()
		if err != nil {
			if err == io.EOF {
				log.Debug("[GBTCPClientManager.handleClientConnection] "+
					"数据读写异常，断开连接"+
					" ReadLen[%d] ServerID[%d] Error[%s]",
					n, client.Tempid, err.Error())
				this.onClientDisconnected(client, clienter)
				return
			} else {
				continue
			}
		}
		functiontime := util.FunctionTime{}
		functiontime.Start("handleClientConnection")

		err = msgReader.RangeMsgBinary(func(msgbinary *msg.MessageBinary) {
			// 判断消息是否阻塞严重
			curtime := uint32(time.Now().Unix())
			if curtime > msgbinary.TimeStamp+1 &&
				curtime > this.lastwarningtime1 {
				log.Error("[GBTCPClientManager.handleClientConnection] "+
					"[消耗时间统计] 服务器TCPClient"+
					"接收消息延迟很严重 TimeInc[%d]",
					curtime-msgbinary.TimeStamp)
				this.lastwarningtime1 = curtime
			}
			// 重置消息标签为接收时间，用于处理阻塞判断
			msgbinary.TimeStamp = curtime
			// 解析消息
			this.msgParseTCPClient(client, msgbinary, clienter)
		})
		functiontime.Stop()
		if err != nil {
			log.Error("[GBTCPClientManager.handleClientConnection] "+
				"RangeMsgBinary读消息失败，断开连接 Err[%s]", err.Error())
			this.onClientDisconnected(client, clienter)
			return
		}
	}
}

func (this *GBTCPClientManager) msgParseTCPClient(client *tcpconn.ServerConn,
	msgbin *msg.MessageBinary, clienter IServerHandler) {
	switch msgbin.CmdID {
	case comm.STestCommandID:
		performance_test := this.subnetManager.moudleConf.
			GetProp("performance_test")
		if performance_test == "true" {
			recvmsg := &comm.STestCommand{}
			recvmsg.ReadBinary([]byte(msgbin.ProtoData))
			log.Debug("[GBTCPClientManager.msgParseTCPClient] "+
				"Server 收到测试消息 CmdLen[%d] No.[%d]",
				msgbin.CmdLen, recvmsg.Testno)
		}
		return
	case comm.STimeTickCommandID:
		recvmsg := &comm.STimeTickCommand{}
		recvmsg.ReadBinary([]byte(msgbin.ProtoData))
		return
	case comm.SLoginRetCommandID:
		// 收到登陆服务器返回的消息
		recvmsg := &comm.SLoginRetCommand{}
		recvmsg.ReadBinary([]byte(msgbin.ProtoData))
		if recvmsg.Loginfailed > 0 {
			this.RemoveTCPClient(client.Tempid)
			log.Error("[GBTCPClientManager.msgParseTCPClient] " +
				"连接验证失败,断开连接")
			return
		}
		client.Serverinfo = recvmsg.Taskinfo
		log.Debug("[GBTCPClientManager.msgParseTCPClient] "+
			"连接服务器验证成功,id:%d,ip:%s,port:%d",
			client.Serverinfo.Serverid, client.Serverinfo.Serverip,
			client.Serverinfo.Serverport)

		// if client.Serverinfo.Servertype == def.TypeSuperServer &&
		// 	recvmsg.Taskinfo.Servertype == def.TypeSuperServer &&
		// 	msg.GetGBServerConfigM().Myserverinfo.Servertype != def.TypeSuperServer {
		// 	// 如果登陆 SuperServer 成功，且当前不是 SuperServer
		// 	msg.GetGBServerConfigM().Myserverinfo = recvmsg.Clientinfo
		// 	log.Debug("[GBTCPClientManager.msgParseTCPClient] "+
		// 		"获取本机信息成功"+
		// 		" ServerID[%d] IP[%s] Port[%d] HTTPPort[%d] Name[%s]",
		// 		recvmsg.Clientinfo.Serverid,
		// 		recvmsg.Clientinfo.Serverip,
		// 		recvmsg.Clientinfo.Serverport,
		// 		recvmsg.Clientinfo.Httpport,
		// 		recvmsg.Clientinfo.Servername)
		// 	msg.GetGBServerConfigM().RedisConfig = recvmsg.Redisinfo
		// }
		return
	case comm.SLogoutCommandID:
		// 服务器已主动关闭，不再尝试连接它了
		log.Debug("[msgParseTCPClient] 服务器已主动关闭，不再尝试连接它了 "+
			"ServerInfo[%s]", client.Serverinfo.GetJson())
		this.serverexitchanmutex.Lock()
		defer this.serverexitchanmutex.Unlock()
		this.subnetManager.ServerConfigs.
			RemoveServerConfig(client.Serverinfo.Serverid)
		return
	case comm.SStartMyNotifyCommandID:
		// 收到依赖我的服务器的信息
		recvmsg := &comm.SStartMyNotifyCommand{}
		recvmsg.ReadBinary([]byte(msgbin.ProtoData))
		serverinfo := recvmsg.Serverinfo
		log.Debug("[GBTCPClientManager.msgParseTCPClient] "+
			"收到服务器启动信息"+
			" ServerID[%d] IP[%s] Port[%d] HTTPPort[%d] Name[%s]",
			serverinfo.Serverid, serverinfo.Serverip,
			serverinfo.Serverport, serverinfo.Httpport,
			serverinfo.Servername)
		this.serverexitchanmutex.Lock()
		defer this.serverexitchanmutex.Unlock()
		this.subnetManager.ServerConfigs.AddServerConfig(serverinfo)
		return
	case comm.SNotifyAllInfoID:
		// 收到所有服务器的配置信息
		recvmsg := &comm.SNotifyAllInfo{}
		recvmsg.ReadBinary([]byte(msgbin.ProtoData))
		this.serverexitchanmutex.Lock()
		defer this.serverexitchanmutex.Unlock()
		// if client.Serverinfo.Servertype == def.TypeSuperServer {
		// 	this.subnetManager.ServerConfigs.CleanServerConfig()
		// }
		log.Debug("[GBTCPClientManager.msgParseTCPClient] " +
			"收到所有服务器列表信息")
		// 所有服务器信息列表
		for i := 0; i < len(recvmsg.Serverinfos); i++ {
			serverinfo := recvmsg.Serverinfos[i]
			this.subnetManager.ServerConfigs.AddServerConfig(serverinfo)
		}
		return
	}
	msgqueues := &ConnectMsgQueueStruct{}
	msgqueues.task = client
	msgqueues.msg = msgbin
	this.MultiQueueControl(msgqueues)
}
