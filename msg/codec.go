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
	DEFAULTMSGHEADSIZE = (4 + 2)
)

// DefaultEncodeObj 默认通过结构构造消息体
func DefaultEncodeObj(v IMsgStruct) *MessageBinary {
	// 通过结构对象构造 json binary
	cmdid := v.GetMsgId()
	// 获取基础数据
	datalen := v.GetSize()
	totalLength := DEFAULTMSGHEADSIZE + datalen
	// 判断数据合法性
	if totalLength >= MessageMaxSize {
		log.Error("[DefaultEncodeObj] Buffer overflow, sending message data too large", log.Uint16("MsgID", cmdid), log.Int("TotalLen", totalLength), log.Int32("MaxSize", MessageMaxSize))
		// 返回一个没有内容的消息
		msgbinary := GetMessageBinary(0)
		msgbinary.Reset()
		return msgbinary
	}
	// 从对象池获取消息对象
	msgbinary := GetMessageBinary(totalLength)
	if msgbinary == nil {
		log.Error("[DefaultEncodeObj] Unable to allocate MsgBinary's memory!!!!", log.Int("TotalLen", totalLength), log.Int("DataLen", datalen))
		return nil
	}

	// MessageBinaryBody
	// 将 protodata 拷贝至 buffer 的数据域
	v.WriteBinary(msgbinary.GetBuffer()[DEFAULTMSGHEADSIZE:totalLength])

	// 初始化消息信息

	// 消息数据字段指针指向 buffer 数据域
	msgbinary.SetProtoDataBound(DEFAULTMSGHEADSIZE, datalen)

	// 初始化消息信息
	// MessageBinaryHeadL1
	msgbinary.SetTotalLength(totalLength)
	msgbinary.SetMsgID(cmdid)
	DefaultWriteHead(msgbinary.GetBuffer(), totalLength, cmdid)

	return msgbinary
}

// DefaultEncodeBytes 默认通过字节流构造消息体
func DefaultEncodeBytes(cmdid uint16, protodata []byte) *MessageBinary {
	// 获取基础数据
	datalen := len(protodata)
	totalLength := DEFAULTMSGHEADSIZE + datalen
	// 判断数据合法性
	if totalLength >= MessageMaxSize {
		log.Error("[DefaultEncodeBytes] Buffer overflow, sending message data too large", log.Uint16("MsgID", cmdid), log.Int("TotalLen", totalLength), log.Int32("MaxSize", MessageMaxSize))
		// 返回一个没有内容的消息
		msgbinary := GetMessageBinary(0)
		msgbinary.Reset()
		return msgbinary
	}
	// 从对象池获取消息对象
	msgbinary := GetMessageBinary(totalLength)
	if msgbinary == nil {
		log.Error("[DefaultEncodeBytes] Unable to allocate MsgBinary's memory!!!!", log.Int("TotalLen", totalLength), log.Int("DataLen", datalen))
		return nil
	}

	// MessageBinaryBody
	// 将 protodata 拷贝至 buffer 的数据域
	msgbinary.Read(DEFAULTMSGHEADSIZE, protodata, 0, datalen)

	// 初始化消息信息
	// 消息数据字段指针指向 buffer 数据域
	msgbinary.SetProtoDataBound(DEFAULTMSGHEADSIZE, datalen)

	// 初始化消息信息
	msgbinary.SetTotalLength(totalLength)
	msgbinary.SetMsgID(cmdid)
	DefaultWriteHead(msgbinary.GetBuffer(), totalLength, cmdid)

	return msgbinary
}

// DefaultWriteHead 默认写头
func DefaultWriteHead(data []byte, totalLen int,
	msgid uint16) (size int) {
	binary.LittleEndian.PutUint32(data[size:], uint32(totalLen)) // 4
	size += 4
	binary.LittleEndian.PutUint16(data[size:], msgid) // 2
	size += 2

	return
}

