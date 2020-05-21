// Package msg micserver 中消息传递的基本单位
package msg

import (
	"encoding/hex"
	"errors"

	"github.com/liasece/micserver/log"
	msgbase "github.com/liasece/micserver/msg/base"
)

const (
	// MessageMaxSize 允许的最大消息大小（包括头部大小）
	// 普通代码中间应留出至少 32 字节提供给消息头部使用
	MessageMaxSize = 8 * 1024 * 1024
)

// MessageBinary micserver 中消息传递的基本单位
type MessageBinary struct {
	msgbase.MessageBase

	totalLength int
	msgID       uint16
	protoLength int

	ProtoData []byte
	buffer    []byte
}

// Free 将消息对象释放到对象池中
func (mb *MessageBinary) Free() {
	// 重置本消息各个属性
	mb.Reset()
	// 根据缓冲区容量归类
	size := len(mb.buffer)
	err := pools.Put(mb, size)
	if err != nil {
		log.Error("[MessageBinary.Free] pools.Put Err[%s]",
			err.Error())
	}
}

// Reset 重置 Message 数据
func (mb *MessageBinary) Reset() {
	mb.MessageBase.Reset()

	mb.totalLength = 0
	mb.msgID = 0
	mb.ProtoData = nil
	// 为了减轻GC压力，不应重置buffer字段
}

// GetTotalLength 获取消息包二进制总长度
func (mb *MessageBinary) GetTotalLength() int {
	return mb.totalLength
}

// SetTotalLength 设置消息包二进制总长度
func (mb *MessageBinary) SetTotalLength(v int) {
	mb.totalLength = v
}

// GetProtoLength 获取消息包数据段总长度
func (mb *MessageBinary) GetProtoLength() int {
	return mb.protoLength
}

// SetProtoLength 设置消息包数据段总长度
func (mb *MessageBinary) SetProtoLength(v int) {
	mb.protoLength = v
}

// GetMsgID 获取消息ID
func (mb *MessageBinary) GetMsgID() uint16 {
	return mb.msgID
}

// SetMsgID 设置消息ID
func (mb *MessageBinary) SetMsgID(v uint16) {
	mb.msgID = v
}

// 从二进制流中读取数据
// soffset 从本地偏移复制
// data 数据源
// doffset 数据源偏移
// Read bytenum 要读取的字节数
func (mb *MessageBinary) Read(soffset int,
	data []byte, doffset int, bytenum int) error {
	dataLen := len(data)
	// 消息结构错误
	if doffset+bytenum > dataLen {
		log.Error("[Read] 缓冲区溢出 DOffset[%d] ByteNum[%d] DataLen[%d]",
			doffset, bytenum, dataLen)
		return errors.New("源数据越界")
	}
	// 检查 buffer
	if mb.buffer == nil ||
		len(mb.buffer) < soffset+bytenum {
		tmpmsg := GetMessageBinary(soffset + bytenum)
		if tmpmsg == nil {
			log.Error("[Read] 无法分配MsgBinary的内存！！！ Size[%d]",
				soffset+bytenum)
			return errors.New("分配内存失败")
		}
		// 重新构建合理的 buffer
		mb.buffer = tmpmsg.buffer
	}
	// 复制 MessageBinaryHeadL2+protodata 数据域
	copy(mb.buffer[soffset:soffset+bytenum], data[doffset:doffset+bytenum])
	// 将数据指针字段指向buffer数据域
	return nil
}

// GetBuffer 获取消息的二进制流数据
func (mb *MessageBinary) GetBuffer() []byte {
	return mb.buffer
}

// SetProtoDataBound 设置消息内容数据域
func (mb *MessageBinary) SetProtoDataBound(offset int, bytenum int) error {
	if len(mb.buffer) < offset+bytenum || offset+bytenum < 0 {
		log.Error("[MessageBinary.SetProtoDataBound] 缓冲区越界 Len[%d] Need[%d]",
			len(mb.buffer), offset+bytenum)
		return errors.New("缓冲区越界")
	}
	mb.ProtoData = mb.buffer[offset : offset+bytenum]
	mb.SetProtoLength(bytenum)
	return nil
}

// String 获取消息的所有二进制内容的16进制字符串
func (mb *MessageBinary) String() string {
	if mb.buffer == nil {
		return ""
	}
	return hex.EncodeToString(mb.buffer)
}
