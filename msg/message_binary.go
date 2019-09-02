package msg

import (
	"errors"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/util"
)

const (
	// 允许的最大消息大小（包括头部大小）
	// 普通代码中间应留出至少 32 字节提供给消息头部使用
	MessageMaxSize = 8 * 1024 * 1024
)

var (
	defaultHead1 MessageBinaryHeadL1
	defaultBody  MessageBinaryBody
)

// 灵活对象池的对象大小分割。
// 从对象池中获取对象时，将会从至少满足需求大小的对象的对象池中获取对象，
// 有效减少对象池中闲置对象带来的内存占用。
// 粒度控制 单位：字节
var sizeControl []int = []int{32, 64, 128, 256, 512, 1024, 2 * 1024,
	4 * 1024, 6 * 1024, 8 * 1024, 10 * 1024, 15 * 1024, 20 * 1024,
	25 * 1024, 30 * 1024, 35 * 1024, 40 * 1024, 45 * 1024, 50 * 1024,
	55 * 1024, 60 * 1024, 64 * 1024, 128 * 1024, 256 * 1024, 512 * 1024,
	1024 * 1024, 2 * 1024 * 1024, 4 * 1024 * 1024, 8 * 1024 * 1024}
var pools *util.FlexiblePool

// 初始化灵活对象池
func init() {
	pools = util.NewFlexiblePool(sizeControl, newMsgBinaryBySize)
}

// 根据 MessageBinary.buffer 的大小来创建一个对象池中的对象
func newMsgBinaryBySize(size int) interface{} {
	msg := new(MessageBinary)
	msg.buffer = make([]byte, size)
	return msg
}

// 根据消息内容大小从对象池获取对应的消息对象
func getMessageBinaryByProtoDataLength(protoDataSize int) *MessageBinary {
	totalSize := protoDataSize + MSG_HEADSIZE // 加上协议头长度
	msg, err := pools.Get(totalSize)
	if err != nil {
		log.Error("[MakeMessageByBytes] "+
			"[getMessageBinaryByProtoDataLength] CmdLen[%d] Err[%s]",
			totalSize, err.Error())
		return nil
	}
	if msg == nil {
		log.Error("[MakeMessageByBytes] "+
			"[getMessageBinaryByProtoDataLength] nil return!!! CmdLen[%d]",
			totalSize)
		return nil
	}
	return msg.(*MessageBinary)
}

type MessageBinary struct {
	MessageBinaryHeadL1
	MessageBinaryBody
	buffer []byte

	regSendDone *SendCompletedAgent // 用于优化的临时数据指针请注意使用！
}

func (this *MessageBinary) RegSendDone(cb func(interface{}), argv interface{}) {
	if cb == nil {
		return
	}
	this.regSendDone = &SendCompletedAgent{
		F:    cb,
		Argv: argv,
	}
}

func (this *MessageBinary) OnSendDone() {
	if this.regSendDone != nil {
		this.regSendDone.F(this.regSendDone.Argv)
	}
}

// 将消息对象释放到对象池中
func (this *MessageBinary) Free() {
	// 重置本消息各个属性
	this.Reset()
	// 根据缓冲区容量归类
	size := len(this.buffer)
	err := pools.Put(this, size)
	if err != nil {
		log.Error("[MessageBinary.Free] pools.Put Err[%s]",
			err.Error())
	}
}

// 重置 Message 数据
func (this *MessageBinary) Reset() {
	this.MessageBinaryHeadL1 = defaultHead1
	this.regSendDone = nil
	this.MessageBinaryBody = defaultBody
	// 为了减轻GC压力，不应重置buffer字段
}

// 从二进制流中读取 Message 结构，带消息头
func (this *MessageBinary) ReadBinary(cmddata []byte) error {
	// 获取基础数据
	offset, err := this.MessageBinaryHeadL1.ReadBinary(cmddata)
	// 过小的长度
	if err != nil {
		log.Error("[MakeMessageByBytes] "+
			"[ReadBinary] 错误的消息头L1 Err[%s]", err.Error())
		return errors.New("错误的消息头L1")
	}
	// 用于消息读取完毕之后的校验
	// 读取消息
	err = this.readBinaryNoHeadL1(cmddata[offset:])
	if err != nil {
		log.Error("[MakeMessageByBytes] "+
			"[ReadBinary] ReadBinaryNoMessageBinaryHeadL1错误 Err[%s]",
			err.Error())
		return nil
	}
	return nil
}

// 从二进制流中读取 Message 结构，无消息头
func (this *MessageBinary) readBinaryNoHeadL1(cmddata []byte) error {
	contentBufSize := len(cmddata)
	// 消息结构错误
	if this.MessageBinaryHeadL1.LowerSize() > contentBufSize {
		log.Error("[readBinaryNoHeadL1] "+
			"[缓冲区溢出] 接收消息格式错误 MessageBinaryHeadL1[%+v] "+
			"ContentBufSize[%d]",
			this.MessageBinaryHeadL1, contentBufSize)
		// 清空本消息信息
		this.Reset()
		return errors.New("消息头标注大小小于整体大小，消息体不完整")
	}
	// 检查 buffer
	if this.buffer == nil ||
		len(this.buffer) < int(this.MessageBinaryHeadL1.CmdLen) {
		tmpmsg := getMessageBinaryByProtoDataLength(
			int(this.MessageBinaryHeadL1.LowerSize()))
		if tmpmsg == nil {
			log.Error("[readBinaryNoHeadL1] "+
				"无法分配MsgBinary的内存！！！ Head[%+v]",
				this.MessageBinaryHeadL1)
			return nil
		}
		// 重新构建合理的 buffer
		this.buffer = tmpmsg.buffer
	}
	offset := 0
	offset += this.MessageBinaryHeadL1.WriteBinary(this.buffer[offset:])
	// 复制 MessageBinaryHeadL2+protodata 数据域
	copy(this.buffer[offset:this.MessageBinaryHeadL1.CmdLen],
		cmddata[:int(this.MessageBinaryHeadL1.CmdLen)-offset])
	// 将数据指针字段指向buffer数据域
	this.MessageBinaryBody.ProtoData =
		this.buffer[MSG_HEADSIZE:this.MessageBinaryHeadL1.CmdLen]

	return nil
}

