package msg

import (
	"encoding/binary"
	"fmt"

	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/util/buffer"
	"github.com/liasece/micserver/util/sysutil"
)

// 默认的消息头大小
const (
	DEFAULT_MSG_HEADSIZE = (4 + 2)
)

// 默认通过结构构造消息体
func DefaultEncodeObj(v MsgStruct) *MessageBinary {
	// 通过结构对象构造 json binary
	cmdid := v.GetMsgId()
	// 获取基础数据
	datalen := v.GetSize()
	totalLength := DEFAULT_MSG_HEADSIZE + datalen
	// 判断数据合法性
	if totalLength >= MessageMaxSize {
		log.Error("[EncodeObj] "+
			"[缓冲区溢出] 发送消息数据过大 MsgID[%d] TotalLen[%d] MaxSize[%d]",
			cmdid, totalLength, MessageMaxSize)
		// 返回一个没有内容的消息
		msgbinary := GetMessageBinary(0)
		msgbinary.Reset()
		return msgbinary
	}
	// 从对象池获取消息对象
	msgbinary := GetMessageBinary(totalLength)
	if msgbinary == nil {
		log.Error("[GetByObj] "+
			"无法分配MsgBinary的内存！！！ TotalLen[%d] DataLen[%d]",
			totalLength, datalen)
		return nil
	}

	// MessageBinaryBody
	// 将 protodata 拷贝至 buffer 的数据域
	v.WriteBinary(msgbinary.GetBuffer()[DEFAULT_MSG_HEADSIZE:totalLength])

	// 初始化消息信息

	// 消息数据字段指针指向 buffer 数据域
	msgbinary.SetProtoDataBound(DEFAULT_MSG_HEADSIZE, datalen)

	// 初始化消息信息
	// MessageBinaryHeadL1
	msgbinary.SetTotalLength(totalLength)
	msgbinary.SetMsgID(cmdid)
	DefaultWriteHead(msgbinary.GetBuffer(), totalLength, cmdid)

	return msgbinary
}

// 默认通过字节流构造消息体
func DefaultEncodeBytes(cmdid uint16, protodata []byte) *MessageBinary {
	// 获取基础数据
	datalen := len(protodata)
	totalLength := DEFAULT_MSG_HEADSIZE + datalen
	// 判断数据合法性
	if totalLength >= MessageMaxSize {
		log.Error("[EncodeObj] "+
			"[缓冲区溢出] 发送消息数据过大 MsgID[%d] TotalLen[%d] MaxSize[%d]",
			cmdid, totalLength, MessageMaxSize)
		// 返回一个没有内容的消息
		msgbinary := GetMessageBinary(0)
		msgbinary.Reset()
		return msgbinary
	}
	// 从对象池获取消息对象
	msgbinary := GetMessageBinary(totalLength)
	if msgbinary == nil {
		log.Error("[GetByObj] "+
			"无法分配MsgBinary的内存！！！ TotalLen[%d] DataLen[%d]",
			totalLength, datalen)
		return nil
	}

	// MessageBinaryBody
	// 将 protodata 拷贝至 buffer 的数据域
	msgbinary.Read(DEFAULT_MSG_HEADSIZE, protodata, 0, datalen)

	// 初始化消息信息
	// 消息数据字段指针指向 buffer 数据域
	msgbinary.SetProtoDataBound(DEFAULT_MSG_HEADSIZE, datalen)

	// 初始化消息信息
	msgbinary.SetTotalLength(totalLength)
	msgbinary.SetMsgID(cmdid)
	DefaultWriteHead(msgbinary.GetBuffer(), totalLength, cmdid)

	return msgbinary
}

// 默认写头
func DefaultWriteHead(data []byte, totalLen int,
	msgid uint16) (size int) {
	binary.LittleEndian.PutUint32(data[size:], uint32(totalLen)) // 4
	size += 4
	binary.LittleEndian.PutUint16(data[size:], msgid) // 2
	size += 2

	return
}

// micserver 默认的消息编解码器
type DefaultCodec struct {
	inMsg bool

	totalLen int
	protoLen int
	msgID    uint16
}

// 遍历目标缓冲区中的消息
func (this *DefaultCodec) RangeMsgBinary(
	buf *buffer.IOBuffer, cb func(*MessageBinary)) (reerr error) {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := sysutil.GetPanicInfo(recover()); err != nil {
			log.Error("[DefaultCodec.RangeMsgBinary] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
			reerr = err
		}
	}()

	// 遍历数据流中的消息体
	for {
		// 读消息头
		// 当前不在消息体中，且当前缓冲区长度已大于消息头长度
		if !this.inMsg && buf.Len() >= DEFAULT_MSG_HEADSIZE {

			// 读头部4个字节
			headData, err := buf.Read(0, DEFAULT_MSG_HEADSIZE)
			if err != nil {
				return err
			}
			_, err = this.readHead(headData)
			if err != nil {
				return fmt.Errorf("Head dec err:%s. headdata:%#v",
					err.Error(), headData)
			}

			// 进入消息处理逻辑
			this.inMsg = true
		}

		// 读消息体
		if this.inMsg && buf.Len() >= this.protoLen {
			cmdbuff, err := buf.Read(0, this.protoLen)
			if err != nil {
				return err
			}

			// 获取合适大小的消息体
			msgbinary := GetMessageBinary(this.totalLen)
			if msgbinary != nil {
				msgbinary.SetTotalLength(this.totalLen)
				msgbinary.SetMsgID(this.msgID)
				// 解析消息（无6个字节的头）
				err := msgbinary.Read(DEFAULT_MSG_HEADSIZE, cmdbuff, 0,
					this.protoLen)
				if err != nil {
					log.Error("[DefaultCodec.RangeMsgBinary] "+
						"解析消息错误 Err[%s]", err.Error())
					return err
				} else {
					// 设置内容边界
					msgbinary.SetProtoDataBound(DEFAULT_MSG_HEADSIZE,
						this.protoLen)
					// 调用回调函数处理消息
					cb(msgbinary)
				}
			} else {
				log.Error("[DefaultCodec.RangeMsgBinary] "+
					"无法分配MsgBinary的内存！！！ TotalLen[%d]", this.totalLen)
			}
			// 退出消息处理状态
			this.inMsg = false
		} else {
			break
		}
	}
	return nil
}

func (this *DefaultCodec) readHead(data []byte) (int, error) {
	if len(data) < DEFAULT_MSG_HEADSIZE {
		return 0, fmt.Errorf("data not enough")
	}
	size := 0
	this.totalLen = int(binary.LittleEndian.Uint32(data[size:]))
	this.protoLen = this.totalLen - DEFAULT_MSG_HEADSIZE
	size += 4
	this.msgID = binary.LittleEndian.Uint16(data[size:])
	size += 2
	// cmdlen must  include head layer 1 size
	if int(this.totalLen) < DEFAULT_MSG_HEADSIZE {
		return 0, fmt.Errorf("error cmdlen")
	}
	return DEFAULT_MSG_HEADSIZE, nil
}

// 编码一个消息对象
func (this *DefaultCodec) EncodeObj(v MsgStruct) *MessageBinary {
	return DefaultEncodeObj(v)
}

// 编码一个由消息号及二进制内容构成的消息
func (this *DefaultCodec) EncodeBytes(cmdid uint16,
	protodata []byte) *MessageBinary {
	return DefaultEncodeBytes(cmdid, protodata)
}
