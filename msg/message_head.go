package msg

import (
	"encoding/binary"
	"fmt"
)

const (
	MSG_HEADSIZE_L1 = (4 + 2)
	MSG_HEADSIZE    = (MSG_HEADSIZE_L1)
)

// layer 1
type MessageBinaryHeadL1 struct {
	CmdLen uint32 // 4  消息长度
	CmdID  uint16 // 2

	// 不存在于二进制数据中，由 CmdLen - MSG_HEADSIZE 得到
	DataLen uint32
}

func (this *MessageBinaryHeadL1) WriteToBuffer(data []byte) (size int) {
	binary.BigEndian.PutUint32(data[size:], this.CmdLen) // 4
	size += 4
	binary.BigEndian.PutUint16(data[size:], this.CmdID) // 2
	size += 2

	return
}

func (this *MessageBinaryHeadL1) ReadFromBuffer(data []byte) (int, error) {
	if len(data) < MSG_HEADSIZE_L1 {
		return 0, fmt.Errorf("data not enough")
	}
	size := 0
	this.CmdLen = binary.BigEndian.Uint32(data[size:])
	this.DataLen = this.CmdLen - MSG_HEADSIZE
	size += 4
	this.CmdID = binary.BigEndian.Uint16(data[size:])
	size += 2
	// cmdlen must  include head layer 1 size
	if int(this.CmdLen) < MSG_HEADSIZE_L1 {
		return 0, fmt.Errorf("error cmdlen")
	}
	return MSG_HEADSIZE_L1, nil
}

func (this *MessageBinaryHeadL1) LowerSize() int {
	return int(this.CmdLen) - MSG_HEADSIZE_L1
}
