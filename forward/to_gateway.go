package forward

import (
	"github.com/liasece/micserver/msg"
	// "base/log"
	// "base/subnet"
	// "encoding/json"
	// "fmt"
	// "jsonmsg"
	// "reflect"
	// "encoding/hex"
	"comm"
	// "strconv"
	// "sync"
	// "time"
)

// 服务器发送数据接口
type TargetTask interface {
	SendCmd(base.MsgStruct) error
	SendCmdWithCallback(base.MsgStruct, func(interface{}), interface{}) error
}

// 用于将转发的消息存于对象池中，避免GC
func SendCallBack(iv interface{}) {
	if iv == nil {
		return
	}
	msg := iv.(*base.MessageBinary)
	// log.Debug("[SendCallBack] %d 回收成功", msg.CmdID)
	msg.Free()
}

// 将消息转发到 Gateway
func GatewayToWebSocket(task TargetTask, towsid uint64, v base.MsgStruct) error {

	submsg := base.MakeMessageByJson(v)

	sendmsg := &comm.SGatewayForwardCommand{}
	sendmsg.ClientConnID = towsid
	sendmsg.Cmdid = v.GetMsgId()
	sendmsg.Cmdlen = submsg.DataLen
	sendmsg.Cmddatas = submsg.ProtoData
	// log.Debug("[GatewayToWebSocket] 发送 [%s] [%+v]", v.GetMsgName(), v)
	return task.SendCmdWithCallback(sendmsg, SendCallBack, submsg)
}

// 将消息转发到 Gateway
func GatewayBroadcastToClient(task TargetTask, touuid []uint64,
	threadHash uint32, v base.MsgStruct) error {

	submsg := base.MakeMessageByJson(v)

	sendmsg := &comm.SGatewayForwardBroadcastCommand{}
	sendmsg.UUIDList = touuid
	sendmsg.Cmdid = v.GetMsgId()
	sendmsg.Cmdlen = submsg.DataLen
	sendmsg.Cmddatas = submsg.ProtoData
	sendmsg.ThreadHash = threadHash
	// log.Debug("[GatewayToWebSocket] 发送 [%s] [%+v]", v.GetMsgName(), v)
	return task.SendCmdWithCallback(sendmsg, SendCallBack, submsg)
}

// 将消息转发到 Gateway
func GatewayBroadcastBytesToClient(task TargetTask, touuid []uint64,
	threadHash uint32, msgid uint16, data []byte) error {

	submsg := base.MakeMessageByBytes(msgid, data)

	sendmsg := &comm.SGatewayForwardBroadcastCommand{}
	sendmsg.UUIDList = touuid
	sendmsg.Cmdid = submsg.CmdID
	sendmsg.Cmdlen = submsg.DataLen
	sendmsg.Cmddatas = submsg.ProtoData
	sendmsg.ThreadHash = threadHash
	// log.Debug("[GatewayToWebSocket] 发送 [%s] [%+v]", v.GetMsgName(), v)
	return task.SendCmdWithCallback(sendmsg, SendCallBack, submsg)
}
