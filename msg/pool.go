package msg

import (
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/util/pool"
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
var pools *pool.FlexiblePool

// 初始化灵活对象池
func init() {
	pools = pool.NewFlexiblePool(sizeControl, newMsgBinaryBySize)
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

// 通过二进制流创建 MessageBinary
func GetByBytes(cmdid uint16, protodata []byte) *MessageBinary {
	// 获取基础数据
	datalen := uint32(len(protodata))
	totalLength := uint32(MSG_HEADSIZE + datalen)
	// 判断数据合法性
	if totalLength >= MessageMaxSize {
		log.Error("[GetByBytes] "+
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
		log.Error("[GetByBytes] "+
			"无法分配MsgBinary的内存！！！ CmdID[%d] Len[%d]",
			cmdid, totalLength)
		return nil
	}
	// 将 protodata 拷贝至 buffer 的数据域
	copy(msgbinary.buffer[MSG_HEADSIZE:totalLength], protodata)

	// 初始化消息信息

	// MessageBinaryBody
	// 消息数据字段指针指向 buffer 数据域
	msgbinary.ProtoData =
		msgbinary.buffer[MSG_HEADSIZE:totalLength]

	// MessageBinaryHeadL1
	msgbinary.MessageBinaryHeadL1.CmdLen = totalLength
	msgbinary.MessageBinaryHeadL1.CmdID = cmdid

	// 将结构数据填入 buffer
	msgbinary.writeHeadBuffer()

	return msgbinary
}

// 通过结构体创建 MessageBinary
func GetByObj(v MsgStruct) *MessageBinary {
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
		log.Error("[GetByObj] "+
			"无法分配MsgBinary的内存！！！ CmdLen[%d] DataLen[%d]",
			totalLength, datalen)
		return nil
	}

	// MessageBinaryBody
	// 将 protodata 拷贝至 buffer 的数据域
	v.WriteBinary(msgbinary.buffer[MSG_HEADSIZE:totalLength])

	// 初始化消息信息

	// 消息数据字段指针指向 buffer 数据域
	msgbinary.ProtoData =
		msgbinary.buffer[MSG_HEADSIZE:totalLength]

	// 初始化消息信息
	// MessageBinaryHeadL1
	msgbinary.MessageBinaryHeadL1.CmdLen = totalLength
	msgbinary.MessageBinaryHeadL1.CmdID = cmdid

	// 将结构数据填入 buffer
	msgbinary.writeHeadBuffer()

	return msgbinary
}
