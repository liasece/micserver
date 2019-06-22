package util

import (
	"errors"
	"sync"
)

// 粒度控制 单位：字节
var sizeControl []int = []int{32, 64, 128, 256, 512, 1024, 2 * 1024,
	4 * 1024, 8 * 1024, 16 * 1024, 32 * 1024, 64 * 1024, 128 * 1024,
	256 * 1024, 512 * 1024, 1024 * 1024, 2 * 1024 * 1024, 4 * 1024 * 1024,
	8 * 1024 * 1024, 16 * 1024 * 1024, 32 * 1024 * 1024, 64 * 1024 * 1024,
	128 * 1024 * 1024, 256 * 1024 * 1024, 512 * 1024 * 1024,
	1024 * 1024 * 1024}

type FlexiblePool struct {
	pools         []sync.Pool
	sizeControler []int
	New           func(int) interface{}
}

func (this *FlexiblePool) MaxSize() int {
	return this.sizeControler[len(this.sizeControler)-1]
}

// 获取某大小所在的区间下标，容忍最大值模式
func (this *FlexiblePool) getIndex(size int) int {
	if size <= this.sizeControler[0] {
		return 0
	}
	controlSize := len(this.sizeControler)
	for i := 1; i < controlSize; i++ {
		if size > this.sizeControler[i-1] && size <= this.sizeControler[i] {
			// 大小介于 (i-1,i]
			return i
		}
	}
	return -1
}

// 获取某大小所在的区间下标，保证安全
func (this *FlexiblePool) getIndexSafety(size int) int {
	if size < this.sizeControler[0] {
		return -1
	}
	controlSize := len(this.sizeControler)
	for i := 1; i < controlSize; i++ {
		if size >= this.sizeControler[i-1] && size < this.sizeControler[i] {
			// 大小介于 (i-1,i]
			return i - 1
		}
	}
	return controlSize - 1
}

func (this *FlexiblePool) Get(size int) (interface{}, error) {
	if size > this.MaxSize() {
		return nil, errors.New("Cap size out of maximum.")
	}
	index := this.getIndex(size)
	if index < 0 {
		return nil, errors.New("SizeControl index out of range.")
	}
	res := this.pools[index].Get()
	if res == nil {
		res = this.New(this.sizeControler[index])
	}
	return res, nil
}

func (this *FlexiblePool) Put(data interface{}, size int) error {
	if size > this.MaxSize() {
		return errors.New("Cap size out of maximum.")
	}
	index := this.getIndexSafety(size)
	if index < 0 {
		return errors.New("SizeControl index out of range.")
	}
	this.pools[index].Put(data)
	return nil
}

func NewFlexiblePool(sizeControler []int,
	Newer func(int) interface{}) *FlexiblePool {
	if sizeControler == nil {
		sizeControler = sizeControl
	}
	res := &FlexiblePool{}
	res.New = Newer
	// 构造大小控制列表
	res.sizeControler = make([]int, len(sizeControler))
	copy(res.sizeControler, sizeControler)
	// 构造sync.Pool列表
	res.pools = make([]sync.Pool, len(res.sizeControler))
	return res
}
