package msg

import (
	"github.com/liasece/micserver/util/buffer"
)

// MicServer消息接口，所有在micserver架构间的消息类都要满足这些接口
type MsgStruct interface {
	WriteBinary(data []byte) int
	GetMsgId() uint16
	GetMsgName() string
	GetSize() int
	GetJson() string
}

type IMsgCodec interface {
	RangeMsgBinary(buf *buffer.IOBuffer, cb func(*MessageBinary)) error
	EncodeBytes(cmdid uint16, protodata []byte) *MessageBinary
	EncodeObj(v MsgStruct) *MessageBinary
}