func (this *MessageBinary) writeHeadBuffer() int {
	// 将结构数据填入 buffer
	offset := 0
	offset += this.MessageBinaryHeadL1.WriteBinary(this.buffer[offset:])
	return offset
}

// 将数据构造成为 Message 结构
// 在写入binary之前，必须经过 MakeMessage* 或 ReadBinary*
func (this *MessageBinary) WriteBinary() ([]byte, int) {
	// 如果缓冲区大小不合适，说明数据被篡改
	if this.buffer == nil ||
		len(this.buffer) < int(this.MessageBinaryHeadL1.CmdLen) {
		log.Error("[MakeMessageByBytes] "+
			"[WriteBinary] 错误的缓冲区大小，数据被篡改！ "+
			"BufferLen[%d] CmdLen[%d]",
			len(this.buffer), this.MessageBinaryHeadL1.CmdLen)
		return make([]byte, 1), 0
	}
	return this.buffer[:this.MessageBinaryHeadL1.CmdLen],
		int(this.MessageBinaryHeadL1.CmdLen)
}

// 通过二进制流创建 MessageBinary
func MakeMessageByBytes(cmdid uint16, protodata []byte) *MessageBinary {
	// 获取基础数据
	datalen := uint32(len(protodata))
	totalLength := uint32(MSG_HEADSIZE + datalen)
	// 判断数据合法性
	if totalLength >= MessageMaxSize {
		log.Error("[MakeMessageByBytes] "+
			"[缓冲区溢出] 发送消息数据过大 CmdID[%d] CmdLen[%d]",
			cmdid, totalLength)
		// 返回一个没有内容的消息
		msgbinary := getMessageBinaryByProtoDataLength(0)
		msgbinary.Reset()
		return msgbinary
	}
	// 从对象池获取消息对象
	msgbinary := getMessageBinaryByProtoDataLength(int(datalen))
	if msgbinary == nil {
		log.Error("[MakeMessageByBytes] "+
			"无法分配MsgBinary的内存！！！ CmdID[%d] Len[%d]",
			cmdid, totalLength)
		return nil
	}
	// 将 protodata 拷贝至 buffer 的数据域
	copy(msgbinary.buffer[MSG_HEADSIZE:totalLength], protodata)

	// 初始化消息信息

	// MessageBinaryBody
	// 消息数据字段指针指向 buffer 数据域
	msgbinary.MessageBinaryBody.ProtoData =
		msgbinary.buffer[MSG_HEADSIZE:totalLength]

	// MessageBinaryHeadL1
	msgbinary.MessageBinaryHeadL1.CmdLen = totalLength
	msgbinary.MessageBinaryHeadL1.CmdID = cmdid

	// 将结构数据填入 buffer
	msgbinary.writeHeadBuffer()

	return msgbinary
}

// 通过结构体创建 MessageBinary
func MakeMessageByObj(v MsgStruct) *MessageBinary {
	// 通过结构对象构造 json binary
	cmdid := v.GetMsgId()
	// 获取基础数据
	datalen := v.GetSize()
	totalLength := uint32(MSG_HEADSIZE + datalen)
	// 判断数据合法性
	if totalLength >= MessageMaxSize {
		log.Error("[MakeMessageByBytes] "+
			"[缓冲区溢出] 发送消息数据过大 MsgID[%d] CmdLen[%d] MaxSize[%d]",
			cmdid, totalLength, MessageMaxSize)
		// 返回一个没有内容的消息
		msgbinary := getMessageBinaryByProtoDataLength(0)
		msgbinary.Reset()
		return msgbinary
	}
	// 从对象池获取消息对象
	msgbinary := getMessageBinaryByProtoDataLength(int(datalen))
	if msgbinary == nil {
		log.Error("[MakeMessageByObj] "+
			"无法分配MsgBinary的内存！！！ CmdLen[%d] DataLen[%d]",
			totalLength, datalen)
		return nil
	}

	// MessageBinaryBody
	// 将 protodata 拷贝至 buffer 的数据域
	v.WriteBinary(msgbinary.buffer[MSG_HEADSIZE:totalLength])

	// 初始化消息信息

	// 消息数据字段指针指向 buffer 数据域
	msgbinary.MessageBinaryBody.ProtoData =
		msgbinary.buffer[MSG_HEADSIZE:totalLength]

	// 初始化消息信息
	// MessageBinaryHeadL1
	msgbinary.MessageBinaryHeadL1.CmdLen = totalLength
	msgbinary.MessageBinaryHeadL1.CmdID = cmdid

	// 将结构数据填入 buffer
	msgbinary.writeHeadBuffer()

	return msgbinary
}
