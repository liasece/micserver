package util

import (
	"sync"
)

// 用户组，用于优化锁
type GroupItem struct {
	sync.Map
}

func (this *GroupItem) Len() int {
	res := 0
	this.Map.Range(func(k, v interface{}) bool {
		res++
		return true
	})
	return res
}

type MapPool struct {
	groups   []*GroupItem
	groupSum uint32
}

func (this *MapPool) Init(groupSum uint32) {
	this.groups = make([]*GroupItem, groupSum)
	for i := uint32(0); i < groupSum; i++ {
		this.groups[i] = &GroupItem{}
	}
	this.groupSum = groupSum
}

// 遍历所有
func (this *MapPool) RangeAll(f func(interface{}, interface{}) bool) {
	for _, v := range this.groups {
		v.Range(func(key interface{}, value interface{}) bool {
			return f(key, value)
		})
	}
}

func (this *MapPool) LenTotal() int {
	res := 0
	for _, v := range this.groups {
		res += v.Len()
	}
	return res
}

func (this *MapPool) GetGroupIndex(key interface{}) uint32 {
	if s, ok := key.(string); ok {
		return GetStringHash(s) % this.groupSum
	}
	return 0
}

func (this *MapPool) Len(key interface{}) int {
	groupIndex := this.GetGroupIndex(key)
	return this.groups[groupIndex].Len()
}

func (this *MapPool) Store(key interface{}, value interface{}) {
	groupIndex := this.GetGroupIndex(key)
	this.groups[groupIndex].Store(key, value)
}

func (this *MapPool) Delete(key interface{}) {
	groupIndex := this.GetGroupIndex(key)
	this.groups[groupIndex].Delete(key)
}

func (this *MapPool) Load(key interface{}) (interface{}, bool) {
	groupIndex := this.GetGroupIndex(key)
	return this.groups[groupIndex].Load(key)
}

func (this *MapPool) LoadOrStore(key interface{},
	value interface{}) (interface{}, bool) {
	groupIndex := this.GetGroupIndex(key)
	return this.groups[groupIndex].LoadOrStore(key, value)
}

func NewMapPool(groupSum uint32) *MapPool {
	res := &MapPool{}
	res.Init(groupSum)
	return res
}
