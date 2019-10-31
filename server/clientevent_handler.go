package server

import (
	"net"

	"github.com/liasece/micserver/connect"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/server/gate/base"
	"github.com/liasece/micserver/util"
)

type clientEventHandler struct {
	server *Server

	gateHook base.GateHook
}

func (this *clientEventHandler) HookGate(gateHook base.GateHook) {
	this.gateHook = gateHook
}

func (this *clientEventHandler) OnAcceptClientConnect(conn net.Conn) {
	if this.gateHook != nil {
		this.gateHook.OnAcceptClientConnect(conn)
	}
}

func (this *clientEventHandler) OnNewClient(client *connect.Client) {
	servertype := util.GetServerIDType(this.server.serverid)
	client.SetBindServer(servertype, this.server.serverid)

	if this.gateHook != nil {
		this.gateHook.OnNewClient(client)
	}
}

func (this *clientEventHandler) OnRecvClientMsg(
	client *connect.Client, msgbin *msg.MessageBinary) {
	if this.gateHook != nil {
		this.gateHook.OnRecvClientMsg(client, msgbin)
	}
}
