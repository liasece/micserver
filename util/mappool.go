package util

import (
	"sync"
)

// 用户组，用于优化锁
type GroupItem struct {
	datas  sync.Map
	length uint32
}

// 遍历组内所有用户
func (this *GroupItem) Range(
	f func(interface{}, interface{}) bool) {
	this.datas.Range(f)
}

func (this *GroupItem) Len() uint32 {
	return this.length
}

func (this *GroupItem) Push(key interface{}, value interface{}) {
	this.add(key, value)
}

func (this *GroupItem) Pop(key interface{}) {
	this.remove(key)
}

func (this *GroupItem) Get(key interface{}) (interface{}, bool) {
	if value, found := this.datas.Load(key); found {
		// 在列表中找到了 User 对象
		return value, true
	}
	return nil, false
}

func (this *GroupItem) LoadOfStroe(key interface{}, value interface{}) (interface{}, bool) {
	vi, isLoad := this.datas.LoadOrStore(key, value)
	if !isLoad {
		this.length++
	} else {
		this.datas.Store(key, value)
	}
	return vi, isLoad
}

func (this *GroupItem) remove(key interface{}) {
	if _, ok := this.datas.Load(key); !ok {
		return
	}
	// 删除
	this.datas.Delete(key)
	this.length--
}

func (this *GroupItem) add(key interface{}, value interface{}) {
	if _, ok := this.datas.Load(key); !ok {
		// 是新增的
		this.length++
	}
	this.datas.Store(key, value)
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

// 遍历组内所有
func (this *MapPool) Range(groupIndex uint32,
	f func(interface{}, interface{}) bool) {
	this.groups[groupIndex].Range(f)
}

func (this *MapPool) LenTotal() uint32 {
	res := uint32(0)
	for _, v := range this.groups {
		res += v.Len()
	}
	return res
}

func (this *MapPool) Len(groupIndex uint32) uint32 {
	return this.groups[groupIndex].Len()
}

func (this *MapPool) Push(groupIndex uint32,
	key interface{}, value interface{}) {
	this.groups[groupIndex].Push(key, value)
}

func (this *MapPool) Pop(groupIndex uint32,
	key interface{}) {
	this.groups[groupIndex].Pop(key)
}

func (this *MapPool) Get(groupIndex uint32,
	key interface{}) (interface{}, bool) {
	return this.groups[groupIndex].Get(key)
}

func (this *MapPool) LoadOfStroe(groupIndex uint32,
	key interface{}, value interface{}) (interface{}, bool) {
	return this.groups[groupIndex].LoadOfStroe(key, value)
}

func NewMapPool(groupSum uint32) *MapPool {
	res := &MapPool{}
	res.Init(groupSum)
	return res
}
