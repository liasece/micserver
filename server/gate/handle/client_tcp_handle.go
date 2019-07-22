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
	now := time.Now().Unix()
	// 设置连接活动过期时间 5分钟
	task.SetTerminateTime(now + 5*60)

	// 如果客户端没有通过验证，后续代码不能被执行
	if !task.IsVertify() {
		task.Debug("由于没有通过验证,消息没有继续处理 MsgID[%d] MsgName[%s]",
			msgbin.CmdID, cmdname)
		return
	}

	// 以下处理需要在客户端通过验证之后才能运行
}
