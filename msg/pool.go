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
func GetMessageBinary(totalLength int) *MessageBinary {
	msg, err := pools.Get(totalLength)
	if err != nil {
		log.Error("[GetMessageBinary] TotalLen[%d] Err[%s]",
			totalLength, err.Error())
		return nil
	}
	if msg == nil {
		log.Error("[GetMessageBinary] nil return!!! TotalLen[%d]",
			totalLength)
		return nil
	}
	return msg.(*MessageBinary)
}
