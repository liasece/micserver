/*
Package base 网关的基本接口
*/
package base

import (
	"net"

	"github.com/liasece/micserver/connect"
	"github.com/liasece/micserver/msg"
)

// GateHook 上层网关服务需要实现的处理网关事件的接口
type GateHook interface {
	// 接受到客户端tcp连接
	OnAcceptClientConnect(conn net.Conn)
	// 新的客户端连接对象
	OnNewClient(client *connect.Client)
	// 关闭客户端连接对象
	OnCloseClient(client *connect.Client)
	// 收到客户端消息
	OnRecvClientMsg(client *connect.Client, msgbin *msg.MessageBinary)
}
