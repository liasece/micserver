package network

import (
	"github.com/liasece/micserver/msg"
)

type IConnection interface {
	IsAlive() bool
	Shutdown() error
	Read(toData []byte) (int, error)
	StartRecv()
	SendMessageBinary(msgbinary *msg.MessageBinary) error
	SendBytes(cmdid uint16, protodata []byte) error
}
