// Remote Project Call
package roc

import (
	"github.com/liasece/micserver/util"
)

const (
	ROC_POOL_GROUP_SUM = 8
)

type ROC struct {
	objPool   util.MapPool
	catch     Catch
	onfRegObj func(IObj)
}

func (this *ROC) Init() {
	this.objPool.Init(ROC_POOL_GROUP_SUM)
}

func (this *ROC) RegOnRegObj(cb func(IObj)) {
	this.onfRegObj = cb
}

// 在使用远程对象调用前，需要先注册
func (this *ROC) RegObj(obj IObj) error {
	id := obj.GetObjID()
	this.objPool.Store(id, obj)
	if this.onfRegObj != nil {
		this.onfRegObj(obj)
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

func (this *ROC) GetOrRegBoj(id string, obj IObj) (IObj, bool) {
	vi, isLoad := this.objPool.LoadOrStore(id, obj)
	if !isLoad {
		if this.onfRegObj != nil {
			this.onfRegObj(obj)
		}
	}
	if vi == nil {
		return nil, isLoad
	}
	return vi.(IObj), isLoad
}
