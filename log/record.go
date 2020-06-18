package log

import (
	"fmt"
	"sync"
)

var (
	recordPool *sync.Pool
)

// Record 一条日志记录
type Record struct {
	time     string
	name     string
	code     string
	info     string
	level    Level
	timeUnix int64
}

// String 格式化该日志记录字符串
func (r *Record) String() string {
	if r.name == "" {
		return fmt.Sprintf("%s %s %s\n", r.time, levelFlags[r.level], r.info)
	}
	return fmt.Sprintf("%s [%s] %s %s\n", r.time, r.name, levelFlags[r.level], r.info)
}
