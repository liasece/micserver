package connect

import (
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/network/baseio"
)

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
}

type ConnectHook interface {
	OnRecvConnectMessage(*Client, *msg.MessageBinary)
	OnConnectClose(*Client)
}
