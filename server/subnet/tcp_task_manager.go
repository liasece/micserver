/**
 * \file GBTCPTaskManager.go
 * \version
 * \author wzy
 * \date  2018年01月17日 11:55:27
 * \brief 处理TCP连接相关
 *
 */

package subnet

import (
	"base"
	"base/functime"
	"base/logger"
	"base/tcpconn"
	"base/util"
	// "bytes"
	// "encoding/binary"
	// "encoding/json"
	// "fmt"
	"io"
	// "math/rand"
	"net"
	"servercomm"
	// "sort"
	// "strconv"
	// "strings"
	// "encoding/hex"
	// "sync"
	"time"
)

// websocket连接管理器
type GBTCPTaskManager struct {
	connPool tcpconn.ServerConnPool

	serverhandler IServerHandler

	// runningMsgQueue  int
	runningMsgChan   []chan *ConnectMsgQueueStruct
	maxRunningMsgNum int

	// lastchan         int
	lastwarningtime1 uint32
	lastwarningtime2 uint32
}

func (this *GBTCPTaskManager) InitMsgQueue(sum int) {
	// 最大同时处理的消息数量
	this.maxRunningMsgNum = sum // 消息队列线程数
	if this.maxRunningMsgNum < 1 {
		this.maxRunningMsgNum = 1
		logger.Error("[GBTCPTaskManager.InitMsgQueue] " +
			"消息处理线程数量过小，置为1...")
	}
	this.runningMsgChan = make([]chan *ConnectMsgQueueStruct,
		this.maxRunningMsgNum)
	logger.Debug("[GBTCPTaskManager.InitMsgQueue] "+
		"Task 消息处理线程数量 ThreadNum[%d]", this.maxRunningMsgNum)
	i := 0
	for i < this.maxRunningMsgNum {
		this.runningMsgChan[i] = make(chan *ConnectMsgQueueStruct,
			15000)
		go this.RecvmsgProcess(i)
		i++
	}
}

func (this *GBTCPTaskManager) RangeConnect(
	callback func(*tcpconn.ServerConn) bool) {
	this.connPool.Range(callback)
}

// 并行处理接收消息队列数据
func (this *GBTCPTaskManager) MultiRecvmsgQueue(
	index int) (normalreturn bool) {
	if this.runningMsgChan == nil || this.runningMsgChan[index] == nil {
		return true
	}
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := util.GetPanicInfo(recover()); err != nil {
			// 这里的err其实就是panic传入的内容
			logger.Error("[GBTCPTaskManager.MultiRecvmsgQueue] "+
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
			logger.Error("[GBTCPTaskManager.MultiRecvmsgQueue] "+
				"[消耗时间统计]服务器TCPTask处理消息延迟很严重"+
				" TimeInc[%d] in Thread[%d]",
				curtime-msgqueues.msg.TimeStamp, index)
			this.lastwarningtime2 = curtime
		}
		functiontime := functime.FunctionTime{}
		functiontime.Start("MultiRecvmsgQueue")
		this.serverhandler.TCPMsgParse(msgqueues.task, msgqueues.msg)
		functiontime.Stop()
	}
	return true
}

func (this *GBTCPTaskManager) RecvmsgProcess(index int) {
	for {
		if this.MultiRecvmsgQueue(index) {
			// 正常退出
			break
		}
	}
}

// 分配消息处理线程
func (this *GBTCPTaskManager) MultiQueueControl(
	msgqueues *ConnectMsgQueueStruct) {
	if this.maxRunningMsgNum < 1 {
		this.serverhandler.TCPMsgParse(msgqueues.task, msgqueues.msg)
		return
	}
	who := this.serverhandler.GetTCPMsgParseChan(msgqueues.task,
		this.maxRunningMsgNum, msgqueues.msg)
	this.runningMsgChan[who] <- msgqueues
}

func (this *GBTCPTaskManager) AddTCPTask(
	conn net.Conn) *tcpconn.ServerConn {
	tcptask := this.connPool.NewServerConn(conn, 0)

	// 2秒以后没有收到验证消息就断开连接
	// 10秒还未通过验证就断开连接
	curtime := uint64(time.Now().Unix())
	tcptask.SetTerminateTime(uint64(curtime + 10))

	logger.Debug("[GBTCPTaskManager.AddTCPTask] "+
		"增加新的连接数 TmpID[%d] 当前连接数量 ConnNum[%d]",
		tcptask.Tempid, this.connPool.Len())

	return tcptask
}
func (this *GBTCPTaskManager) GetTCPTask(
	tempid uint64) *tcpconn.ServerConn {
	return this.connPool.Get(tempid)
}

