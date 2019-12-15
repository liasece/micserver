package log

import (
	"fmt"
	"sync"
)

var (
	recordPool *sync.Pool
)

type Record struct {
	time     string
	name     string
	code     string
	info     string
	level    int32
	timeUnix int64
}

func (r *Record) String() string {
	return fmt.Sprintf("%s [%s] %s %s\n", r.time, r.name, LEVEL_FLAGS[r.level], r.info)
}
