package msg

import (
	"fmt"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/util"
)

type MessageBinaryReader struct {
	inMsg  bool
	HeadL1 MessageBinaryHeadL1

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
		if !this.inMsg && this.netbuffer.Len() >= MSG_HEADSIZE {

			// 读头部4个字节
			MessageBinaryHeadL1buf, err := this.netbuffer.Read(0, MSG_HEADSIZE)
			if err != nil {
				return err
			}
			_, err = this.HeadL1.ReadBinary(MessageBinaryHeadL1buf)
			if err != nil {
				return fmt.Errorf("Head dec err:%s. headdata:%#v",
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

			// 获取合适大小的消息体
			msgbinary := getMessageBinaryByProtoDataLength(
				this.HeadL1.LowerSize())
			if msgbinary != nil {
				msgbinary.MessageBinaryHeadL1 = this.HeadL1
				// 解析消息（无6个字节的头）
				err := msgbinary.readBinaryNoHeadL1(cmdbuff)
				if err != nil {
					log.Error("[MessageBinaryReader.RangeMsgBinary] "+
						"解析消息错误 Err[%s] RecvLen[%d] HeadL1[%+v]",
						err.Error(), len(cmdbuff), this.HeadL1)
					return err
				} else {
					// 调用回调函数处理消息
					callback(msgbinary)
				}
			} else {
				log.Error("[MessageBinaryReader.RangeMsgBinary] "+
					"无法分配MsgBinary的内存！！！ RecvLen[%d] HeadL1[%+v]",
					len(cmdbuff), this.HeadL1)
			}
			// 退出消息处理状态
			this.inMsg = false
		} else {
			break
		}
	}
	return nil
}
