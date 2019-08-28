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
	regRecvMsg         func(*connect.Client, *msg.MessageBinary)
	regNewClient       func(*connect.Client)
	regAcceptConnect   func(net.Conn)
}

func (this *ClientTcpHandler) RegRecvMsg(
	cb func(*connect.Client, *msg.MessageBinary)) {
	this.regRecvMsg = cb
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

	if this.regRecvMsg != nil {
		this.regRecvMsg(client, msgbin)
	}
}

func (this *ClientTcpHandler) RegNewClient(
	cb func(*connect.Client)) {
	this.regNewClient = cb
}

func (this *ClientTcpHandler) OnNewClient(client *connect.Client) {
	if this.regNewClient != nil {
		this.regNewClient(client)
	}
}

func (this *ClientTcpHandler) RegAcceptConnect(
	cb func(net.Conn)) {
	this.regAcceptConnect = cb
}

func (this *ClientTcpHandler) OnAcceptConnect(conn net.Conn) {
	if this.regAcceptConnect != nil {
		this.regAcceptConnect(conn)
	}
}
