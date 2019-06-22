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
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/tcpconn"
	"github.com/liasece/micserver/util"
	// "bytes"
	// "encoding/binary"
	// "encoding/json"
	// "fmt"
	"io"
	// "math/rand"
	"github.com/liasece/micserver/comm"
	"github.com/liasece/micserver/msg"
	"net"
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

	subnetManager *SubnetManger
}

func (this *GBTCPTaskManager) InitMsgQueue(sum int) {
	// 最大同时处理的消息数量
	this.maxRunningMsgNum = sum // 消息队列线程数
	if this.maxRunningMsgNum < 1 {
		this.maxRunningMsgNum = 1
		log.Error("[GBTCPTaskManager.InitMsgQueue] " +
			"消息处理线程数量过小，置为1...")
	}
	this.runningMsgChan = make([]chan *ConnectMsgQueueStruct,
		this.maxRunningMsgNum)
	log.Debug("[GBTCPTaskManager.InitMsgQueue] "+
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
			log.Error("[GBTCPTaskManager.MultiRecvmsgQueue] "+
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
			log.Error("[GBTCPTaskManager.MultiRecvmsgQueue] "+
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

	log.Debug("[GBTCPTaskManager.AddTCPTask] "+
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
	servertype uint32, v msg.MsgStruct) {
	this.connPool.BroadcastByType(servertype, v)
}

func (this *GBTCPTaskManager) BroadcastAll(v msg.MsgStruct) {
	this.connPool.BroadcastCmd(v)
}

//通知所有服务器列表信息
func (this *GBTCPTaskManager) NotifyAllServerInfo(
	tcptask *tcpconn.ServerConn) {
	retmsg := &comm.SNotifyAllInfo{}
	retmsg.Serverinfos = make([]comm.SServerInfo, 0)
	this.subnetManager.ServerConfigs.RangeServerConfig(func(
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

func (this *GBTCPTaskManager) handleConnection(
	tcptask *tcpconn.ServerConn, server IServerHandler) {
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
			log.Error("[GBTCPTaskManager.handleConnection] "+
				"长时间未通过验证，断开连接 TmpID[%d]",
				tcptask.Tempid)
			return
		}
		if tcptask.IsTerminateForce() {
			server.OnRemoveTCPConnect(tcptask)
			this.RemoveTCPTask(tcptask.Tempid)
			log.Error("[GBTCPTaskManager.handleConnection] "+
				"服务器主动断开连接 TmpID[%d]", tcptask.Tempid)
			return
		}
		if this.subnetManager.moudleConf.TerminateServer {
			log.Debug("[GBTCPTaskManager.handleConnection] "+
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
			server.OnRemoveTCPConnect(tcptask)
			this.RemoveTCPTask(tcptask.Tempid)
			return
		}
		_, err := netbuffer.ReadFromReader()
		if err != nil {
			if err == io.EOF {
				log.Debug("[GBTCPTaskManager.handleConnection] "+
					"tcptask数据读写异常 scoket返回 TmpID[%d] Error[%s]",
					tcptask.Tempid, err.Error())
				server.OnRemoveTCPConnect(tcptask)
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
				log.Error("[GBTCPTaskManager.handleConnection] "+
					"[消耗时间统计]服务器TCPTask接收消息延迟很严重"+
					" TimeInc[%d]",
					curtime-msgbinary.TimeStamp)
				this.lastwarningtime1 = curtime
			}
			// 重置消息标签为接收时间，用于处理阻塞判断
			msgbinary.TimeStamp = curtime
			// 处理消息
			this.msgParse(tcptask, msgbinary, server)
			pcknum++

			thismsgtimes += 1
		})
		functiontime.Stop()
		if err != nil {
			log.Error("[GBTCPTaskManager.handleConnection] "+
				"RangeMsgBinary读消息失败，断开连接 Err[%s]", err.Error())
			server.OnRemoveTCPConnect(tcptask)
			this.RemoveTCPTask(tcptask.Tempid)
			return
		}
	}
}

func (this *GBTCPTaskManager) msgParse(tcptask *tcpconn.ServerConn, msgbin *msg.MessageBinary,
	server IServerHandler) {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := util.GetPanicInfo(recover()); err != nil {
			// 这里的err其实就是panic传入的内容
			log.Error("[GBTCPTaskManager.msgParse] "+
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
		this.subnetManager.ServerConfigs.
			RemoveServerConfig(tcptask.Serverinfo.Serverid)
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
