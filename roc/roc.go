// Remote Project Call
package roc

import (
	"github.com/liasece/micserver/util/pool"
)

const (
	ROC_POOL_GROUP_SUM = 8
)

type ROC struct {
	objPool   pool.MapPool
	onfRegObj func(IObj)
	onfDelObj func(IObj)
}

func (this *ROC) Init() {
	this.objPool.Init(ROC_POOL_GROUP_SUM)
}

func (this *ROC) RegOnRegObj(cb func(IObj)) {
	this.onfRegObj = cb
}

func (this *ROC) onRegObj(obj IObj) {
	if this.onfRegObj != nil {
		this.onfRegObj(obj)
	}
}

func (this *ROC) RegOnDelObj(cb func(IObj)) {
	this.onfDelObj = cb
}

func (this *ROC) onDelObj(obj IObj) {
	if this.onfDelObj != nil {
		this.onfDelObj(obj)
	}
}

// 在使用远程对象调用前，需要先注册
func (this *ROC) RegObj(obj IObj) error {
	id := obj.GetObjID()
	this.objPool.Store(id, obj)
	this.onRegObj(obj)
	return nil
}

// 删除远程调用对象
func (this *ROC) DelObj(obj IObj) error {
	id := obj.GetObjID()
	this.objPool.Delete(id)
	this.onDelObj(obj)
	return nil
}

func (this *ROC) GetObj(id string) (IObj, bool) {
	vi, ok := this.objPool.Load(id)
	if !ok || vi == nil {
		return nil, ok
	}
	return vi.(IObj), ok
}

func (this *ROC) GetOrRegBoj(id string, obj IObj) (IObj, bool) {
	vi, isLoad := this.objPool.LoadOrStore(id, obj)
	if !isLoad {
		this.onRegObj(obj)
	}
	if vi == nil {
		return nil, isLoad
	}
	return vi.(IObj), isLoad
}
