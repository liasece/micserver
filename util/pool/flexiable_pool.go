package pool

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

// FlexiblePool Flexible Pool
type FlexiblePool struct {
	pools         []sync.Pool
	sizeControler []int
	New           func(int) interface{}
}

// MaxSize max pool item size
func (fp *FlexiblePool) MaxSize() int {
	return fp.sizeControler[len(fp.sizeControler)-1]
}

// getIndex 获取某大小所在的区间下标，容忍最大值模式
func (fp *FlexiblePool) getIndex(size int) int {
	if size <= fp.sizeControler[0] {
		return 0
	}
	controlSize := len(fp.sizeControler)
	for i := 1; i < controlSize; i++ {
		if size > fp.sizeControler[i-1] && size <= fp.sizeControler[i] {
			// 大小介于 (i-1,i]
			return i
		}
	}
	return -1
}

// getIndexSafety 获取某大小所在的区间下标，保证安全
func (fp *FlexiblePool) getIndexSafety(size int) int {
	if size < fp.sizeControler[0] {
		return -1
	}
	controlSize := len(fp.sizeControler)
	for i := 1; i < controlSize; i++ {
		if size >= fp.sizeControler[i-1] && size < fp.sizeControler[i] {
			// 大小介于 (i-1,i]
			return i - 1
		}
	}
	return controlSize - 1
}

// Get func
func (fp *FlexiblePool) Get(size int) (interface{}, error) {
	if size > fp.MaxSize() {
		return nil, errors.New("cap size out of maximum")
	}
	index := fp.getIndex(size)
	if index < 0 {
		return nil, errors.New("sizeControl index out of range")
	}
	res := fp.pools[index].Get()
	if res == nil {
		res = fp.New(fp.sizeControler[index])
	}
	return res, nil
}

// Put func
func (fp *FlexiblePool) Put(data interface{}, size int) error {
	if size > fp.MaxSize() {
		return errors.New("cap size out of maximum")
	}
	index := fp.getIndexSafety(size)
	if index < 0 {
		return errors.New("SizeControl index out of range")
	}
	fp.pools[index].Put(data)
	return nil
}

// NewFlexiblePool func
func NewFlexiblePool(sizeControler []int, Newer func(int) interface{}) *FlexiblePool {
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
