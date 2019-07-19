package handle

import (
	"github.com/liasece/micserver/comm"
	"github.com/liasece/micserver/conf"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/server/gate/manager"
	"github.com/liasece/micserver/tcpconn"
	"github.com/liasece/micserver/util"
	"io"
	"net"
	"time"
)

type ClientTcpHandler struct {
	clientSocketManager *manager.ClientSocketManager
	moduleConfig        *conf.TopConfig
}

func (this *ClientTcpHandler) Init(clientSocketManager *manager.ClientSocketManager,
	config *conf.TopConfig) {
	this.clientSocketManager = clientSocketManager
	this.moduleConfig = config
}

// socket消息处理
func (this *ClientTcpHandler) clientsocketJsonHandle(conn net.Conn) {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := util.GetPanicInfo(recover()); err != nil {
			log.Error("[clientsocketJsonHandle] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
		}
	}()
	log.Debug("[clientsocketJsonHandle] Receive one conn connect json")
	task, err := this.clientSocketManager.AddClientTcpSocket(conn)
	if err != nil || task == nil {
		log.Error("[clientsocketJsonHandle] "+
			"创建 ClientTcpSocket 对象失败，断开连接 Err[%s]", err.Error())
		return
	}
	netbuffer := util.NewIOBuffer(conn, 64*1024)
	msgReader := msg.NewMessageBinaryReader(netbuffer)

	// 所有连接都需要经过加密
	// task.Encryption = base.EncryptionTypeXORSimple

	for {
		if !task.Check() {
			// 强制移除客户端连接
			// manager.NotifyClientUserOffline(task)
			this.clientSocketManager.RemoveTaskByTmpID(task.Tempid)
			return
		}
		// 设置阻塞读取过期时间
		err := conn.SetReadDeadline(
			time.Now().Add(time.Duration(time.Millisecond * 250)))
		if err != nil {
			task.Error("[clientsocketJsonHandle] SetReadDeadline Err[%s]",
				err.Error())
		}
		// buffer从连接中读取socket数据
		_, err = netbuffer.ReadFromReader()

		// 异常
		if err != nil {
			if err == io.EOF {
				task.Debug("[clientsocketJsonHandle] "+
					"Scoket数据读写异常,断开连接了,"+
					"scoket返回 Err[%s]", err.Error())
				// manager.NotifyClientUserOffline(task)
				this.clientSocketManager.RemoveTaskByTmpID(task.Tempid)
				return
			} else {
				continue
			}
		}

		err = msgReader.RangeMsgBinary(func(msgbinary *msg.MessageBinary) {
			if task.Encryption != msg.EncryptionTypeNone &&
				msgbinary.CmdMask != task.Encryption {
				task.Error("加密方式错误，加密方式应为 %d 此消息为 %d "+
					"MsgID[%d]", task.Encryption,
					msgbinary.CmdMask, msgbinary.CmdID)
			} else {
				// 解析消息
				this.ParseClientJsonMsg(msgbinary, task)
			}
		})
		if err != nil {
			task.Error("[clientsocketJsonHandle] 解析消息错误，断开连接 "+
				"Err[%s]", err.Error())
			// 强制移除客户端连接
			// manager.NotifyClientUserOffline(task)
			this.clientSocketManager.RemoveTaskByTmpID(task.Tempid)
			return
		}
	}
}

