package base

import (
	"github.com/liasece/micserver/connect"
	"github.com/liasece/micserver/msg"
	"net"
)

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