func (this *GBTCPTaskManager) RemoveTCPTask(tempid uint64) {
	this.connPool.Remove(tempid)
}

// 修改链接的 tempip
func (this *GBTCPTaskManager) ChangeTCPTaskTempid(
	tcptask *tcpconn.ServerConn, newtempid uint64) {
	this.connPool.ChangeTempid(tcptask, newtempid)
}

func (this *GBTCPTaskManager) BroadcastByType(
	servertype uint32, v base.MsgStruct) {
	this.connPool.BroadcastByType(servertype, v)
}

func (this *GBTCPTaskManager) BroadcastAll(v base.MsgStruct) {
	this.connPool.BroadcastCmd(v)
}

//通知所有服务器列表信息
func (this *GBTCPTaskManager) NotifyAllServerInfo(
	tcptask *tcpconn.ServerConn) {
	retmsg := &servercomm.SNotifyAllInfo{}
	retmsg.Serverinfos = make([]servercomm.SServerInfo, 0)
	GetSubnetManager().ServerConfigs.RangeServerConfig(func(
		value servercomm.SServerInfo) bool {
		retmsg.Serverinfos = append(retmsg.Serverinfos, value)
		return true
	})
	if len(retmsg.Serverinfos) > 0 {
		logger.Debug("[NotifyAllServerInfo] 发送所有服务器列表信息 Msg[%s]",
			retmsg.GetJson())
		tcptask.SendCmd(retmsg)
	}
}

func handleConnection(
	tcptask *tcpconn.ServerConn, server IServerHandler) {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := util.GetPanicInfo(recover()); err != nil {
			logger.Error("[handleConnection] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
		}
	}()
	timer_10_sec := uint64(time.Now().Unix()) + 10
	netbuffer := base.NewIOBuffer(tcptask.Conn.Conn, 6400*1024)
	msgReader := base.NewMessageBinaryReader(netbuffer)
	for {
		curtime := uint64(time.Now().Unix())
		if tcptask.IsTerminateTimeout(curtime) {
			GetGBTCPTaskManager().RemoveTCPTask(tcptask.Tempid)
			logger.Error("[GBTCPTaskManager.handleConnection] "+
				"长时间未通过验证，断开连接 TmpID[%d]",
				tcptask.Tempid)
			return
		}
		if tcptask.IsTerminateForce() {
			server.OnRemoveTCPConnect(tcptask)
			GetGBTCPTaskManager().RemoveTCPTask(tcptask.Tempid)
			logger.Error("[GBTCPTaskManager.handleConnection] "+
				"服务器主动断开连接 TmpID[%d]", tcptask.Tempid)
			return
		}
		if base.GetGBServerConfigM().TerminateServer {
			logger.Debug("[GBTCPTaskManager.handleConnection] "+
				"服务器准备停机了,退出连接处理 TmpID[%d]",
				tcptask.Tempid)
			return
		}
		// 发送心跳消息
		if timer_10_sec <= curtime {
			timer_10_sec = curtime + 10
			sendmsg := &servercomm.STimeTickCommand{}
			sendmsg.Testno = uint32(curtime)
			tcptask.SendCmd(sendmsg)
		}
		// 250毫秒以后读不到数据超时
		derr := tcptask.Conn.Conn.SetReadDeadline(time.Now().
			Add(time.Duration(time.Millisecond * 250)))
		if derr != nil {
			logger.Error("[handleConnection] SetReadDeadline Err[%s]"+
				"ServerName[%s] ServerID[%d]",
				derr.Error(), tcptask.Serverinfo.Servername,
				tcptask.Serverinfo.Serverid)
			server.OnRemoveTCPConnect(tcptask)
			GetGBTCPTaskManager().RemoveTCPTask(tcptask.Tempid)
			return
		}
		_, err := netbuffer.ReadFromReader()
		if err != nil {
			if err == io.EOF {
				logger.Debug("[GBTCPTaskManager.handleConnection] "+
					"tcptask数据读写异常 scoket返回 TmpID[%d] Error[%s]",
					tcptask.Tempid, err.Error())
				server.OnRemoveTCPConnect(tcptask)
				GetGBTCPTaskManager().RemoveTCPTask(tcptask.Tempid)
				return
			} else {
				continue
			}
		}
		thismsgtimes := 0

		functiontime := functime.FunctionTime{}
		functiontime.Start("handleTaskConnection")

		pcknum := 0
		// pcksize := msgbuf.Len()

		err = msgReader.RangeMsgBinary(func(msgbinary *base.MessageBinary) {
			// 判断消息是否阻塞严重
			curtime := uint32(time.Now().Unix())
			if curtime > msgbinary.TimeStamp+1 &&
				curtime > GetGBTCPTaskManager().lastwarningtime1 {
				logger.Error("[GBTCPTaskManager.handleConnection] "+
					"[消耗时间统计]服务器TCPTask接收消息延迟很严重"+
					" TimeInc[%d]",
					curtime-msgbinary.TimeStamp)
				GetGBTCPTaskManager().lastwarningtime1 = curtime
			}
			// 重置消息标签为接收时间，用于处理阻塞判断
			msgbinary.TimeStamp = curtime
			// 处理消息
			msgParse(tcptask, msgbinary, server)
			pcknum++

			thismsgtimes += 1
		})
		functiontime.Stop()
		if err != nil {
			logger.Error("[GBTCPTaskManager.handleConnection] "+
				"RangeMsgBinary读消息失败，断开连接 Err[%s]", err.Error())
			server.OnRemoveTCPConnect(tcptask)
			GetGBTCPTaskManager().RemoveTCPTask(tcptask.Tempid)
			return
		}
	}
}

