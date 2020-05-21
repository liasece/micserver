package monitor

import (
	"sync/atomic"
)

// Load 负载
type Load struct {
	value int64
}

// GetLoad func
func (load *Load) GetLoad() int64 {
	return atomic.LoadInt64(&load.value)
}

// AddLoad func
func (load *Load) AddLoad(add int64) {
	atomic.AddInt64(&load.value, add)
}

// SetLoad func
func (load *Load) SetLoad(value int64) {
	atomic.StoreInt64(&load.value, value)
}
