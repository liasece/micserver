/*
Package base 服务器子网基础
*/
package base

import (
	"github.com/liasece/micserver/connect"
	"github.com/liasece/micserver/msg"
)

// SubnetHook 服务器子网管理器需要实现的接口
type SubnetHook interface {
	// 当一个服务器成功加入网络时调用
	OnServerJoinSubnet(server *connect.Server)
	// 收到子网消息
	OnRecvSubnetMsg(server *connect.Server, msgbin *msg.MessageBinary)
}
