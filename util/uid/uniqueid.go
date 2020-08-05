/*
Package uid UUID生成器
*/
package uid

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type uniqueIDBuilder struct {
	lastLowlevelID uint16
	lastMidlevelID uint32
	// lastHeightlevelID uint16

	lastMidLowLevelID uint64
	mutex             sync.Mutex
}

var unique uniqueIDBuilder

// GenUniqueID 生成一个至少保证当前模块唯一的ID
func (builder *uniqueIDBuilder) GenUniqueID(heightlevelID uint16) (string, error) {
	builder.mutex.Lock()
	defer builder.mutex.Unlock()
	nowtime := time.Now().Unix() - 1514736000
	if nowtime <= 0 {
		return "", errors.New("Server time error!!! Must late than 2018/1/1 00:00:00")
	}
	now := uint32(nowtime & 0x0ffffffff)
	if now > builder.lastMidlevelID {
		// 时间已过去最后的兼容秒数
		builder.lastLowlevelID = 0
		builder.lastMidlevelID = now
	} else {
		if builder.lastLowlevelID == 0x0ffff {
			// 生成频率超限
			return "", errors.New("本秒内已随机出超出限制的唯一ID数量")
		}
		builder.lastLowlevelID++
	}

	subvalue := uint64(0)
	subvalue |= uint64(builder.lastMidlevelID) << (16)
	subvalue |= uint64(builder.lastLowlevelID) << (0)
	if subvalue <= builder.lastMidLowLevelID {
		return "", errors.New("生成的ID可能重复了")
	}
	builder.lastMidLowLevelID = subvalue

	res := uint64(0)
	res |= uint64(subvalue) << (16)
	res |= uint64(heightlevelID) << (0)
	return fmt.Sprint(res), nil
}

// GenUniqueID 生成一个至少保证当前模块唯一的ID
func GenUniqueID(heightlevelID uint16) (string, error) {
	return unique.GenUniqueID(heightlevelID)
}
