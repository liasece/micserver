package base

import (
	"github.com/liasece/micserver/connect"
	"github.com/liasece/micserver/msg"
)

type SubnetHook interface {
	// 收到子网消息
	OnRecvSubnetMsg(server *connect.Server, msgbin *msg.MessageBinary)
}
