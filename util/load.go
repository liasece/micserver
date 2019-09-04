package util

import (
	"sync/atomic"
)

// 负载
type Load struct {
	value int64
}

func (this *Load) GetLoad() int64 {
	return atomic.LoadInt64(&this.value)
}

func (this *Load) AddLoad(add int64) {
	atomic.AddInt64(&this.value, add)
}

func (this *Load) SetLoad(value int64) {
	atomic.StoreInt64(&this.value, value)
}
