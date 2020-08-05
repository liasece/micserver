package subnet

import (
	"time"

	"github.com/liasece/micserver/connect"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/servercomm"
	"github.com/liasece/micserver/util/hash"
	"github.com/liasece/micserver/util/monitor"
	"github.com/liasece/micserver/util/sysutil"
)

// ConnectMsgQueueStruct 服务器消息处理封包
type ConnectMsgQueueStruct struct {
	conn *connect.Server
	msg  *msg.MessageBinary
}

// onConnectClose 当TCP连接被移除时调用
func (manager *Manager) onConnectClose(conn *connect.Server) {
	manager.RemoveServer(conn.GetTempID())
}

// OnRecvTCPMsg 当收到TCP消息时调用
func (manager *Manager) OnRecvTCPMsg(conn *connect.Server,
	msgbinary *msg.MessageBinary) {
	if manager.subnetHook != nil {
		manager.subnetHook.OnRecvSubnetMsg(conn, msgbinary)
	} else {
		manager.Syslog("manager.SubnetCallback.fonRecvMsg == nil", log.Uint16("MsgID", msgbinary.GetMsgID()))
	}
}

// getRecvTCPMsgParseChan 获取TCP消息的消息处理通道
func (manager *Manager) getRecvTCPMsgParseChan(conn *connect.Server,
	maxChan int32, msgbinary *msg.MessageBinary) int32 {
	chankey := ""
	if conn.ModuleInfo != nil {
		chankey = conn.ModuleInfo.ModuleID
	}
	switch msgbinary.GetMsgID() {
	case servercomm.SForwardToClientID:
		layerMsg := &servercomm.SForwardToClient{}
		layerMsg.ReadBinary(msgbinary.ProtoData)
		msgbinary.SetObj(layerMsg)
		chankey = layerMsg.ToClientID
	}
	if chankey != "" {
		hash := int32(hash.GetStringHash(chankey))
		if hash < 0 {
			hash = -hash
		}
		return hash % maxChan
	}
	return 0
}

// OnCreateNewServer 当新增一个服务器连接时调用
func (manager *Manager) OnCreateNewServer(conn *connect.Server) {
}

// onConnectRecv 当收到了一个服务器消息时调用
func (manager *Manager) onConnectRecv(conn *connect.Server,
	msgbin *msg.MessageBinary) {
	if conn.GetSCType() == connect.ServerSCTypeTask {
		curtime := time.Now().Unix()
		if conn.IsTerminateTimeout(curtime) {
			manager.onClientDisconnected(conn)
			manager.Error("[Manager.onConnectRecv] Prolonged failure to verify, disconnect", log.String("TmpID[%s]", conn.GetTempID()))
			return
		}
		if conn.IsTerminateForce() {
			manager.onClientDisconnected(conn)
			manager.Syslog("[Manager.onConnectRecv] The server is actively disconnected", log.String("TmpID[%s]", conn.GetTempID()))
			return
		}
	}
	switch msgbin.GetMsgID() {
	case servercomm.STestCommandID:
		recvmsg := &servercomm.STestCommand{}
		recvmsg.ReadBinary([]byte(msgbin.ProtoData))
		manager.Syslog("[Manager.onConnectRecv] Server test message received", log.Int("MsgLen", msgbin.GetTotalLength()), log.Uint32("No", recvmsg.Testno))
		return
	case servercomm.STimeTickCommandID:
		recvmsg := &servercomm.STimeTickCommand{}
		recvmsg.ReadBinary([]byte(msgbin.ProtoData))
		return
	case servercomm.SLoginRetCommandID:
		manager.connectMutex.Lock()
		defer manager.connectMutex.Unlock()
		// 收到登陆服务器返回的消息
		recvmsg := &servercomm.SLoginRetCommand{}
		recvmsg.ReadBinary([]byte(msgbin.ProtoData))
		if recvmsg.Loginfailed > 0 {
			conn.Terminate()
			if recvmsg.Loginfailed == servercomm.LoginRetCodeIdentical {
				conn.IsNormalDisconnect = true
				manager.Syslog("[Manager.onConnectRecv] Duplicate connection, no need to connect", log.String("TmpID", conn.GetTempID()))
			} else {
				manager.Error("[Manager.onConnectRecv] Connection verification failed, disconnected", log.String("TmpID", conn.GetTempID()))
			}
			return
		}
		conn.ModuleInfo = recvmsg.Destination
		manager.Syslog("[Manager.onConnectRecv] Connection server verification successful",
			log.String("ModuleID", conn.ModuleInfo.ModuleID), log.String("ipport", conn.ModuleInfo.ModuleAddr))
		manager.subnetHook.OnServerJoinSubnet(conn)
		return
	case servercomm.SLoginCommandID:
		recvmsg := &servercomm.SLoginCommand{}
		recvmsg.ReadBinary([]byte(msgbin.ProtoData))
		manager.OnServerLogin(conn, recvmsg)
		return
	case servercomm.SLogoutCommandID:
		// 服务器已主动关闭，不再尝试连接它了
		conn.IsNormalDisconnect = true
		manager.connectMutex.Lock()
		defer manager.connectMutex.Unlock()
		manager.connInfos.Delete(conn.ModuleInfo.ModuleID)
		manager.Syslog("[Manager.onConnectRecv] The server is actively shutting down. No more attempts to connect to it", log.Reflect("ModuleInfo", conn.ModuleInfo))
		return
	case servercomm.SNotifyAllInfoID:
		// 收到所有服务器的配置信息
		recvmsg := &servercomm.SNotifyAllInfo{}
		recvmsg.ReadBinary([]byte(msgbin.ProtoData))
		manager.connectMutex.Lock()
		defer manager.connectMutex.Unlock()
		manager.Syslog("[Manager.onConnectRecv] Receive all server listings")
		// 所有服务器信息列表
		for i := 0; i < len(recvmsg.ServerInfos); i++ {
			serverinfo := recvmsg.ServerInfos[i]
			manager.connInfos.Add(serverinfo)
		}
		return
	}
	msgqueues := &ConnectMsgQueueStruct{}
	msgqueues.conn = conn
	msgqueues.msg = msgbin
	manager.MultiQueueControl(msgqueues)
}

