// Remote Project Call
package roc

import (
	"github.com/liasece/micserver/util/pool"
)

type ROCObjType string

const (
	ROC_POOL_GROUP_SUM = 8
)

type ROC struct {
	objPool   pool.MapPool
	eventHook IROCObjEventHook
}

func (this *ROC) Init() {
	this.objPool.Init(ROC_POOL_GROUP_SUM)
}

func (this *ROC) HookObjEvent(hook IROCObjEventHook) {
	this.eventHook = hook
}

func (this *ROC) onRegObj(obj IObj) {
	if this.eventHook != nil {
		this.eventHook.OnROCObjAdd(obj)
	}
}

func (this *ROC) onDelObj(obj IObj) {
	if this.eventHook != nil {
		this.eventHook.OnROCObjDel(obj)
	}
}

// 在使用远程对象调用前，需要先注册
func (this *ROC) RegObj(obj IObj) error {
	id := obj.GetROCObjID()
	this.objPool.Store(id, obj)
	this.onRegObj(obj)
	return nil
}

// 删除远程调用对象
func (this *ROC) DelObj(obj IObj) error {
	id := obj.GetROCObjID()
	return this.DelObjByID(id)
}

func (this *ROC) DelObjByID(id string) error {
	localObj, ok := this.GetObj(id)
	this.objPool.Delete(id)
	if ok && localObj != nil {
		this.onDelObj(localObj)
	}
	return nil
}

func (this *ROC) GetObj(id string) (IObj, bool) {
	vi, ok := this.objPool.Load(id)
	if !ok || vi == nil {
		return nil, ok
	}
	return vi.(IObj), ok
}

func (this *ROC) GetOrRegObj(id string, obj IObj) (IObj, bool) {
	vi, isLoad := this.objPool.LoadOrStore(id, obj)
	if !isLoad {
		this.onRegObj(obj)
	}
	if vi == nil {
		return nil, isLoad
	}
	return vi.(IObj), isLoad
}

func (this *ROC) RangeObj(f func(obj IObj) bool) {
	this.objPool.Range(func(ki, vi interface{}) bool {
		v := vi.(IObj)
		return f(v)
	})
}
