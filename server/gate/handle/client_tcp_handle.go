package handle

import (
	"github.com/liasece/micserver/connect"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/servercomm"
	"net"
	"time"
)

type ClientTcpHandler struct {
	*log.Logger

	Analysiswsmsgcount uint32
	fonRecvMsg         func(*connect.Client, *msg.MessageBinary)
	fonNewClient       func(*connect.Client)
	fonAcceptConnect   func(net.Conn)
}

func (this *ClientTcpHandler) RegOnRecvMsg(
	cb func(*connect.Client, *msg.MessageBinary)) {
	this.fonRecvMsg = cb
}

func (this *ClientTcpHandler) OnConnectRecv(client *connect.Client,
	msgbin *msg.MessageBinary) {
	cmdname := servercomm.MsgIdToString(msgbin.CmdID)
	this.Analysiswsmsgcount++
	defer msgbin.Free()

	if !client.Check() {
		client.Shutdown()
		return
	}
	client.Debug("[ParseClientJsonMsg] 收到数据 "+
		"MsgID[%d] Msgname[%s] CmdLen[%d] DataLen[%d]",
		msgbin.CmdID, cmdname, msgbin.CmdLen, msgbin.DataLen)
	// 接收到有效消息，开始处理
	now := time.Now().Unix()
	// 设置连接活动过期时间 5分钟
	client.SetTerminateTime(now + 5*60)

	if this.fonRecvMsg != nil {
		this.fonRecvMsg(client, msgbin)
	}
}

func (this *ClientTcpHandler) RegOnNewClient(
	cb func(*connect.Client)) {
	this.fonNewClient = cb
}

func (this *ClientTcpHandler) OnNewClient(client *connect.Client) {
	if this.fonNewClient != nil {
		this.fonNewClient(client)
	}
}

func (this *ClientTcpHandler) RegOnAcceptConnect(
	cb func(net.Conn)) {
	this.fonAcceptConnect = cb
}

func (this *ClientTcpHandler) OnAcceptConnect(conn net.Conn) {
	if this.fonAcceptConnect != nil {
		this.fonAcceptConnect(conn)
	}
}
