package msg

import (
	"encoding/binary"
	"fmt"
)

const (
	MSG_HEADSIZE_L1 = (4 + 1 + 1)
	MSG_HEADSIZE_L2 = (2 + 4 + 4)
	MSG_HEADSIZE    = (MSG_HEADSIZE_L1 + MSG_HEADSIZE_L2)
)

// layer 1
type MessageBinaryHeadL1 struct {
	CmdLen  uint32          // 4  消息长度
	CmdMask TEncryptionType // 1  消息是否加密
	CmdZip  byte            // 1  是否压缩
}

func (this *MessageBinaryHeadL1) WriteToBuffer(data []byte) (size int) {
	binary.BigEndian.PutUint32(data[size:], this.CmdLen) // 4
	size += 4
	data[size] = byte(this.CmdMask)
	size += 1
	data[size] = this.CmdZip
	size += 1
	return
}

func (this *MessageBinaryHeadL1) ReadFromBuffer(data []byte) (int, error) {
	if len(data) < MSG_HEADSIZE_L1 {
		return 0, fmt.Errorf("data not enough")
	}
	this.CmdLen = binary.BigEndian.Uint32(data[0:4])
	this.CmdMask = TEncryptionType(data[4])
	this.CmdZip = data[5]
	// cmdlen must  include head layer 1 size
	if int(this.CmdLen) < MSG_HEADSIZE_L1 {
		return 0, fmt.Errorf("error cmdlen")
	}
	return MSG_HEADSIZE_L1, nil
}

func (this *MessageBinaryHeadL1) LowerSize() int {
	return int(this.CmdLen) - MSG_HEADSIZE_L1
}

// layer 2
type MessageBinaryHeadL2 struct {
	CmdID     uint16 // 2
	TimeStamp uint32 // 4
	DataLen   uint32 // 4
}

func (this *MessageBinaryHeadL2) WriteToBuffer(data []byte) (size int) {
	binary.BigEndian.PutUint16(data[size:], this.CmdID) // 2
	size += 2
	binary.BigEndian.PutUint32(data[size:], this.TimeStamp) // 4
	size += 4
	binary.BigEndian.PutUint32(data[size:], this.DataLen) // 4
	size += 4
	return
}

func (this *MessageBinaryHeadL2) ReadFromBuffer(data []byte) (int, error) {
	if len(data) < MSG_HEADSIZE_L2 {
		return 0, fmt.Errorf("data not enough")
	}
	size := 0
	this.CmdID = binary.BigEndian.Uint16(data[size:])
	size += 2
	this.TimeStamp = binary.BigEndian.Uint32(data[size:])
	size += 4
	this.DataLen = binary.BigEndian.Uint32(data[size:])
	size += 4
	return MSG_HEADSIZE_L2, nil
}

func (this *MessageBinaryHeadL2) LowerSize() int {
	return int(this.DataLen)
}
