package handle

import (
	"github.com/liasece/micserver/comm"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/tcpconn"
	"time"
)

type ClientTcpHandler struct {
	*log.Logger

	Analysiswsmsgcount uint32
}

func (this *ClientTcpHandler) OnRecvSocketPackage(msgbin *msg.MessageBinary,
	task *tcpconn.ClientConn) {
	cmdname := comm.MsgIdToString(msgbin.CmdID)
	this.Analysiswsmsgcount++
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

	// this.Debug()

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
