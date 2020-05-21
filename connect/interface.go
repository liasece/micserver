package connect

import (
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/network/baseio"
)

// IConnection 模块间/客户端 连接都实现了该接口
type IConnection interface {
	IsAlive() bool
	Shutdown() error
	RemoteAddr() string
	Read(toData []byte) (int, error)
	StartRecv()
	GetRecvMessageChannel() chan *msg.MessageBinary
	SendMessageBinary(msgbinary *msg.MessageBinary) error
	SendBytes(cmdid uint16, protodata []byte) error
	Write(data []byte) (int, error)
	HookProtocal(p baseio.Protocal)
	SetBanAutoResize(value bool)
	GetMsgCodec() msg.IMsgCodec
	SetMsgCodec(msg.IMsgCodec)
}

// IConnectHook 连接事件钩子需要满足的接口
type IConnectHook interface {
	OnRecvConnectMessage(*Client, *msg.MessageBinary)
	OnConnectClose(*Client)
}
