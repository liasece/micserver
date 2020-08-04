package log

import (
	"fmt"
	"sync"

	"github.com/liasece/micserver/log/core"
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
	fields   []core.Field
}

// String 格式化该日志记录字符串
func (r *Record) String() string {
	if r.name == "" {
		return fmt.Sprintf("%s %s %s", r.time, levelFlags[r.level], r.info)
	}
	return fmt.Sprintf("%s [%s] %s %s", r.time, r.name, levelFlags[r.level], r.info)
}
