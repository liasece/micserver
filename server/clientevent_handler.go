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

// HookGate 设置网关事件监听者，如果本服务没有启用网关，将不会收到任何事件
func (handler *clientEventHandler) HookGate(gateHook base.GateHook) {
	handler.gateHook = gateHook
}

// OnAcceptClientConnect 当接收一个客户端的TCP连接时调用
func (handler *clientEventHandler) OnAcceptClientConnect(conn net.Conn) {
	if handler.gateHook != nil {
		handler.gateHook.OnAcceptClientConnect(conn)
	}
}

// OnNewClient 当新增了一个客户端连接时调用
func (handler *clientEventHandler) OnNewClient(client *connect.Client) {
	moduleType := util.GetModuleIDType(handler.server.moduleID)
	client.SetBind(moduleType, handler.server.moduleID)

	if handler.gateHook != nil {
		handler.gateHook.OnNewClient(client)
	}
}

// OnCloseClient 当关闭了一个客户端连接时调用
func (handler *clientEventHandler) OnCloseClient(client *connect.Client) {
	if handler.gateHook != nil {
		handler.gateHook.OnCloseClient(client)
	}
}

// OnRecvClientMsg 当收到了一个客户端消息时调用
func (handler *clientEventHandler) OnRecvClientMsg(
	client *connect.Client, msgbin *msg.MessageBinary) {
	if handler.gateHook != nil {
		handler.gateHook.OnRecvClientMsg(client, msgbin)
	}
}