// MultiQueueControl 分配消息处理线程
func (manager *Manager) MultiQueueControl(
	msgqueues *ConnectMsgQueueStruct) {
	if manager.maxRunningMsgNum < 1 {
		manager.OnRecvTCPMsg(msgqueues.conn, msgqueues.msg)
		return
	}
	who := manager.getRecvTCPMsgParseChan(msgqueues.conn,
		manager.maxRunningMsgNum, msgqueues.msg)
	if who >= int32(len(manager.runningMsgChan)) || who < 0 {
		log.Panic("who >= len(manager.runningMsgChan)", log.Int32("who", who), log.Int("len", len(manager.runningMsgChan)))
	}
	manager.runningMsgChan[who] <- msgqueues
}

// InitMsgQueue 初始化消息处理队列
func (manager *Manager) InitMsgQueue(sum int32) {
	// 最大同时处理的消息数量
	manager.maxRunningMsgNum = sum // 消息队列线程数
	if manager.maxRunningMsgNum < 1 {
		manager.maxRunningMsgNum = 1
		manager.Warn("[Manager.InitMsgQueue] The number of message processing threads is too small, set to 1...")
	}
	manager.runningMsgChan = make([]chan *ConnectMsgQueueStruct,
		manager.maxRunningMsgNum)
	manager.Syslog("[Manager.InitMsgQueue] Task number of message processing threads", log.Int32("ThreadNum", manager.maxRunningMsgNum))
	for i := int32(0); i < manager.maxRunningMsgNum; i++ {
		manager.runningMsgChan[i] = make(chan *ConnectMsgQueueStruct,
			15000)
		go manager.RecvmsgProcess(i)
	}
}

// MultiRecvmsgQueue 并行处理接收消息队列数据
func (manager *Manager) MultiRecvmsgQueue(
	index int32) (normalreturn bool) {
	if manager.runningMsgChan == nil || manager.runningMsgChan[index] == nil {
		log.Panic("manager.runningMsgChan[%d] == nil", index)
	}
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if stackInfo, err := sysutil.GetPanicInfo(recover()); err != nil {
			// 这里的err其实就是panic传入的内容
			manager.Error("[Manager.MultiRecvmsgQueue] Panic", log.ErrorField(err), log.String("Stack", stackInfo))
			normalreturn = false
		}
	}()
	normalreturn = true

	msgchan := manager.runningMsgChan[index]
	for msgqueues := range msgchan {
		functiontime := monitor.FunctionTime{}
		functiontime.Start("MultiRecvmsgQueue")
		manager.OnRecvTCPMsg(msgqueues.conn, msgqueues.msg)
		functiontime.Stop()
	}
	return true
}

// RecvmsgProcess 保持服务器消息处理线程
func (manager *Manager) RecvmsgProcess(index int32) {
	for {
		if manager.MultiRecvmsgQueue(index) {
			// 正常退出
			break
		}
	}
}
