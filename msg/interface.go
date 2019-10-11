package msg

// MicServer消息接口，所有在micserver架构间的消息类都要满足这些接口
type MsgStruct interface {
	WriteBinary(data []byte) int
	GetMsgId() uint16
	GetMsgName() string
	GetSize() int
	GetJson() string
}
