package msg

import (
	"errors"
	"fmt"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/util"
	"time"
)

type TEncryptionType byte

const (
	EncryptionTypeNone      TEncryptionType = 0x00
	EncryptionTypeXORSimple TEncryptionType = 0x01
)

// 粒度控制 单位：字节
var sizeControl []int = []int{32, 64, 128, 256, 512, 1024, 2 * 1024,
	4 * 1024, 6 * 1024, 8 * 1024, 10 * 1024, 15 * 1024, 20 * 1024,
	25 * 1024, 30 * 1024, 35 * 1024, 40 * 1024, 45 * 1024, 50 * 1024,
	55 * 1024, 60 * 1024, 64 * 1024, 128 * 1024, 256 * 1024, 512 * 1024,
	1024 * 1024, 2 * 1024 * 1024, 4 * 1024 * 1024, 8 * 1024 * 1024,
	16 * 1024 * 1024, 32 * 1024 * 1024, 64 * 1024 * 1024, 128 * 1024 * 1024,
	256 * 1024 * 1024, 512 * 1024 * 1024}
var pools *util.FlexiblePool

func init() {
	pools = util.NewFlexiblePool(sizeControl, NewMsgBinaryBySize)
}

func NewMsgBinaryBySize(size int) interface{} {
	msg := new(MessageBinary)
	msg.buffers = make([]byte, size)
	return msg
}

type MessageBinary struct {
	MessageBinaryHeadL1
	MessageBinaryHeadL2
	MessageBinaryBody
	buffers []byte

	TmpData       interface{}       // 用于优化的临时数据指针请注意使用！
	TmpData1      interface{}       // 用于优化的临时数据指针请注意使用！
	OnSendDone    func(interface{}) // 用于优化的临时数据指针请注意使用！
	OnSendDoneArg interface{}       // 用于优化的临时数据指针请注意使用！
}

func (this *MessageBinary) Free() {
	// 重置本消息各个属性
	this.Reset()
	// 根据缓冲区容量归类
	size := len(this.buffers)
	err := pools.Put(this, size)
	if err != nil {
		log.Error("[MessageBinary.Free] 放入对象池错误 Err[%s]",
			err.Error())
	}
}

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

var defaultHead1 MessageBinaryHeadL1
var defaultHead2 MessageBinaryHeadL2
var defaultBody MessageBinaryBody

// 重置 Message 数据
func (this *MessageBinary) Reset() {
	this.MessageBinaryHeadL1 = defaultHead1
	this.MessageBinaryHeadL2 = defaultHead2
	this.TmpData = nil
	this.TmpData1 = nil
	this.OnSendDone = nil
	this.OnSendDoneArg = nil
	this.MessageBinaryBody = defaultBody
	// 为了减轻GC压力，不应重置buffers字段
}

// 从二进制流中读取 Message 结构，带消息头
func (this *MessageBinary) ReadBinary(cmddata []byte) error {
	// 获取基础数据
	offset, err := this.MessageBinaryHeadL1.ReadFromBuffer(cmddata)
	// 过小的长度
	if err != nil {
		log.Error("[MakeMessageByBytes] "+
			"[ReadBinary] 错误的消息头L1 Err[%s]", err.Error())
		return errors.New("错误的消息头L1")
	}
	// 用于消息读取完毕之后的校验
	// 读取消息
	err = this.ReadBinaryNoMessageBinaryHeadL1(cmddata[offset:])
	if err != nil {
		log.Error("[MakeMessageByBytes] "+
			"[ReadBinary] ReadBinaryNoMessageBinaryHeadL1错误 Err[%s]",
			err.Error())
		return nil
	}
	return nil
}

