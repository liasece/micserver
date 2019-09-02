package msg

// MicServer消息接口，所有在micserver架构间的消息类都要满足这些接口
type MsgStruct interface {
	WriteBinary(data []byte) int
	GetMsgId() uint16
	GetMsgName() string
	GetSize() int
	GetJson() string
}

type SendCompletedAgent struct {
	F    func(interface{}) // 用于优化的临时数据指针请注意使用！
	Argv interface{}       // 用于优化的临时数据指针请注意使用！
}
