package module

import (
	"github.com/liasece/micserver/connect"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/util"
)

type clientEventHandler struct {
	mod *BaseModule

	regRecvMsg func(client *connect.ClientConn, msgbin *msg.MessageBinary)
}

func (this *clientEventHandler) OnNewClient(client *connect.ClientConn) {
	servertype := util.GetServerIDType(this.mod.ModuleID)
	client.SetBindServer(servertype, this.mod.ModuleID)
}

func (this *clientEventHandler) RegRecvMsg(
	cb func(client *connect.ClientConn, msgbin *msg.MessageBinary)) {
	this.regRecvMsg = cb
}

func (this *clientEventHandler) OnRecvMsg(
	client *connect.ClientConn, msgbin *msg.MessageBinary) {
	if this.regRecvMsg != nil {
		this.regRecvMsg(client, msgbin)
	}
}
