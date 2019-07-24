package handle

import (
	"github.com/liasece/micserver/comm"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/tcpconn"
	"time"
)

type TFuncHandleSocketPackage func(*tcpconn.ClientConn, *msg.MessageBinary)

type ClientTcpHandler struct {
	*log.Logger

	Analysiswsmsgcount     uint32
	regHandleSocketPackage TFuncHandleSocketPackage
}

func (this *ClientTcpHandler) RegHandleSocketPackage(
	cb TFuncHandleSocketPackage) {
	this.regHandleSocketPackage = cb
}

func (this *ClientTcpHandler) OnRecvSocketPackage(task *tcpconn.ClientConn,
	msgbin *msg.MessageBinary) {
	cmdname := comm.MsgIdToString(msgbin.CmdID)
	this.Analysiswsmsgcount++
	defer msgbin.Free()

	task.Debug("[ParseClientJsonMsg] 收到数据 "+
		"MsgMask[%d] MsgID[%d] Msgname[%s] CmdLen[%d] DataLen[%d]",
		msgbin.CmdMask, msgbin.CmdID, cmdname, msgbin.CmdLen, msgbin.DataLen)
	if msgbin.CmdID == 0 {
		task.Error("[ParseClientJsonMsg] 错误的 MsgID[%d]", msgbin.CmdID)
		return
	}

	// 接收到有效消息，开始处理
	now := time.Now().Unix()
	// 设置连接活动过期时间 5分钟
	task.SetTerminateTime(now + 5*60)

	if this.regHandleSocketPackage != nil {
		this.regHandleSocketPackage(task, msgbin)
	}
}
