// Package roc Remote Object Call.
// ROC ，是 micserver 重要的分布式目标调用思想。
// 如果一个对象，例如房间/商品/玩家/工会，需要提供一个可供远程执行的操作，
// 这在以前称之为 RPC 调用，在 micserver 中，任意一个构造了这种对象的模块，
// 均可以通过 BaseModule.GetROC(objtype).GetOrRegObj(obj) 来
// 注册一个 ROC 对象，在其他模块中，只需要通过 BaseModule.ROCCallNR 等方法，提供
// 目标对象的类型及ID，即可发起针对该对象的远程操作。
// 因此，在任意模块中，发起的任意针对其他模块的调用，都不应该使用模块ID来操作，
// 因为使用统一的 ROC 至少包含以下好处：
// 		无需知道目标对象在哪个模块上；
// 		只需要关心目标对象的ID（目标的类型你当然是知道的）；
// 		在模块更新时，可以统一将该模块的 ROC 对象迁移到新版本模块中实现热更；
// 		可以将对象存储到数据库并且在其他模块中加载（基于第一点好处）；
// 		对象的位置/调用路由等由底层系统维护，可提供一个高可用/强一致/易维护的分布式网络。
package roc

import (
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/util/pool"
)

// ObjType ROC 对象的类型，本质上是字符串
type ObjType string

// ROC 对象分组数
const (
	PoolGroupSum = 8
)

// ROC 一个类型的ROC，维护了这个类型的所有 ROC 对象
type ROC struct {
	objPool   pool.MapPool
	eventHook IROCObjEventHook
}

// Init 初始化该类型的ROC
func (roc *ROC) Init() {
	roc.objPool.Init(PoolGroupSum)
}

// HookObjEvent 设置消息
func (roc *ROC) HookObjEvent(hook IROCObjEventHook) {
	roc.eventHook = hook
}

func (roc *ROC) onRegObj(obj IObj) {
	if roc.eventHook != nil {
		roc.eventHook.OnROCObjAdd(obj)
	}
}

func (roc *ROC) onDelObj(obj IObj) {
	if roc.eventHook != nil {
		roc.eventHook.OnROCObjDel(obj)
	}
}

// RegObj 在使用远程对象调用前，需要先注册
func (roc *ROC) RegObj(obj IObj) error {
	id := obj.GetROCObjID()
	roc.objPool.Store(id, obj)
	roc.onRegObj(obj)
	return nil
}

// DelObj 删除目标ROC对象
func (roc *ROC) DelObj(obj IObj) error {
	id := obj.GetROCObjID()
	return roc.DelObjByID(id)
}

// DelObjByID 删除指定ID的ROC对象
func (roc *ROC) DelObjByID(id string) error {
	localObj, ok := roc.GetObj(id)
	roc.objPool.Delete(id)
	if ok && localObj != nil {
		roc.onDelObj(localObj)
	}
	return nil
}

// GetObj 获取指定ID的ROC对象
func (roc *ROC) GetObj(id string) (IObj, bool) {
	vi, ok := roc.objPool.Load(id)
	if !ok || vi == nil {
		return nil, ok
	}
	return vi.(IObj), ok
}

// GetOrRegObj 获取或者注册一个ROC对象
func (roc *ROC) GetOrRegObj(id string, obj IObj) (IObj, bool) {
	vi, isLoad := roc.objPool.LoadOrStore(id, obj)
	if !isLoad {
		roc.onRegObj(obj)
	}
	if vi == nil {
		return nil, isLoad
	}
	return vi.(IObj), isLoad
}

// RangeObj 遍历该类型的ROC对象
func (roc *ROC) RangeObj(f func(obj IObj) bool) {
	roc.objPool.Range(func(ki, vi interface{}) bool {
		v, ok := vi.(IObj)
		if !ok {
			log.Error("interface conversion: %+v is not roc.IObj", vi)
			return true
		}
		return f(v)
	})
}
