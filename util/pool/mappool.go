package pool

import (
	"sync"

	"github.com/liasece/micserver/util/hash"
	"github.com/liasece/micserver/util/strings"
)

// SubPool sub pool
type SubPool struct {
	sync.Map
}

// Len func
func (p *SubPool) Len() int {
	res := 0
	p.Map.Range(func(k, v interface{}) bool {
		res++
		return true
	})
	return res
}

// MapPool map pool
type MapPool struct {
	pool         SubPool
	groupPool    sync.Map
	groupSum     int
	getGroupFunc func(k interface{}) int
}

// Init func
func (p *MapPool) Init(groupSum int) {
	if groupSum == 0 || groupSum == 1 || p.getGroupFunc == nil {
		// 不分组
		p.groupSum = 1
	} else {
		// 分组
		p.groupSum = groupSum
		for i := 0; i < p.groupSum; i++ {
			p.groupPool.Store(i, &SubPool{})
		}
	}
	p.getGroupFunc = p.getGroupIndexFunc
}

func (p *MapPool) getGroupIndexFunc(k interface{}) int {
	str := strings.MustInterfaceToString(k)
	hash := hash.GetStringHash(str)
	return int(hash) % p.groupSum
}

// GetPool func
func (p *MapPool) GetPool(k interface{}) *SubPool {
	if p.groupSum > 1 && p.getGroupFunc != nil {
		index := p.getGroupFunc(k)
		vi, ok := p.groupPool.Load(index)
		if !ok {
			return nil
		}
		return vi.(*SubPool)
	}
	return &p.pool
}

// Load func
func (p *MapPool) Load(k interface{}) (interface{}, bool) {
	return p.GetPool(k).Load(k)
}

// Store func
func (p *MapPool) Store(k, v interface{}) {
	p.GetPool(k).Store(k, v)
}

// LoadOrStore func
func (p *MapPool) LoadOrStore(k, v interface{}) (interface{}, bool) {
	return p.GetPool(k).LoadOrStore(k, v)
}

// Delete func
func (p *MapPool) Delete(k interface{}) {
	p.GetPool(k).Delete(k)
}

// Range func
func (p *MapPool) Range(cb func(k, v interface{}) bool) {
	if p.groupSum > 1 && p.getGroupFunc != nil {
		goon := true
		p.groupPool.Range(func(_, pooli interface{}) bool {
			pool := pooli.(*SubPool)
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
		p.pool.Range(cb)
	}
}

// LenTotal func
func (p *MapPool) LenTotal() int {
	res := 0
	if p.groupSum > 1 && p.getGroupFunc != nil {
		p.groupPool.Range(func(_, pooli interface{}) bool {
			pool := pooli.(*SubPool)
			res += pool.Len()
			return true
		})
	} else {
		res += p.pool.Len()
	}
	return res
}
