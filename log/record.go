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
	level    int32
	timeUnix int64
}

// String 格式化该日志记录字符串
func (r *Record) String() string {
	return fmt.Sprintf("%s [%s] %s %s\n", r.time, r.name, LEVELFLAGS[r.level], r.info)
}
