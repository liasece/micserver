package subnet

import (
	"github.com/liasece/micserver/connect"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/servercomm"
	"github.com/liasece/micserver/util"
	"time"
)

type ConnectMsgQueueStruct struct {
	conn *connect.ServerConn
	msg  *msg.MessageBinary
}

// 当TCP连接被移除时调用
func (this *SubnetManager) onConnectClose(conn *connect.ServerConn) {
	this.RemoveServerConn(conn.Tempid)
}

// 当收到TCP消息时调用
func (this *SubnetManager) OnRecvTCPMsg(conn *connect.ServerConn,
	msgbinary *msg.MessageBinary) {
	switch msgbinary.CmdID {
	case servercomm.SForwardToServerID:
		{
			// 服务器间用户空间消息转发
			if this.SubnetCallback.regForwardToServer != nil {
				layerMsg := &servercomm.SForwardToServer{}
				layerMsg.ReadBinary(msgbinary.ProtoData)
				this.SubnetCallback.regForwardToServer(layerMsg)
			}
		}
	case servercomm.SForwardFromGateID:
		{
			// Gateway 转发过来的客户端消息
			if this.SubnetCallback.regForwardFromGate != nil {
				layerMsg := &servercomm.SForwardFromGate{}
				layerMsg.ReadBinary(msgbinary.ProtoData)
				this.SubnetCallback.regForwardFromGate(layerMsg)
			}
		}
	case servercomm.SForwardToClientID:
		{
			// 其他服务器转发过来的，要发送到客户端的消息
			if this.SubnetCallback.regForwardToClient != nil {
				layerMsg := &servercomm.SForwardToClient{}
				layerMsg.ReadBinary(msgbinary.ProtoData)
				this.SubnetCallback.regForwardToClient(layerMsg)
			}
		}
	case servercomm.SUpdateSessionID:
		{
			// 客户端会话更新
			if this.SubnetCallback.regUpdateSession != nil {
				layerMsg := &servercomm.SUpdateSession{}
				layerMsg.ReadBinary(msgbinary.ProtoData)
				this.SubnetCallback.regUpdateSession(layerMsg)
			}
		}
	case servercomm.SStartMyNotifyCommandID:
	default:
		{
			msgid := msgbinary.CmdID
			msgname := servercomm.MsgIdToString(msgid)
			this.Error("[SubnetManager.OnRecvTCPMsg] 未知消息 %d:%s",
				msgid, msgname)
		}
	}
}

// 获取TCP消息的消息处理通道
func (this *SubnetManager) OnGetRecvTCPMsgParseChan(conn *connect.ServerConn,
	maxChan int32, msgbinary *msg.MessageBinary) int32 {
	return 0
}

func (this *SubnetManager) OnCreateTCPConnect(conn *connect.ServerConn) {
}

