package network

import (
	"github.com/liasece/micserver/msg"
	"io"
)

type IConnection interface {
	IsAlive() bool
	Shutdown() error
	RemoteAddr() string
	Read(toData []byte) (int, error)
	StartRecv()
	SendMessageBinary(msgbinary *msg.MessageBinary) error
	SendBytes(cmdid uint16, protodata []byte) error
	Write(data []byte) (int, error)
	RegDoSendTCPBytes(cb func(writer io.ReadWriter, state interface{},
		data []byte) (int, interface{}, error))
	RegDoReadTCPBytes(cb func(reader io.ReadWriter, state interface{},
		toData []byte) (int, interface{}, error))
}
