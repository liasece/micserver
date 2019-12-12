package base

import (
	"github.com/liasece/micserver/connect"
	"github.com/liasece/micserver/msg"
)

type SubnetHook interface {
	// 当一个服务器成功加入网络时调用
	OnServerJoinSubnet(server *connect.Server)
	// 收到子网消息
	OnRecvSubnetMsg(server *connect.Server, msgbin *msg.MessageBinary)
}
