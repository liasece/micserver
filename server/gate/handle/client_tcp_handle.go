package handle

import (
	"github.com/liasece/micserver/connect"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/servercomm"
	"time"
)

type TFuncHandleSocketPackage func(*connect.ClientConn, *msg.MessageBinary)

type ClientTcpHandler struct {
	*log.Logger

	Analysiswsmsgcount     uint32
	regHandleSocketPackage TFuncHandleSocketPackage
}

func (this *ClientTcpHandler) RegHandleSocketPackage(
	cb TFuncHandleSocketPackage) {
	this.regHandleSocketPackage = cb
}

func (this *ClientTcpHandler) OnRecvSocketPackage(task *connect.ClientConn,
	msgbin *msg.MessageBinary) {
	cmdname := servercomm.MsgIdToString(msgbin.CmdID)
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
