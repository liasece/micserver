package handle

import (
	"net"

	"github.com/liasece/micserver/connect"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/server/gate/base"
	"github.com/liasece/micserver/servercomm"
)

type ClientTcpHandler struct {
	*log.Logger

	Analysiswsmsgcount uint32

	gateHook base.GateHook
}

func (this *ClientTcpHandler) HookGate(gateHook base.GateHook) {
	this.gateHook = gateHook
}

func (this *ClientTcpHandler) OnRecvMessage(client *connect.Client,
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

	if this.gateHook != nil {
		this.gateHook.OnRecvClientMsg(client, msgbin)
	}
}

func (this *ClientTcpHandler) OnNewClient(client *connect.Client) {
	if this.gateHook != nil {
		this.gateHook.OnNewClient(client)
	}
}

func (this *ClientTcpHandler) OnClose(client *connect.Client) {
	if this.gateHook != nil {
		this.gateHook.OnCloseClient(client)
	}
}

func (this *ClientTcpHandler) OnAcceptClientConnect(conn net.Conn) {
	if this.gateHook != nil {
		this.gateHook.OnAcceptClientConnect(conn)
	}
}
