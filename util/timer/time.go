/*
Package timer 系统时间
*/
package timer

import (
	"time"
)

// GetTimeMs 获取系统毫秒时间
func GetTimeMs() uint64 {
	return uint64(time.Now().UnixNano()) / 1000000
}