func (this *SubnetManager) onConnectRecv(conn *connect.ServerConn,
	msgbin *msg.MessageBinary) {
	if conn.GetSCType() == connect.ServerSCTypeTask {
		curtime := uint64(time.Now().Unix())
		if conn.IsTerminateTimeout(curtime) {
			this.onClientDisconnected(conn)
			this.Error("[SubnetManager.handleConnection] "+
				"长时间未通过验证，断开连接 TmpID[%d]",
				conn.Tempid)
			return
		}
		if conn.IsTerminateForce() {
			this.onClientDisconnected(conn)
			this.Debug("[SubnetManager.handleConnection] "+
				"服务器主动断开连接 TmpID[%s]", conn.Tempid)
			return
		}
	}
	// this.Debug("[SubnetManager.msgParseTCPConn] 收到消息 %s",
	// 	servercomm.MsgIdToString(msgbin.CmdID))
	switch msgbin.CmdID {
	case servercomm.STestCommandID:
		recvmsg := &servercomm.STestCommand{}
		recvmsg.ReadBinary([]byte(msgbin.ProtoData))
		this.Debug("[SubnetManager.msgParseTCPConn] "+
			"Server 收到测试消息 CmdLen[%d] No.[%d]",
			msgbin.CmdLen, recvmsg.Testno)
		return
	case servercomm.STimeTickCommandID:
		recvmsg := &servercomm.STimeTickCommand{}
		recvmsg.ReadBinary([]byte(msgbin.ProtoData))
		return
	case servercomm.SLoginRetCommandID:
		this.connectMutex.Lock()
		defer this.connectMutex.Unlock()
		// 收到登陆服务器返回的消息
		recvmsg := &servercomm.SLoginRetCommand{}
		recvmsg.ReadBinary([]byte(msgbin.ProtoData))
		if recvmsg.Loginfailed > 0 {
			conn.Terminate()
			if recvmsg.Loginfailed == servercomm.LOGINRETCODE_IDENTICAL {
				conn.IsNormalDisconnect = true
				this.Debug("[SubnetManager.msgParseTCPConn] " +
					"重复连接,不必连接")
			} else {
				this.Error("[SubnetManager.msgParseTCPConn] " +
					"连接验证失败,断开连接")
			}
			return
		}
		conn.Serverinfo = recvmsg.Destination
		this.Debug("[SubnetManager.msgParseTCPConn] "+
			"连接服务器验证成功,id:%s,ipport:%s",
			conn.Serverinfo.ServerID, conn.Serverinfo.ServerAddr)
		return
	case servercomm.SLoginCommandID:
		recvmsg := &servercomm.SLoginCommand{}
		recvmsg.ReadBinary([]byte(msgbin.ProtoData))
		this.OnServerLogin(conn, recvmsg)
		return
	case servercomm.SLogoutCommandID:
		// 服务器已主动关闭，不再尝试连接它了
		conn.IsNormalDisconnect = true
		this.connectMutex.Lock()
		defer this.connectMutex.Unlock()
		this.connInfos.RemoveConnInfo(conn.Serverinfo.ServerID)
		this.Debug("[msgParseTCPConn] 服务器已主动关闭，不再尝试连接它了 "+
			"ServerInfo[%s]", conn.Serverinfo.GetJson())
		return
	case servercomm.SNotifyAllInfoID:
		// 收到所有服务器的配置信息
		recvmsg := &servercomm.SNotifyAllInfo{}
		recvmsg.ReadBinary([]byte(msgbin.ProtoData))
		this.connectMutex.Lock()
		defer this.connectMutex.Unlock()
		this.Debug("[SubnetManager.msgParseTCPConn] " +
			"收到所有服务器列表信息")
		// 所有服务器信息列表
		for i := 0; i < len(recvmsg.Serverinfos); i++ {
			serverinfo := recvmsg.Serverinfos[i]
			this.connInfos.AddConnInfo(serverinfo)
		}
		return
	}
	msgqueues := &ConnectMsgQueueStruct{}
	msgqueues.conn = conn
	msgqueues.msg = msgbin
	this.MultiQueueControl(msgqueues)
}

// 分配消息处理线程
func (this *SubnetManager) MultiQueueControl(
	msgqueues *ConnectMsgQueueStruct) {
	if this.maxRunningMsgNum < 1 {
		this.OnRecvTCPMsg(msgqueues.conn, msgqueues.msg)
		return
	}
	who := this.OnGetRecvTCPMsgParseChan(msgqueues.conn,
		this.maxRunningMsgNum, msgqueues.msg)
	this.runningMsgChan[who] <- msgqueues
}

func (this *SubnetManager) InitMsgQueue(sum int32) {
	// 最大同时处理的消息数量
	this.maxRunningMsgNum = sum // 消息队列线程数
	if this.maxRunningMsgNum < 1 {
		this.maxRunningMsgNum = 1
		this.Error("[SubnetManager.InitMsgQueue] " +
			"消息处理线程数量过小，置为1...")
	}
	this.runningMsgChan = make([]chan *ConnectMsgQueueStruct,
		this.maxRunningMsgNum)
	this.Debug("[SubnetManager.InitMsgQueue] "+
		"Task 消息处理线程数量 ThreadNum[%d]", this.maxRunningMsgNum)
	i := int32(0)
	for i < this.maxRunningMsgNum {
		this.runningMsgChan[i] = make(chan *ConnectMsgQueueStruct,
			15000)
		go this.RecvmsgProcess(i)
		i++
	}
}

// 并行处理接收消息队列数据
func (this *SubnetManager) MultiRecvmsgQueue(
	index int32) (normalreturn bool) {
	if this.runningMsgChan == nil || this.runningMsgChan[index] == nil {
		return true
	}
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := util.GetPanicInfo(recover()); err != nil {
			// 这里的err其实就是panic传入的内容
			this.Error("[SubnetManager.MultiRecvmsgQueue] "+
				"Panic: ErrName[%v] \n Stack[%s]", err, stackInfo)
			normalreturn = false
		}
	}()
	normalreturn = true

	msgchan := this.runningMsgChan[index]
	for msgqueues := range msgchan {
		functiontime := util.FunctionTime{}
		functiontime.Start("MultiRecvmsgQueue")
		this.OnRecvTCPMsg(msgqueues.conn, msgqueues.msg)
		functiontime.Stop()
	}
	return true
}

func (this *SubnetManager) RecvmsgProcess(index int32) {
	for {
		if this.MultiRecvmsgQueue(index) {
			// 正常退出
			break
		}
	}
}
