package server

import (
	"net"

	"github.com/liasece/micserver/connect"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/server/gate/base"
	"github.com/liasece/micserver/util"
)

// 客户端事件处理
type clientEventHandler struct {
	server *Server

	gateHook base.GateHook
}

// 设置网关事件监听者，如果本服务没有启用网关，将不会收到任何事件
func (this *clientEventHandler) HookGate(gateHook base.GateHook) {
	this.gateHook = gateHook
}

// 当接收一个客户端的TCP连接时调用
func (this *clientEventHandler) OnAcceptClientConnect(conn net.Conn) {
	if this.gateHook != nil {
		this.gateHook.OnAcceptClientConnect(conn)
	}
}

// 当新增了一个客户端连接时调用
func (this *clientEventHandler) OnNewClient(client *connect.Client) {
	moduleType := util.GetModuleIDType(this.server.moduleid)
	client.SetBind(moduleType, this.server.moduleid)

	if this.gateHook != nil {
		this.gateHook.OnNewClient(client)
	}
}

// 当关闭了一个客户端连接时调用
func (this *clientEventHandler) OnCloseClient(client *connect.Client) {
	if this.gateHook != nil {
		this.gateHook.OnCloseClient(client)
	}
}

// 当收到了一个客户端消息时调用
func (this *clientEventHandler) OnRecvClientMsg(
	client *connect.Client, msgbin *msg.MessageBinary) {
	if this.gateHook != nil {
		this.gateHook.OnRecvClientMsg(client, msgbin)
	}
}
