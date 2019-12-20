package msg

import (
	"encoding/hex"
	"errors"

	"github.com/liasece/micserver/log"
	msgbase "github.com/liasece/micserver/msg/base"
)

const (
	// 允许的最大消息大小（包括头部大小）
	// 普通代码中间应留出至少 32 字节提供给消息头部使用
	MessageMaxSize = 8 * 1024 * 1024
)

type MessageBinary struct {
	msgbase.MessageBase

	totalLength int
	msgID       uint16
	protoLength int

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

	this.totalLength = 0
	this.msgID = 0
	this.ProtoData = nil
	// 为了减轻GC压力，不应重置buffer字段
}

func (this *MessageBinary) GetTotalLength() int {
	return this.totalLength
}

func (this *MessageBinary) SetTotalLength(v int) {
	this.totalLength = v
}

func (this *MessageBinary) GetProtoLength() int {
	return this.protoLength
}

func (this *MessageBinary) SetProtoLength(v int) {
	this.protoLength = v
}

func (this *MessageBinary) GetMsgID() uint16 {
	return this.msgID
}

func (this *MessageBinary) SetMsgID(v uint16) {
	this.msgID = v
}

// 从二进制流中读取数据
// soffset 从本地偏移复制
// data 数据源
// doffset 数据源偏移
// bytenum 要读取的字节数
func (this *MessageBinary) Read(soffset int,
	data []byte, doffset int, bytenum int) error {
	dataLen := len(data)
	// 消息结构错误
	if doffset+bytenum > dataLen {
		log.Error("[Read] 缓冲区溢出 DOffset[%d] ByteNum[%d] DataLen[%d]",
			doffset, bytenum, dataLen)
		return errors.New("源数据越界")
	}
	// 检查 buffer
	if this.buffer == nil ||
		len(this.buffer) < soffset+bytenum {
		tmpmsg := GetMessageBinary(soffset + bytenum)
		if tmpmsg == nil {
			log.Error("[Read] 无法分配MsgBinary的内存！！！ Size[%d]",
				soffset+bytenum)
			return errors.New("分配内存失败")
		}
		// 重新构建合理的 buffer
		this.buffer = tmpmsg.buffer
	}
	// 复制 MessageBinaryHeadL2+protodata 数据域
	copy(this.buffer[soffset:soffset+bytenum], data[doffset:doffset+bytenum])
	// 将数据指针字段指向buffer数据域
	return nil
}

// 获取消息的二进制流数据
func (this *MessageBinary) GetBuffer() []byte {
	return this.buffer
}

// 设置消息内容数据域
func (this *MessageBinary) SetProtoDataBound(offset int, bytenum int) error {
	if len(this.buffer) < offset+bytenum || offset+bytenum < 0 {
		log.Error("[MessageBinary.SetProtoDataBound] 缓冲区越界 Len[%d] Need[%d]",
			len(this.buffer), offset+bytenum)
		return errors.New("缓冲区越界")
	}
	this.ProtoData = this.buffer[offset : offset+bytenum]
	this.SetProtoLength(bytenum)
	return nil
}

// 获取消息的所有二进制内容的16进制字符串
func (this *MessageBinary) String() string {
	if this.buffer == nil {
		return ""
	}
	return hex.EncodeToString(this.buffer)
}
