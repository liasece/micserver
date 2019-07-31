package msg

import (
	"encoding/binary"
	"encoding/hex"
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

const MSG_HEAD_SIZE = 16

func init() {
	pools = util.NewFlexiblePool(sizeControl, NewMsgBinaryBySize)
}

func NewMsgBinaryBySize(size int) interface{} {
	msg := new(MessageBinary)
	msg.buffers = make([]byte, size)
	return msg
}

type MessageBinary struct {
	// LV1 头
	CmdLen  uint32          // 4  消息长度
	CmdMask TEncryptionType // 1  消息是否加密
	CmdZip  byte            // 1  是否压缩

	// LV2 头
	CmdID     uint16 // 2
	TimeStamp uint32 // 4
	DataLen   uint32 // 4

	TmpData       interface{}       // 用于优化的临时数据指针请注意使用！
	TmpData1      interface{}       // 用于优化的临时数据指针请注意使用！
	OnSendDone    func(interface{}) // 用于优化的临时数据指针请注意使用！
	OnSendDoneArg interface{}       // 用于优化的临时数据指针请注意使用！
	ProtoData     []byte
	buffers       []byte
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
	totalSize := protoDataSize + MSG_HEAD_SIZE // 加上协议头长度
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

// 重置 Message 数据
func (this *MessageBinary) Reset() {
	this.CmdLen = 0
	this.CmdMask = 0
	this.CmdZip = 0
	this.CmdID = 0
	this.TimeStamp = 0
	this.DataLen = 0
	this.TmpData = nil
	this.TmpData1 = nil
	this.OnSendDone = nil
	this.OnSendDoneArg = nil
	this.ProtoData = nil
	// 为了减轻GC压力，不应重置buffers字段
}

// 从二进制流中读取 Message 结构，带消息头
func (this *MessageBinary) ReadBinary(cmddata []byte) error {
	// 获取基础数据
	maxlen := uint16(len(cmddata))
	// 过小的长度
	if maxlen < 6 {
		log.Error("[MakeMessageByBytes] "+
			"[ReadBinary] 错误的二进制数据,过小的[]byte DataLen[%d] Data[%s]",
			maxlen, hex.EncodeToString(cmddata))
		return errors.New("消息头接收不完整")
	}
	// 用于消息读取完毕之后的校验
	checklen := binary.BigEndian.Uint32(cmddata[0:4])
	this.CmdMask = TEncryptionType(cmddata[4])
	this.CmdZip = cmddata[3]
	// 读取消息
	err := this.ReadBinaryNoHead(cmddata[6:])
	if err != nil {
		log.Error("[MakeMessageByBytes] "+
			"[ReadBinary] ReadBinaryNoHead错误 Err[%s]",
			err.Error())
		return nil
	}
	// 长度检查
	if this.CmdLen != checklen {
		log.Error("[MakeMessageByBytes] "+
			"[ReadBinary] 错误的头部大小[%d] [%d]",
			checklen, this.CmdLen)
		return errors.New("消息头标注大小与实际大小不匹配")
	}
	return nil
}

// 从二进制流中读取 Message 结构，无消息头
func (this *MessageBinary) ReadBinaryNoHead(cmddata []byte) error {
	tmask := this.CmdMask
	tzip := this.CmdZip
	// 重置先前的消息
	this.Reset()
	maxlen := uint32(len(cmddata))
	if maxlen < 10 {
		log.Error("[ReadBinaryNoHead] "+
			"[ReadBinary] 错误的二进制数据,过小的[]byte NoHeadDataLen[%d]",
			maxlen)
		return errors.New("消息头接收不完整")
	}
	// 读消息名
	this.CmdID = binary.BigEndian.Uint16(cmddata[0:2])
	// 消息构造时间戳
	this.TimeStamp = binary.BigEndian.Uint32(cmddata[2:6])
	// 消息数据长度
	this.DataLen = binary.BigEndian.Uint32(cmddata[6:10])
	// 消息结构错误
	if this.DataLen+10 > maxlen {
		log.Error("[ReadBinaryNoHead] "+
			"[缓冲区溢出] 接收消息格式错误 CmdID[%d] CmdLen[%d] "+
			"DataLen[%d] RecvLen[%d]",
			this.CmdID, this.CmdLen, this.DataLen, maxlen)
		// 清空本消息信息
		this.Reset()
		return errors.New("消息头标注大小小于整体大小，消息体不完整")
	}
	// 总协议体长度
	this.CmdLen = MSG_HEAD_SIZE + this.DataLen
	// 检查 buffer
	if this.buffers == nil || len(this.buffers) < int(this.CmdLen) {
		tmpmsg := getMessageBinaryByProtoDataLength(int(this.DataLen))
		if tmpmsg == nil {
			log.Error("[ReadBinaryNoHead] "+
				"无法分配MsgBinary的内存！！！ Len[%d]", this.DataLen)
			return nil
		}
		// 重新构建合理的 buffer
		this.buffers = tmpmsg.buffers
	}
	// 复制 buffer 数据域
	copy(this.buffers[MSG_HEAD_SIZE:this.CmdLen], cmddata[10:10+this.DataLen])
	// 将数据指针字段指向buffer数据域
	this.ProtoData = this.buffers[MSG_HEAD_SIZE:this.CmdLen]

	this.CmdMask = tmask
	this.CmdZip = tzip
	// 将结构数据填入 buffer
	this.MakeMessageHead()

	// 解密消息
	this.Decrypt()
	return nil
}

func (this *MessageBinary) MakeMessageHead() {
	// 将结构数据填入 buffer
	binary.BigEndian.PutUint32(this.buffers[0:], this.CmdLen) // 4
	this.buffers[4] = byte(this.CmdMask)
	this.buffers[5] = this.CmdZip
	binary.BigEndian.PutUint16(this.buffers[6:], this.CmdID)     // 2
	binary.BigEndian.PutUint32(this.buffers[8:], this.TimeStamp) // 4
	binary.BigEndian.PutUint32(this.buffers[12:], this.DataLen)  // 4
}

// 将数据构造成为 Message 结构
// 在写入binary之前，必须经过 MakeMessage* 或 ReadBinary*
func (this *MessageBinary) WriteBinary() ([]byte, int) {
	// 如果缓冲区大小不合适，说明数据被篡改
	if this.buffers == nil ||
		len(this.buffers) < int(MSG_HEAD_SIZE+this.DataLen) {
		log.Error("[MakeMessageByBytes] "+
			"[WriteBinary] 错误的缓冲区大小，数据被篡改！ "+
			"BufferLen[%d] CmdLen[%d]",
			len(this.buffers), int(MSG_HEAD_SIZE+this.DataLen))
		return make([]byte, 1), 0
	}
	return this.buffers[:this.CmdLen], int(this.CmdLen)
}

func (this *MessageBinary) Encryption(t TEncryptionType) error {
	if this.CmdMask != 0 || t == 0x00 {
		return fmt.Errorf("[加密] 无效 %d -> %d ", this.CmdMask, t)
	}

	if t == EncryptionTypeXORSimple {
		// 异或加密

		// 第3个字节 加密标志字节
		this.CmdMask = t
		this.buffers[4] = byte(this.CmdMask)

		// 计算加密组长度
		modn := byte(this.TimeStamp&0x0FF)%5 + 5
		// 异或值
		xor := byte(this.DataLen & 0x0FF)
		for i, b := range this.ProtoData {
			// 异或
			this.ProtoData[i] = b ^ xor
			// 加上加密组长度
			this.ProtoData[i] += 10 - (byte(i&0x0FF) % modn)
		}
		return nil
	}
	return fmt.Errorf("[加密] 未知的消息加密类型 %d ", t)
}

func (this *MessageBinary) Decrypt() error {
	if this.CmdMask == 0 {
		return fmt.Errorf("[解密] 无效 %d", this.CmdMask)
	}
	if this.CmdMask == 0x01 {
		// 计算加密组长度
		modn := byte(this.TimeStamp&0x0FF)%5 + 5
		// 异或值
		xor := byte(this.DataLen & 0x0FF)
		for i, b := range this.ProtoData {
			// 先减去加密组长度
			b -= 10 - (byte(i&0x0FF) % modn)
			// 异或
			this.ProtoData[i] = b ^ xor
		}
		return nil
	}
	return fmt.Errorf("[解密] 未知的消息加密类型 %d ", this.CmdMask)
}

// 通过二进制流创建 Message
func MakeMessageByBytes(cmdid uint16, protodata []byte) *MessageBinary {
	// 获取基础数据
	datalen := uint32(len(protodata))
	totalLength := uint32(MSG_HEAD_SIZE + datalen)
	// 判断数据合法性
	if totalLength >= 64*1024 {
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

	totallen := uint32(totalLength)
	// 将 protodata 拷贝至 buffer 的数据域
	copy(msgbinary.buffers[MSG_HEAD_SIZE:totallen], protodata)
	// 消息数据字段指针指向 buffer 数据域
	msgbinary.ProtoData = msgbinary.buffers[MSG_HEAD_SIZE:totallen]

	// 初始化消息信息
	msgbinary.DataLen = uint32(datalen)
	msgbinary.CmdLen = totallen
	msgbinary.TimeStamp = uint32(time.Now().Unix())
	msgbinary.CmdID = cmdid

	// 将结构数据填入 buffer
	msgbinary.MakeMessageHead()

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
	totalLength := uint32(MSG_HEAD_SIZE + datalen)
	// 判断数据合法性
	if totalLength >= 64*1024 {
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

	totallen := uint32(totalLength)
	// 将 protodata 拷贝至 buffer 的数据域
	v.WriteBinary(msgbinary.buffers[MSG_HEAD_SIZE:totallen])
	// 消息数据字段指针指向 buffer 数据域
	msgbinary.ProtoData = msgbinary.buffers[MSG_HEAD_SIZE:totallen]

	// 初始化消息信息
	msgbinary.DataLen = uint32(datalen)
	msgbinary.CmdLen = totallen
	msgbinary.TimeStamp = uint32(time.Now().Unix())
	msgbinary.CmdID = cmdid

	// 将结构数据填入 buffer
	msgbinary.MakeMessageHead()
	return msgbinary
}

type MessageBinaryReader struct {
	inMsg     bool
	msglength int

	mask TEncryptionType
	zip  byte

	netbuffer *util.IOBuffer
}

func NewMessageBinaryReader(netbuffer *util.IOBuffer) *MessageBinaryReader {
	return &MessageBinaryReader{false, 0, 0x00, 0x00, netbuffer}
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
			headbuf, err := this.netbuffer.Read(0, 6)
			if err != nil {
				return err
			}
			// 消息总长度
			ulength := binary.BigEndian.Uint32(headbuf[0:4])
			if ulength < 6 {
				// 一个消息至少有头部6字节大小，如果小于6，说明数据已经不正确
				return fmt.Errorf("CmdLen too small. CmdLen:%d headdata:%#v",
					ulength, headbuf[0:6])
			}
			// 把头的6个字节去掉
			this.msglength = int(ulength) - 6
			if this.msglength < 0 {
				// 说明数据已经不正确
				return fmt.Errorf("msglength too small. "+
					"msglength:%d headdata:%#v",
					ulength, headbuf[0:6])
			}
			this.mask = TEncryptionType(headbuf[4])
			this.zip = headbuf[5]
			// 检查超长消息
			// 不会出现超长消息，最大长度即包头长度字段所能表示的最大长度

			// 进入消息处理逻辑
			this.inMsg = true
		}

		// 读消息体
		if this.inMsg && this.netbuffer.Len() >= this.msglength {
			// 取出消息体（无4个字节的头）
			cmdbuff, err := this.netbuffer.Read(0, this.msglength)
			if err != nil {
				return err
			}
			dataLength := len(cmdbuff) - (MSG_HEAD_SIZE - 6)
			// 获取合适大小的消息体
			msgbinary := getMessageBinaryByProtoDataLength(dataLength)
			if msgbinary != nil {
				msgbinary.CmdMask = this.mask
				msgbinary.CmdZip = this.zip
				// 解析消息（无6个字节的头）
				err := msgbinary.ReadBinaryNoHead(cmdbuff)
				if err != nil {
					log.Error("[MessageBinaryReader.RangeMsgBinary] "+
						"解析消息错误 Err[%s] RecvLen[%d] NoHeadLen[%d] "+
						"DataLen[%d]",
						err.Error(), len(cmdbuff), this.msglength, dataLength)
					return err
				} else {
					// 调用回调函数处理消息
					callback(msgbinary)
				}
			} else {
				log.Error("[MessageBinaryReader.RangeMsgBinary] "+
					"无法分配MsgBinary的内存！！！ RecvLen[%d] NoHeadLen[%d] "+
					"DataLen[%d]",
					len(cmdbuff), this.msglength, dataLength)
			}
			// 退出消息处理状态
			this.inMsg = false
		} else {
			break
		}
	}
	return nil
}