// DefaultCodec micserver 默认的消息编解码器
type DefaultCodec struct {
	inMsg bool

	totalLen int
	protoLen int
	msgID    uint16
}

// RangeMsgBinary 遍历目标缓冲区中的消息
func (defaultCodec *DefaultCodec) RangeMsgBinary(
	buf *buffer.IOBuffer, cb func(*MessageBinary)) (reerr error) {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if stackInfo, err := sysutil.GetPanicInfo(recover()); err != nil {
			log.Error("[DefaultCodec.RangeMsgBinary] Panic", log.ErrorField(err), log.String("Stack", stackInfo))
			reerr = err
		}
	}()

	// 遍历数据流中的消息体
	for {
		// 读消息头
		// 当前不在消息体中，且当前缓冲区长度已大于消息头长度
		if !defaultCodec.inMsg && buf.Len() >= DEFAULTMSGHEADSIZE {

			// 读头部4个字节
			headData, err := buf.Read(0, DEFAULTMSGHEADSIZE)
			if err != nil {
				return err
			}
			_, err = defaultCodec.readHead(headData)
			if err != nil {
				return fmt.Errorf("Head dec err:%s. headdata:%#v", err.Error(), headData)
			}

			// 进入消息处理逻辑
			defaultCodec.inMsg = true
		}

		// 读消息体
		if defaultCodec.inMsg && buf.Len() >= defaultCodec.protoLen {
			cmdbuff, err := buf.Read(0, defaultCodec.protoLen)
			if err != nil {
				return err
			}

			// 获取合适大小的消息体
			msgbinary := GetMessageBinary(defaultCodec.totalLen)
			if msgbinary != nil {
				msgbinary.SetTotalLength(defaultCodec.totalLen)
				msgbinary.SetMsgID(defaultCodec.msgID)
				// 解析消息（无6个字节的头）
				err := msgbinary.Read(DEFAULTMSGHEADSIZE, cmdbuff, 0,
					defaultCodec.protoLen)
				if err != nil {
					log.Error("[DefaultCodec.RangeMsgBinary] Parse message error", log.ErrorField(err))
					return err
				}
				// 设置内容边界
				msgbinary.SetProtoDataBound(DEFAULTMSGHEADSIZE, defaultCodec.protoLen)
				// 调用回调函数处理消息
				cb(msgbinary)
			} else {
				log.Error("[DefaultCodec.RangeMsgBinary] Unable to allocate MsgBinary's memory!!!!", log.Int("TotalLen", defaultCodec.totalLen))
			}
			// 退出消息处理状态
			defaultCodec.inMsg = false
		} else {
			break
		}
	}
	return nil
}

func (defaultCodec *DefaultCodec) readHead(data []byte) (int, error) {
	if len(data) < DEFAULTMSGHEADSIZE {
		return 0, fmt.Errorf("data not enough")
	}
	size := 0
	defaultCodec.totalLen = int(binary.LittleEndian.Uint32(data[size:]))
	defaultCodec.protoLen = defaultCodec.totalLen - DEFAULTMSGHEADSIZE
	size += 4
	defaultCodec.msgID = binary.LittleEndian.Uint16(data[size:])
	size += 2
	// cmdlen must  include head layer 1 size
	if int(defaultCodec.totalLen) < DEFAULTMSGHEADSIZE {
		return 0, fmt.Errorf("error cmdlen")
	}
	return DEFAULTMSGHEADSIZE, nil
}

// EncodeObj 编码一个消息对象
func (defaultCodec *DefaultCodec) EncodeObj(v IMsgStruct) *MessageBinary {
	return DefaultEncodeObj(v)
}

// EncodeBytes 编码一个由消息号及二进制内容构成的消息
func (defaultCodec *DefaultCodec) EncodeBytes(cmdid uint16,
	protodata []byte) *MessageBinary {
	return DefaultEncodeBytes(cmdid, protodata)
}
