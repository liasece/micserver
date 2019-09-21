// Remote Project Call
package roc

import (
	"github.com/liasece/micserver/util"
)

const (
	POOL_GROUP_SUM = 8
)

type ROC struct {
	objPool util.MapPool
}

func (this *ROC) Init() {
	this.objPool.Init(POOL_GROUP_SUM)
}

// 在使用远程对象调用前，需要先注册
func (this *ROC) RegObj(obj IObj) error {
	id := obj.GetObjID()
	this.objPool.Store(id, obj)
	return nil
}

func (this *ROC) GetObj(id string) (IObj, bool) {
	vi, ok := this.objPool.Load(id)
	if !ok || vi == nil {
		return nil, ok
	}
	return vi.(IObj), ok
}
