package msg

import (
	"github.com/liasece/micserver/util/buffer"
)

// IMsgStruct micserver消息接口，所有在micserver架构间的消息类都要满足这些接口
type IMsgStruct interface {
	WriteBinary(data []byte) int
	GetMsgId() uint16
	GetMsgName() string
	GetSize() int
	GetJSON() string
}

// IMsgCodec 消息编解码器
type IMsgCodec interface {
	RangeMsgBinary(buf *buffer.IOBuffer, cb func(*MessageBinary)) error
	EncodeBytes(cmdid uint16, protodata []byte) *MessageBinary
	EncodeObj(v IMsgStruct) *MessageBinary
}