// 从二进制流中读取 Message 结构，无消息头
func (this *MessageBinary) ReadBinaryNoMessageBinaryHeadL1(cmddata []byte) error {
	offset, err := this.MessageBinaryHeadL2.ReadFromBuffer(cmddata)
	contentBufSize := len(cmddata) - offset
	if err != nil {
		log.Error("[ReadBinaryNoMessageBinaryHeadL1] "+
			"[ReadBinary] 错误的消息头L2 Err[%s]", err.Error())
		return errors.New("错误的消息头L2")
	}
	// 消息结构错误
	if this.MessageBinaryHeadL2.LowerSize() > contentBufSize {
		log.Error("[ReadBinaryNoMessageBinaryHeadL1] "+
			"[缓冲区溢出] 接收消息格式错误 MessageBinaryHeadL1[%+v] MessageBinaryHeadL2[%+v] "+
			"ContentBufSize[%d]",
			this.MessageBinaryHeadL1, this.MessageBinaryHeadL2, contentBufSize)
		// 清空本消息信息
		this.Reset()
		return errors.New("消息头标注大小小于整体大小，消息体不完整")
	}
	// 检查 buffer
	if this.buffers == nil || len(this.buffers) < int(this.MessageBinaryHeadL1.CmdLen) {
		tmpmsg := getMessageBinaryByProtoDataLength(int(this.MessageBinaryHeadL2.DataLen))
		if tmpmsg == nil {
			log.Error("[ReadBinaryNoMessageBinaryHeadL1] "+
				"无法分配MsgBinary的内存！！！ Len[%d]", this.MessageBinaryHeadL2.DataLen)
			return nil
		}
		// 重新构建合理的 buffer
		this.buffers = tmpmsg.buffers
	}
	offset = 0
	offset += this.MessageBinaryHeadL1.WriteToBuffer(this.buffers[offset:])
	// 复制 MessageBinaryHeadL2+protodata 数据域
	copy(this.buffers[offset:this.MessageBinaryHeadL1.CmdLen],
		cmddata[:int(this.MessageBinaryHeadL1.CmdLen)-offset])
	// 将数据指针字段指向buffer数据域
	this.MessageBinaryBody.ProtoData = this.buffers[MSG_HEADSIZE:this.MessageBinaryHeadL1.CmdLen]

	// 解密消息
	this.Decrypt()
	return nil
}

func (this *MessageBinary) writeHeadBuffer() int {
	// 将结构数据填入 buffer
	offset := 0
	offset += this.MessageBinaryHeadL1.WriteToBuffer(this.buffers[offset:])
	offset += this.MessageBinaryHeadL2.WriteToBuffer(this.buffers[offset:])
	return offset
}

// 将数据构造成为 Message 结构
// 在写入binary之前，必须经过 MakeMessage* 或 ReadBinary*
func (this *MessageBinary) WriteBinary() ([]byte, int) {
	// 如果缓冲区大小不合适，说明数据被篡改
	if this.buffers == nil ||
		len(this.buffers) < int(this.MessageBinaryHeadL1.CmdLen) {
		log.Error("[MakeMessageByBytes] "+
			"[WriteBinary] 错误的缓冲区大小，数据被篡改！ "+
			"BufferLen[%d] CmdLen[%d]",
			len(this.buffers), this.MessageBinaryHeadL1.CmdLen)
		return make([]byte, 1), 0
	}
	return this.buffers[:this.MessageBinaryHeadL1.CmdLen], int(this.MessageBinaryHeadL1.CmdLen)
}

func (this *MessageBinary) Encryption(t TEncryptionType) error {
	if this.MessageBinaryHeadL1.CmdMask != 0 || t == 0x00 {
		return fmt.Errorf("[加密] 无效 %d -> %d ", this.MessageBinaryHeadL1.CmdMask, t)
	}

	if t == EncryptionTypeXORSimple {
		// 异或加密

		// 第3个字节 加密标志字节
		this.MessageBinaryHeadL1.CmdMask = t
		this.buffers[4] = byte(this.MessageBinaryHeadL1.CmdMask)

		// 计算加密组长度
		modn := byte(this.MessageBinaryHeadL2.TimeStamp&0x0FF)%5 + 5
		// 异或值
		xor := byte(this.MessageBinaryHeadL2.DataLen & 0x0FF)
		for i, b := range this.MessageBinaryBody.ProtoData {
			// 异或
			this.MessageBinaryBody.ProtoData[i] = b ^ xor
			// 加上加密组长度
			this.MessageBinaryBody.ProtoData[i] += 10 - (byte(i&0x0FF) % modn)
		}
		return nil
	}
	return fmt.Errorf("[加密] 未知的消息加密类型 %d ", t)
}

