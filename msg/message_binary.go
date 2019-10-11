package msg

import (
	"errors"
	"github.com/liasece/micserver/log"
	msgbase "github.com/liasece/micserver/msg/base"
)

const (
	// 允许的最大消息大小（包括头部大小）
	// 普通代码中间应留出至少 32 字节提供给消息头部使用
	MessageMaxSize = 8 * 1024 * 1024
)

var (
	defaultHead1 MessageBinaryHeadL1
)

type MessageBinary struct {
	msgbase.MessageBase

	MessageBinaryHeadL1

	ProtoData []byte
	buffer    []byte
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
	this.MessageBase.Reset()

	this.MessageBinaryHeadL1 = defaultHead1
	this.ProtoData = nil
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
	this.ProtoData =
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
		log.Error("[MessageBinary.WriteBinary] 错误的缓冲区大小，数据被篡改！ "+
			"BufferLen[%d] CmdLen[%d]",
			len(this.buffer), this.MessageBinaryHeadL1.CmdLen)
		return make([]byte, 0), 0
	}
	return this.buffer[:this.MessageBinaryHeadL1.CmdLen],
		int(this.MessageBinaryHeadL1.CmdLen)
}
