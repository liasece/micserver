package subnet

import (
	"fmt"
	"github.com/liasece/micserver/comm"
	"github.com/liasece/micserver/conf"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/server/subnet/serconfs"
	"github.com/liasece/micserver/tcpconn"
	"github.com/liasece/micserver/util"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type ConnectMsgQueueStruct struct {
	task *tcpconn.ServerConn
	msg  *msg.MessageBinary
}

func CheckServerType(servertype uint32) bool {
	if servertype <= 0 || servertype > 10 {
		return false
	}
	return true
}

// websocket连接管理器
type SubnetManager struct {
	serverhandler IServerHandler
	servermanger  ServerManager

	clientManager *GBTCPClientManager
	taskManager   *SubnetManager

	TopConfigs serconfs.TopConfigsManager // 所有服务器信息

	moudleConf    *conf.TopConfig
	subnetHandler SubnetHandler

	connPool tcpconn.ServerConnPool

	// runningMsgQueue  int
	runningMsgChan   []chan *ConnectMsgQueueStruct
	maxRunningMsgNum int

	// lastchan         int
	lastwarningtime1 uint32
	lastwarningtime2 uint32

	subnetManager *SubnetManager
}

func (this *SubnetManager) InitManager() {

	this.clientManager = &GBTCPClientManager{}
	this.clientManager.subnetManager = this
	this.clientManager.connPool.
		Init(tcpconn.ServerSCTypeClient, 1)
	this.clientManager.superexitchan = make(chan bool, 1)

	this.taskManager = &SubnetManager{}
	this.taskManager.subnetManager = this
	this.taskManager.connPool.
		Init(tcpconn.ServerSCTypeTask, 2)
}

func (this *SubnetManager) GetClientManager() *GBTCPClientManager {
	return this.clientManager
}

func (this *SubnetManager) GetTaskManager() *SubnetManager {
	return this.taskManager
}

func (this *SubnetManager) GetLatestVersionTopConfigByType(servertype uint32) uint64 {
	latestVersion := uint64(0)
	this.TopConfigs.RangeTopConfig(
		func(value comm.SServerInfo) bool {
			if value.Servertype == servertype &&
				value.Version > latestVersion {
				latestVersion = value.Version
			}
			return true
		})
	return latestVersion
}

func (this *SubnetManager) GetServerManger() ServerManager {
	return this.servermanger
}

func (this *SubnetManager) GetServerHandler() IServerHandler {
	return this.serverhandler
}

type IServerHandler interface {
	OnInit()
	OnFinal()
	OnCreateTCPConnect(serverconn *tcpconn.ServerConn)
	OnRemoveTCPConnect(serverconn *tcpconn.ServerConn)

	TCPMsgParse(serverconn *tcpconn.ServerConn,
		msgbin *msg.MessageBinary)
	GetTCPMsgParseChan(serverconn *tcpconn.ServerConn,
		maxchan int, msg *msg.MessageBinary) int
}

type ServerManager interface {
	OnSignal(os.Signal)
	NotifyOtherMyInfo()
}

var serverLoginMutex sync.Mutex

func (this *SubnetManager) OnServerLogin(tcptask *tcpconn.ServerConn,
	tarinfo *comm.SLoginCommand) {
	serverLoginMutex.Lock()
	defer serverLoginMutex.Unlock()

	// 来源服务器请求登陆本服务器
	oldtcptask := this.taskManager.GetTCPTask(uint64(tarinfo.Serverid))
	if oldtcptask != nil {
		// 已经连接成功过了，非法连接
		log.Error("[SubNetManager.OnServerLogin] "+
			"收到了重复的Server连接请求 Msg[%s]",
			tarinfo.GetJson())
		retmsg := &comm.SLoginRetCommand{}
		retmsg.Loginfailed = 1
		tcptask.SendCmd(retmsg)
		tcptask.Terminate()
		return
	}
	if !CheckServerType(tarinfo.Servertype) {
		// 未知的服务器类型
		log.Error("[SubNetManager.OnServerLogin] "+
			"收到未知的服务器类型 Msg[%s]",
			tarinfo.GetJson())
		retmsg := &comm.SLoginRetCommand{}
		retmsg.Loginfailed = 1
		tcptask.SendCmd(retmsg)
		tcptask.Terminate()
		return
	}
	var TopConfig comm.SServerInfo

	// 来源服务器地址
	remoteaddr := tcptask.GetConn().RemoteAddr().String()
	// 检查来源服务器IP地址与消息中的地址一致
	sameip := false
	s := strings.Split(remoteaddr, ":")
	if len(s) == 2 && s[0] == tarinfo.Serverip {
		// IP 地址正确
		sameip = true
	}
	// 获取来源服务器ID在本地的配置
	TopConfig = this.TopConfigs.
		GetTopConfigByInfo(tarinfo)
	// 如果成功获取到了一个serverid
	if TopConfig.Serverid != 0 {
		// 检查该配置是否已被占用
		oldtcptask = this.taskManager.GetTCPTask(uint64(tarinfo.Serverid))
		if oldtcptask != nil {
			// 已经连接成功过了，非法连接
			log.Error("[SubNetManager.OnServerLogin] "+
				"收到了重复的Server连接请求 2 Msg[%s]",
				tarinfo.GetJson())
			retmsg := &comm.SLoginRetCommand{}
			retmsg.Loginfailed = 1
			tcptask.SendCmd(retmsg)
			tcptask.Terminate()
			return
		}
	}
	// 检查是否获取信息成功
	if TopConfig.Serverid == 0 || !sameip ||
		TopConfig.Servertype != tarinfo.Servertype {
		// 如果获取信息不成功
		log.Error("[SubNetManager.OnServerLogin] "+
			"连接分配异常 未知服务器连接 "+
			"Addr[%s] Msg[%s]",
			remoteaddr, tarinfo.GetJson())
		retmsg := &comm.SLoginRetCommand{}
		retmsg.Loginfailed = 1
		tcptask.SendCmd(retmsg)
		tcptask.Terminate()

		// if this.moudleConf.Myserverinfo.Servertype !=
		// 	def.TypeSuperServer {
		// 	requestServerInfo := &comm.SRequestServerInfo{}
		// 	this.clientManager.BroadcastByType(def.TypeSuperServer,
		// 		requestServerInfo)
		// }
		return
	}

	TopConfig.Version = tarinfo.Version

	// 来源服务器检查完毕
	// 完善来源服务器在本服务器的信息
	tcptask.SetVertify(true)
	tcptask.SetTerminateTime(0) // 清除终止时间状态
	tcptask.Serverinfo = TopConfig
	this.taskManager.ChangeTCPTaskTempid(tcptask,
		uint64(tcptask.Serverinfo.Serverid))
	log.Debug("[SubNetManager.OnServerLogin] "+
		"客户端连接验证成功 "+
		" SerID[%d] IP[%s] Type[%d] Port[%d] Name[%s]",
		TopConfig.Serverid, TopConfig.Serverip,
		TopConfig.Servertype, TopConfig.Serverport,
		TopConfig.Servername)
	// 向来源服务器回复登陆成功消息
	retmsg := &comm.SLoginRetCommand{}
	retmsg.Loginfailed = 0
	retmsg.Clientinfo = TopConfig
	// retmsg.Taskinfo = this.moudleConf.Myserverinfo
	// retmsg.Redisinfo = this.moudleConf.RedisConfig
	tcptask.SendCmd(retmsg)

	// 通知其他服务器，次服务器登陆完成
	// 如果我是SuperServer
	// 向来源服务器发送本地已有的所有服务器信息
	this.taskManager.NotifyAllServerInfo(tcptask)
	// 把来源服务器信息广播给其它所有服务器
	notifymsg := &comm.SStartMyNotifyCommand{}
	notifymsg.Serverinfo = TopConfig
	this.taskManager.BroadcastAll(notifymsg)
}

// 绑定我的HTTP服务 会阻塞
func (this *SubnetManager) BindMyHttpServer() {
	usehttps := this.moudleConf.GetProp("usehttps")
	httpportstr := fmt.Sprintf(":%d", 8080)
	if usehttps == "true" {
		log.Debug("[SubNetManager.BindMyHttpServer] "+
			"服务器https绑定成功 HTTPPort[%s] use https",
			httpportstr)
		httpscertfile := this.moudleConf.
			GetProp("httpscertfile")
		httpskeyfile := this.moudleConf.
			GetProp("httpskeyfile")
		err := http.ListenAndServeTLS(httpportstr, httpscertfile,
			httpskeyfile, nil)
		if err != nil {
			panic("ListenAndServe: " + err.Error())
			// return err
		}
	} else {
		log.Debug("[SubNetManager.BindMyHttpServer] "+
			"服务器http绑定成功,%s", httpportstr)
		err := http.ListenAndServe(httpportstr, nil)
		if err != nil {
			panic("ListenAndServe: " + err.Error())
			// return err
		}
	}
}

func (this *SubnetManager) BindTCPSubnet(settings map[string]string) error {
	addr, hasconf := settings["subnettcpaddr"]
	if !hasconf {
		return fmt.Errorf("subnettcpaddr hasn't set.")
	}
	// init tcp subnet port
	netlisten, err := net.Listen("tcp", addr)
	if err != nil {
		log.Error("[SubNetManager.BindTCPServer] "+
			"服务器绑定失败 IPPort[%s] Err[%s]",
			addr, err.Error())
		return err
	}
	// myservertype := this.moudleConf.Myserverinfo.Servertype
	log.Debug("[SubNetManager.BindTCPServer] "+
		"服务器绑定成功 IPPort[%s]", addr)

	go this.TCPServerListenerProcess(netlisten)
	return nil
}

func (this *SubnetManager) TCPServerListenerProcess(listener net.Listener) {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := util.GetPanicInfo(recover()); err != nil {
			log.Error("[TCPServerListenerProcess] "+
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
		if err, stackInfo := util.GetPanicInfo(recover()); err != nil {
			// 这里的err其实就是panic传入的内容
			log.Error("mTCPServerListener "+
				"Panic: ErrName[%v] \n Stack[%s]", err, stackInfo)
		}
	}()

	for true {
		conn, err := listener.Accept()
		if err != nil {
			log.Error("[SubNetManager.TCPServerListenerProcess] "+
				"服务器端口监听异常 Err[%s]",
				err.Error())
			continue
		}
		log.Debug("[SubNetManager.BindTCPServer] "+
			"收到新的TCP连接 Addr[%s]",
			conn.RemoteAddr().String())
		tcptask := this.taskManager.AddTCPTask(conn)
		if tcptask != nil {
			this.subnetHandler.OnCreateTCPConnect(tcptask)
			go this.taskManager.handleConnection(tcptask)
		}
	}
}

func (this *SubnetManager) InitMsgQueue(sum int) {
	// 最大同时处理的消息数量
	this.maxRunningMsgNum = sum // 消息队列线程数
	if this.maxRunningMsgNum < 1 {
		this.maxRunningMsgNum = 1
		log.Error("[SubnetManager.InitMsgQueue] " +
			"消息处理线程数量过小，置为1...")
	}
	this.runningMsgChan = make([]chan *ConnectMsgQueueStruct,
		this.maxRunningMsgNum)
	log.Debug("[SubnetManager.InitMsgQueue] "+
		"Task 消息处理线程数量 ThreadNum[%d]", this.maxRunningMsgNum)
	i := 0
	for i < this.maxRunningMsgNum {
		this.runningMsgChan[i] = make(chan *ConnectMsgQueueStruct,
			15000)
		go this.RecvmsgProcess(i)
		i++
	}
}

func (this *SubnetManager) RangeConnect(
	callback func(*tcpconn.ServerConn) bool) {
	this.connPool.Range(callback)
}

// 并行处理接收消息队列数据
func (this *SubnetManager) MultiRecvmsgQueue(
	index int) (normalreturn bool) {
	if this.runningMsgChan == nil || this.runningMsgChan[index] == nil {
		return true
	}
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := util.GetPanicInfo(recover()); err != nil {
			// 这里的err其实就是panic传入的内容
			log.Error("[SubnetManager.MultiRecvmsgQueue] "+
				"Panic: ErrName[%v] \n Stack[%s]", err, stackInfo)
			normalreturn = false
		}
	}()
	normalreturn = true

	msgchan := this.runningMsgChan[index]
	for msgqueues := range msgchan {
		curtime := uint32(time.Now().Unix())
		if curtime > msgqueues.msg.TimeStamp+1 &&
			curtime > this.lastwarningtime2 {
			log.Error("[SubnetManager.MultiRecvmsgQueue] "+
				"[消耗时间统计]服务器TCPTask处理消息延迟很严重"+
				" TimeInc[%d] in Thread[%d]",
				curtime-msgqueues.msg.TimeStamp, index)
			this.lastwarningtime2 = curtime
		}
		functiontime := util.FunctionTime{}
		functiontime.Start("MultiRecvmsgQueue")
		this.serverhandler.TCPMsgParse(msgqueues.task, msgqueues.msg)
		functiontime.Stop()
	}
	return true
}

func (this *SubnetManager) RecvmsgProcess(index int) {
	for {
		if this.MultiRecvmsgQueue(index) {
			// 正常退出
			break
		}
	}
}

// 分配消息处理线程
func (this *SubnetManager) MultiQueueControl(
	msgqueues *ConnectMsgQueueStruct) {
	if this.maxRunningMsgNum < 1 {
		this.serverhandler.TCPMsgParse(msgqueues.task, msgqueues.msg)
		return
	}
	who := this.serverhandler.GetTCPMsgParseChan(msgqueues.task,
		this.maxRunningMsgNum, msgqueues.msg)
	this.runningMsgChan[who] <- msgqueues
}

func (this *SubnetManager) AddTCPTask(
	conn net.Conn) *tcpconn.ServerConn {
	tcptask := this.connPool.NewServerConn(conn, 0)

	// 2秒以后没有收到验证消息就断开连接
	// 10秒还未通过验证就断开连接
	curtime := uint64(time.Now().Unix())
	tcptask.SetTerminateTime(uint64(curtime + 10))

	log.Debug("[SubnetManager.AddTCPTask] "+
		"增加新的连接数 TmpID[%d] 当前连接数量 ConnNum[%d]",
		tcptask.Tempid, this.connPool.Len())

	return tcptask
}
func (this *SubnetManager) GetTCPTask(
	tempid uint64) *tcpconn.ServerConn {
	return this.connPool.Get(tempid)
}

func (this *SubnetManager) RemoveTCPTask(tempid uint64) {
	this.connPool.Remove(tempid)
}

// 修改链接的 tempip
func (this *SubnetManager) ChangeTCPTaskTempid(
	tcptask *tcpconn.ServerConn, newtempid uint64) {
	this.connPool.ChangeTempid(tcptask, newtempid)
}

func (this *SubnetManager) BroadcastByType(
	servertype uint32, v msg.MsgStruct) {
	this.connPool.BroadcastByType(servertype, v)
}

func (this *SubnetManager) BroadcastAll(v msg.MsgStruct) {
	this.connPool.BroadcastCmd(v)
}

//通知所有服务器列表信息
func (this *SubnetManager) NotifyAllServerInfo(
	tcptask *tcpconn.ServerConn) {
	retmsg := &comm.SNotifyAllInfo{}
	retmsg.Serverinfos = make([]comm.SServerInfo, 0)
	this.subnetManager.TopConfigs.RangeTopConfig(func(
		value comm.SServerInfo) bool {
		retmsg.Serverinfos = append(retmsg.Serverinfos, value)
		return true
	})
	if len(retmsg.Serverinfos) > 0 {
		log.Debug("[NotifyAllServerInfo] 发送所有服务器列表信息 Msg[%s]",
			retmsg.GetJson())
		tcptask.SendCmd(retmsg)
	}
}

func (this *SubnetManager) handleConnection(
	tcptask *tcpconn.ServerConn) {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := util.GetPanicInfo(recover()); err != nil {
			log.Error("[handleConnection] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
		}
	}()
	timer_10_sec := uint64(time.Now().Unix()) + 10
	netbuffer := util.NewIOBuffer(tcptask.Conn.Conn, 6400*1024)
	msgReader := msg.NewMessageBinaryReader(netbuffer)
	for {
		curtime := uint64(time.Now().Unix())
		if tcptask.IsTerminateTimeout(curtime) {
			this.RemoveTCPTask(tcptask.Tempid)
			log.Error("[SubnetManager.handleConnection] "+
				"长时间未通过验证，断开连接 TmpID[%d]",
				tcptask.Tempid)
			return
		}
		if tcptask.IsTerminateForce() {
			this.subnetHandler.OnRemoveTCPConnect(tcptask)
			this.RemoveTCPTask(tcptask.Tempid)
			log.Error("[SubnetManager.handleConnection] "+
				"服务器主动断开连接 TmpID[%d]", tcptask.Tempid)
			return
		}
		if false {
			log.Debug("[SubnetManager.handleConnection] "+
				"服务器准备停机了,退出连接处理 TmpID[%d]",
				tcptask.Tempid)
			return
		}
		// 发送心跳消息
		if timer_10_sec <= curtime {
			timer_10_sec = curtime + 10
			sendmsg := &comm.STimeTickCommand{}
			sendmsg.Testno = uint32(curtime)
			tcptask.SendCmd(sendmsg)
		}
		// 250毫秒以后读不到数据超时
		derr := tcptask.Conn.Conn.SetReadDeadline(time.Now().
			Add(time.Duration(time.Millisecond * 250)))
		if derr != nil {
			log.Error("[handleConnection] SetReadDeadline Err[%s]"+
				"ServerName[%s] ServerID[%d]",
				derr.Error(), tcptask.Serverinfo.Servername,
				tcptask.Serverinfo.Serverid)
			this.subnetHandler.OnRemoveTCPConnect(tcptask)
			this.RemoveTCPTask(tcptask.Tempid)
			return
		}
		_, err := netbuffer.ReadFromReader()
		if err != nil {
			if err == io.EOF {
				log.Debug("[SubnetManager.handleConnection] "+
					"tcptask数据读写异常 scoket返回 TmpID[%d] Error[%s]",
					tcptask.Tempid, err.Error())
				this.subnetHandler.OnRemoveTCPConnect(tcptask)
				this.RemoveTCPTask(tcptask.Tempid)
				return
			} else {
				continue
			}
		}
		thismsgtimes := 0

		functiontime := util.FunctionTime{}
		functiontime.Start("handleTaskConnection")

		pcknum := 0
		// pcksize := msgbuf.Len()

		err = msgReader.RangeMsgBinary(func(msgbinary *msg.MessageBinary) {
			// 判断消息是否阻塞严重
			curtime := uint32(time.Now().Unix())
			if curtime > msgbinary.TimeStamp+1 &&
				curtime > this.lastwarningtime1 {
				log.Error("[SubnetManager.handleConnection] "+
					"[消耗时间统计]服务器TCPTask接收消息延迟很严重"+
					" TimeInc[%d]",
					curtime-msgbinary.TimeStamp)
				this.lastwarningtime1 = curtime
			}
			// 重置消息标签为接收时间，用于处理阻塞判断
			msgbinary.TimeStamp = curtime
			// 处理消息
			this.msgParse(tcptask, msgbinary)
			pcknum++

			thismsgtimes += 1
		})
		functiontime.Stop()
		if err != nil {
			log.Error("[SubnetManager.handleConnection] "+
				"RangeMsgBinary读消息失败，断开连接 Err[%s]", err.Error())
			this.subnetHandler.OnRemoveTCPConnect(tcptask)
			this.RemoveTCPTask(tcptask.Tempid)
			return
		}
	}
}

func (this *SubnetManager) msgParse(tcptask *tcpconn.ServerConn, msgbin *msg.MessageBinary) {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := util.GetPanicInfo(recover()); err != nil {
			// 这里的err其实就是panic传入的内容
			log.Error("[SubnetManager.msgParse] "+
				"Panic: ErrName[%v] \n Stack[%s]", err, stackInfo)
		}
	}()

	cmdname := comm.MsgIdToString(msgbin.CmdID)
	log.Debug("[TCPTask.msgParse] "+
		"收到 MsgName[%s] CmdLen[%d]", cmdname, msgbin.CmdLen)
	if msgbin.CmdID == comm.STestCommandID {
		// 来源服务器请求测试
		performance_test := this.subnetManager.moudleConf.
			GetProp("performance_test")
		if performance_test == "true" {
			recvmsg := &comm.STestCommand{}
			recvmsg.ReadBinary([]byte(msgbin.ProtoData))
			log.Debug("[TCPTask.msgParse] "+
				"Server 收到测试消息 CmdLen[%d] No.[%d]",
				msgbin.CmdLen, recvmsg.Testno)
		}
		return
	} else if msgbin.CmdID == comm.SLoginCommandID {
		recvmsg := &comm.SLoginCommand{}
		recvmsg.ReadBinary([]byte(msgbin.ProtoData))
		this.subnetManager.OnServerLogin(tcptask, recvmsg)
		return
	} else if msgbin.CmdID == comm.SLogoutCommandID {
		// 服务器退出登陆
		log.Debug("[TCPTask.msgParse] "+
			"服务器 ServerName[%s] 退出登陆了", tcptask.Serverinfo.Servername)
		tcptask.SetVertify(false)
		this.subnetManager.TopConfigs.
			RemoveTopConfig(tcptask.Serverinfo.Serverid)
		return
	}

	// 如果来源服务器未登陆或身份校验未通过
	if !tcptask.IsVertify() {
		log.Debug("[TCPTask.msgParse] "+
			"客户端连接验证异常 MsgName[%s] ServerID[%d]",
			cmdname, tcptask.Serverinfo.Serverid)
		return
	}

	if msgbin.CmdID == comm.SRequestServerInfoID {
		this.NotifyAllServerInfo(tcptask)
	}
	// 以下消息只有来源服务器通过了连接验证才会处理

	// 消息队列启用
	msgqueues := &ConnectMsgQueueStruct{}
	msgqueues.task = tcptask
	msgqueues.msg = msgbin
	this.MultiQueueControl(msgqueues)
}