func (this *MessageBinary) Decrypt() error {
	if this.MessageBinaryHeadL1.CmdMask == 0 {
		return fmt.Errorf("[解密] 无效 %d", this.MessageBinaryHeadL1.CmdMask)
	}
	if this.MessageBinaryHeadL1.CmdMask == 0x01 {
		// 计算加密组长度
		modn := byte(this.MessageBinaryHeadL2.TimeStamp&0x0FF)%5 + 5
		// 异或值
		xor := byte(this.MessageBinaryHeadL2.DataLen & 0x0FF)
		for i, b := range this.MessageBinaryBody.ProtoData {
			// 先减去加密组长度
			b -= 10 - (byte(i&0x0FF) % modn)
			// 异或
			this.MessageBinaryBody.ProtoData[i] = b ^ xor
		}
		return nil
	}
	return fmt.Errorf("[解密] 未知的消息加密类型 %d ", this.MessageBinaryHeadL1.CmdMask)
}

// 通过二进制流创建 Message
func MakeMessageByBytes(cmdid uint16, protodata []byte) *MessageBinary {
	// 获取基础数据
	datalen := uint32(len(protodata))
	totalLength := uint32(MSG_HEADSIZE + datalen)
	// 判断数据合法性
	if totalLength >= 512*1024*1024 {
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

	// 初始化消息信息

	// MessageBinaryBody
	// 将 protodata 拷贝至 buffer 的数据域
	copy(msgbinary.buffers[MSG_HEADSIZE:totalLength], protodata)
	// 消息数据字段指针指向 buffer 数据域
	msgbinary.MessageBinaryBody.ProtoData =
		msgbinary.buffers[MSG_HEADSIZE:msgbinary.MessageBinaryHeadL1.CmdLen]

	// MessageBinaryHeadL2
	msgbinary.MessageBinaryHeadL2.DataLen = datalen
	msgbinary.MessageBinaryHeadL2.TimeStamp = uint32(time.Now().Unix())
	msgbinary.MessageBinaryHeadL2.CmdID = cmdid

	// MessageBinaryHeadL1
	msgbinary.MessageBinaryHeadL1.CmdLen = totalLength

	// 将结构数据填入 buffer
	msgbinary.writeHeadBuffer()

	return msgbinary
}

type MsgStruct interface {
	WriteBinary(data []byte) int
	GetMsgId() uint16
	GetMsgName() string
	GetSize() int
	GetJson() string
}

// 通过结构体创建 Json Message
func MakeMessageByJson(v MsgStruct) *MessageBinary {
	// 通过结构对象构造 json binary
	cmdid := v.GetMsgId()
	// 获取基础数据
	datalen := v.GetSize()
	totalLength := uint32(MSG_HEADSIZE + datalen)
	// 判断数据合法性
	if totalLength >= 512*1024*1024 {
		log.Error("[MakeMessageByBytes] "+
			"[缓冲区溢出] 发送消息数据过大 MsgID[%d] CmdLen[%d]",
			cmdid, totalLength)
		// 返回一个没有内容的消息
		msgbinary := getMessageBinaryByProtoDataLength(0)
		msgbinary.Reset()
		return msgbinary
	}
	// 从对象池获取消息对象
	msgbinary := getMessageBinaryByProtoDataLength(int(datalen))
	if msgbinary == nil {
		log.Error("[MakeMessageByJson] "+
			"无法分配MsgBinary的内存！！！ CmdLen[%d] DataLen[%d]",
			totalLength, datalen)
		return nil
	}

	// 初始化消息信息

	// MessageBinaryBody
	// 将 protodata 拷贝至 buffer 的数据域

	v.WriteBinary(msgbinary.buffers[MSG_HEADSIZE:totalLength])

	// 消息数据字段指针指向 buffer 数据域
	msgbinary.MessageBinaryBody.ProtoData = msgbinary.buffers[MSG_HEADSIZE:totalLength]

	// MessageBinaryHeadL2
	// 初始化消息信息
	msgbinary.MessageBinaryHeadL2.DataLen = uint32(datalen)
	msgbinary.MessageBinaryHeadL2.TimeStamp = uint32(time.Now().Unix())
	msgbinary.MessageBinaryHeadL2.CmdID = cmdid
	// MessageBinaryHeadL1
	msgbinary.MessageBinaryHeadL1.CmdLen = totalLength

	// 将结构数据填入 buffer
	msgbinary.writeHeadBuffer()

	return msgbinary
}

type MessageBinaryReader struct {
	inMsg  bool
	HeadL1 MessageBinaryHeadL1
	HeadL2 MessageBinaryHeadL2

	netbuffer *util.IOBuffer
}

func NewMessageBinaryReader(netbuffer *util.IOBuffer) *MessageBinaryReader {
	return &MessageBinaryReader{
		netbuffer: netbuffer,
	}
}

func (this *MessageBinaryReader) RangeMsgBinary(
	callback func(*MessageBinary)) (reerr error) {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := util.GetPanicInfo(recover()); err != nil {
			log.Error("[MessageBinaryReader.RangeMsgBinary] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
			reerr = err
		}
	}()

	// 遍历数据流中的消息体
	for {
		// 读消息头
		// 当前不在消息体中，且当前缓冲区长度已大于消息头长度
		if !this.inMsg && this.netbuffer.Len() >= 6 {

			// 读头部4个字节
			MessageBinaryHeadL1buf, err := this.netbuffer.Read(0, 6)
			if err != nil {
				return err
			}
			_, err = this.HeadL1.ReadFromBuffer(MessageBinaryHeadL1buf)
			if err != nil {
				return fmt.Errorf("Head layer 1 dec err:%s. headdata:%#v",
					err.Error(), MessageBinaryHeadL1buf)
			}

			// 进入消息处理逻辑
			this.inMsg = true

		}

		// 读消息体
		if this.inMsg && this.netbuffer.Len() >= this.HeadL1.LowerSize() {

			cmdbuff, err := this.netbuffer.Read(0, this.HeadL1.LowerSize())

			if err != nil {
				return err
			}
			_, err = this.HeadL2.ReadFromBuffer(cmdbuff)

			if err != nil {
				return fmt.Errorf("Head layer 2 dec err:%s.",
					err.Error())
			}

			// 解密解压
			// TODO

			// 获取合适大小的消息体
			msgbinary := getMessageBinaryByProtoDataLength(
				this.HeadL2.LowerSize())
			if msgbinary != nil {
				msgbinary.MessageBinaryHeadL1 = this.HeadL1
				msgbinary.MessageBinaryHeadL2 = this.HeadL2
				// 解析消息（无6个字节的头）
				err := msgbinary.ReadBinaryNoMessageBinaryHeadL1(cmdbuff)
				if err != nil {
					log.Error("[MessageBinaryReader.RangeMsgBinary] "+
						"解析消息错误 Err[%s] RecvLen[%d] HeadL1[%+v] "+
						"HeadL2[%+v]",
						err.Error(), len(cmdbuff), this.HeadL1, this.HeadL2)
					return err
				} else {
					// 调用回调函数处理消息
					callback(msgbinary)
				}
			} else {
				log.Error("[MessageBinaryReader.RangeMsgBinary] "+
					"无法分配MsgBinary的内存！！！ RecvLen[%d] HeadL1[%+v] "+
					"HeadL2[%+v]", len(cmdbuff), this.HeadL1, this.HeadL2)
			}
			// 退出消息处理状态
			this.inMsg = false
		} else {
			break
		}
	}
	return nil
}
