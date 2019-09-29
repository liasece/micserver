package server

import (
	"github.com/liasece/micserver/connect"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/util"
	"net"
)

type clientEventHandler struct {
	server *Server

	fonNewClient     func(client *connect.Client)
	fonRecvClientMsg func(client *connect.Client, msgbin *msg.MessageBinary)
	fonAcceptConnect func(conn net.Conn)
}

// 接受到客户端tcp连接
func (this *clientEventHandler) RegOnAcceptConnect(cb func(conn net.Conn)) {
	this.fonAcceptConnect = cb
}

func (this *clientEventHandler) OnAcceptConnect(conn net.Conn) {
	if this.fonAcceptConnect != nil {
		this.fonAcceptConnect(conn)
	}
}

// 新的客户端连接对象
func (this *clientEventHandler) RegOnNewClient(
	cb func(client *connect.Client)) {
	this.fonNewClient = cb
}

func (this *clientEventHandler) OnNewClient(client *connect.Client) {
	servertype := util.GetServerIDType(this.server.serverid)
	client.SetBindServer(servertype, this.server.serverid)

	if this.fonNewClient != nil {
		this.fonNewClient(client)
	}
}

// 收到客户端消息
func (this *clientEventHandler) RegOnRecvClientMsg(
	cb func(client *connect.Client, msgbin *msg.MessageBinary)) {
	this.fonRecvClientMsg = cb
}

func (this *clientEventHandler) OnRecvClientMsg(
	client *connect.Client, msgbin *msg.MessageBinary) {
	if this.fonRecvClientMsg != nil {
		this.fonRecvClientMsg(client, msgbin)
	}
}
