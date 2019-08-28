package module

import (
	"github.com/liasece/micserver/connect"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/util"
	"net"
)

type clientEventHandler struct {
	mod *BaseModule

	regNewClient     func(client *connect.Client)
	regRecvMsg       func(client *connect.Client, msgbin *msg.MessageBinary)
	regAcceptConnect func(conn net.Conn)
}

func (this *clientEventHandler) RegNewClient(
	cb func(client *connect.Client)) {
	this.regNewClient = cb
}

func (this *clientEventHandler) onNewClient(client *connect.Client) {
	servertype := util.GetServerIDType(this.mod.ModuleID)
	client.SetBindServer(servertype, this.mod.ModuleID)

	if this.regNewClient != nil {
		this.regNewClient(client)
	}
}

func (this *clientEventHandler) RegRecvMsg(
	cb func(client *connect.Client, msgbin *msg.MessageBinary)) {
	this.regRecvMsg = cb
}

func (this *clientEventHandler) onRecvMsg(
	client *connect.Client, msgbin *msg.MessageBinary) {
	if this.regRecvMsg != nil {
		this.regRecvMsg(client, msgbin)
	}
}

func (this *clientEventHandler) RegAcceptConnect(cb func(conn net.Conn)) {
	this.regAcceptConnect = cb
}

func (this *clientEventHandler) onAcceptConnect(conn net.Conn) {
	if this.regAcceptConnect != nil {
		this.regAcceptConnect(conn)
	}
}
