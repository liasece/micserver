package pool

import (
	"sync"

	"github.com/liasece/micserver/util/hash"
	"github.com/liasece/micserver/util/strings"
)

type subPool struct {
	sync.Map
}

func (this *subPool) Len() int {
	res := 0
	this.Map.Range(func(k, v interface{}) bool {
		res++
		return true
	})
	return res
}

type MapPool struct {
	pool         subPool
	groupPool    sync.Map
	groupSum     int
	getGroupFunc func(k interface{}) int
}

func (this *MapPool) Init(groupSum int) {
	if groupSum == 0 || groupSum == 1 || this.getGroupFunc == nil {
		// 不分组
		this.groupSum = 1
	} else {
		// 分组
		this.groupSum = groupSum
		for i := 0; i < this.groupSum; i++ {
			this.groupPool.Store(i, &subPool{})
		}
	}
	this.getGroupFunc = this.getGroupIndexFunc
}

func (this *MapPool) getGroupIndexFunc(k interface{}) int {
	str := strings.MustInterfaceToString(k)
	hash := hash.GetStringHash(str)
	return int(hash) % this.groupSum
}

func (this *MapPool) GetPool(k interface{}) *subPool {
	if this.groupSum > 1 && this.getGroupFunc != nil {
		index := this.getGroupFunc(k)
		vi, ok := this.groupPool.Load(index)
		if !ok {
			return nil
		}
		return vi.(*subPool)
	} else {
		return &this.pool
	}
}

func (this *MapPool) Load(k interface{}) (interface{}, bool) {
	return this.GetPool(k).Load(k)
}

func (this *MapPool) Store(k, v interface{}) {
	this.GetPool(k).Store(k, v)
}

func (this *MapPool) LoadOrStore(k, v interface{}) (interface{}, bool) {
	return this.GetPool(k).LoadOrStore(k, v)
}

func (this *MapPool) Delete(k interface{}) {
	this.GetPool(k).Delete(k)
}

func (this *MapPool) Range(cb func(k, v interface{}) bool) {
	if this.groupSum > 1 && this.getGroupFunc != nil {
		goon := true
		this.groupPool.Range(func(_, pooli interface{}) bool {
			pool := pooli.(*subPool)
			pool.Range(func(k, v interface{}) bool {
				res := cb(k, v)
				if !res {
					goon = false
				}
				return res
			})
			return goon
		})
	} else {
		this.pool.Range(cb)
	}
}

func (this *MapPool) LenTotal() int {
	res := 0
	this.groupPool.Range(func(_, pooli interface{}) bool {
		pool := pooli.(*subPool)
		res += pool.Len()
		return true
	})
	return res
}