func (this *ClientTcpHandler) ParseClientJsonMsg(msgbin *msg.MessageBinary,
	task *tcpconn.ClientConn) {
	cmdname := comm.MsgIdToString(msgbin.CmdID)
	this.clientSocketManager.Analysiswsmsgcount++
	defer msgbin.Free()
	// task.Debug("收到数据 msgname:[%s] openid [%s] <%s>", cmdname,
	// 	task.Openid, hex.EncodeToString(msgbin.ProtoData))
	task.Debug("[ParseClientJsonMsg] 收到数据 "+
		"MsgMask[%d] MsgID[%d] Msgname[%s] CmdLen[%d] DataLen[%d]",
		msgbin.CmdMask, msgbin.CmdID, cmdname, msgbin.CmdLen, msgbin.DataLen)
	if msgbin.CmdID == 0 {
		task.Error("[ParseClientJsonMsg] 错误的 MsgID[%d]", msgbin.CmdID)
		return
	}

	// 接收到有效消息，开始处理
	now := uint64(time.Now().Unix())
	// 设置连接活动过期时间 5分钟
	task.SetTerminateTime(now + 5*60)

	// if msgbin.CmdID == comm.UStartLoginUserID {
	// 	task.Encryption = msgbin.CmdMask

	// 	recvmsg := &comm.UStartLoginUser{}
	// 	recvmsg.ReadBinary(msgbin.ProtoData)
	// 	// manager.OnMsg_UStartLoginUser(task, recvmsg)
	// 	return
	// } else if msgbin.CmdID == comm.GPingGatewayID {
	// 	recvmsg := &comm.GPingGateway{}
	// 	recvmsg.ReadBinary(msgbin.ProtoData)
	// 	task.Debug("收到Ping Msg[%d]", recvmsg.GetJson())
	// 	// 处理ping消息
	// 	syn, ack, seq := task.GetPing().OnRecv(
	// 		recvmsg.SYN, recvmsg.ACK, recvmsg.Seq)
	// 	if syn != 0 {
	// 		sendmsg := &comm.GPingGateway{}
	// 		sendmsg.SYN = syn
	// 		sendmsg.ACK = ack
	// 		sendmsg.Seq = seq
	// 		task.SendCmd(sendmsg)
	// 	}
	// 	return
	// }

	// 如果客户端没有通过验证，后续代码不能被执行
	if !task.IsVertify() {
		task.Debug("由于没有通过验证,消息没有继续处理 MsgID[%d] MsgName[%s]",
			msgbin.CmdID, cmdname)
		return
	}

	// 以下处理需要在客户端通过验证之后才能运行

	// log.Debug()

	// if msgbin.CmdID == comm.GBroadcastToUserID {
	// 	// 用户间消息直接转发
	// 	recvmsg := &comm.GBroadcastToUser{}
	// 	recvmsg.ReadBinary(msgbin.ProtoData)
	// 	forward.ForwardClientCmdToUsers(task, recvmsg.RecipientUUIDs,
	// 		msgbin.CmdID, msgbin.ProtoData)
	// 	task.Debug("客户端间的消息转发 MsgID[%d] MsgName[%s] Recv[%v]",
	// 		recvmsg.GetMsgId(), recvmsg.GetMsgName(), recvmsg.RecipientUUIDs)
	// } else {
	// 	// 用户 - 服务器 间消息转发
	// 	msgtype := comm.MsgIdToType(msgbin.CmdID)

	// 	// 判断消息是发送至RoomServer的
	// 	if msgtype == 'R' {
	// 		// RoomServer
	// 		forward.ForwardClientCmdToRoomServer(task.Roomserverid, msgbin.CmdID,
	// 			msgbin.ProtoData, task)
	// 	} else if msgtype == 'U' {
	// 		// UserServer
	// 		forward.ForwardClientCmdToUserServer(task.Userserverid, msgbin.CmdID,
	// 			msgbin.ProtoData, task)
	// 	} else {
	// 		forward.ForwardClientCmdToUserServer(task.Userserverid, msgbin.CmdID,
	// 			msgbin.ProtoData, task)
	// 		task.Error("未知的消息类型,默认转发到UserServer MsgType[%v] "+
	// 			"MsgID[%d] MsgName[%s]",
	// 			msgtype, msgbin.CmdID, cmdname)
	// 	}
	// }
}

func (this *ClientTcpHandler) StartAddClientTcpSocketHandle(addr string) {
	// 由于部分 NAT 主机没有网卡概念，需要自己配置IP
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Error("[ClientSocket] %s", err.Error())
		return
	}
	log.Debug("Gateway Client TCP服务启动成功 IPPort[%s]", addr)
	go func() {
		for {
			// 接受连接
			conn, err := ln.Accept()
			if err != nil {
				// handle error
				log.Error("[StartAddClientTcpSocketHandle] Accept() ERR:%q",
					err.Error())
				continue
			}
			go this.clientsocketJsonHandle(conn)
		}
	}()
}