func msgParse(tcptask *tcpconn.ServerConn, msgbin *base.MessageBinary,
	server IServerHandler) {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := util.GetPanicInfo(recover()); err != nil {
			// 这里的err其实就是panic传入的内容
			logger.Error("[GBTCPTaskManager.msgParse] "+
				"Panic: ErrName[%v] \n Stack[%s]", err, stackInfo)
		}
	}()

	cmdname := servercomm.MsgIdToString(msgbin.CmdID)
	logger.Debug("[TCPTask.msgParse] "+
		"收到 MsgName[%s] CmdLen[%d]", cmdname, msgbin.CmdLen)
	if msgbin.CmdID == servercomm.STestCommandID {
		// 来源服务器请求测试
		performance_test := base.GetGBServerConfigM().
			GetProp("performance_test")
		if performance_test == "true" {
			recvmsg := &servercomm.STestCommand{}
			recvmsg.ReadBinary([]byte(msgbin.ProtoData))
			logger.Debug("[TCPTask.msgParse] "+
				"Server 收到测试消息 CmdLen[%d] No.[%d]",
				msgbin.CmdLen, recvmsg.Testno)
		}
		return
	} else if msgbin.CmdID == servercomm.SLoginCommandID {
		recvmsg := &servercomm.SLoginCommand{}
		recvmsg.ReadBinary([]byte(msgbin.ProtoData))
		OnServerLogin(tcptask, recvmsg)
		return
	} else if msgbin.CmdID == servercomm.SLogoutCommandID {
		// 服务器退出登陆
		logger.Debug("[TCPTask.msgParse] "+
			"服务器 ServerName[%s] 退出登陆了", tcptask.Serverinfo.Servername)
		tcptask.SetVertify(false)
		GetSubnetManager().ServerConfigs.
			RemoveServerConfig(tcptask.Serverinfo.Serverid)
		return
	}

	// 如果来源服务器未登陆或身份校验未通过
	if !tcptask.IsVertify() {
		logger.Debug("[TCPTask.msgParse] "+
			"客户端连接验证异常 MsgName[%s] ServerID[%d]",
			cmdname, tcptask.Serverinfo.Serverid)
		return
	}

	if msgbin.CmdID == servercomm.SRequestServerInfoID {
		GetGBTCPTaskManager().NotifyAllServerInfo(tcptask)
	}
	// 以下消息只有来源服务器通过了连接验证才会处理

	// 消息队列启用
	msgqueues := &ConnectMsgQueueStruct{}
	msgqueues.task = tcptask
	msgqueues.msg = msgbin
	GetGBTCPTaskManager().MultiQueueControl(msgqueues)
}
