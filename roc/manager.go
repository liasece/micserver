package roc

import (
	"sync"
)

// ROC 管理器，每个 Module 都包含了一个ROC管理器，管理了本 Module 已注册的所有ROC类型
// 及ROC对象。
type ROCManager struct {
	rocs      sync.Map
	eventHook IROCObjEventHook
}

// 新建一种类型的ROC
func (this *ROCManager) NewROC(objtype ROCObjType) *ROC {
	var res *ROC
	newroc := &ROC{}
	vi, isLoad := this.rocs.LoadOrStore(objtype, newroc)
	if !isLoad {
		newroc.Init()
		newroc.HookObjEvent(this)
		res = newroc
	} else {
		res = vi.(*ROC)
	}
	return res
}

// 监听ROC事件
func (this *ROCManager) HookObjEvent(hook IROCObjEventHook) {
	this.eventHook = hook
}

// 当增加一个ROC对象时调用
func (this *ROCManager) OnROCObjAdd(obj IObj) {
	if this.eventHook != nil {
		this.eventHook.OnROCObjAdd(obj)
	}
}

// 当删除一个ROC对象时调用
func (this *ROCManager) OnROCObjDel(obj IObj) {
	if this.eventHook != nil {
		this.eventHook.OnROCObjDel(obj)
	}
}

// 注册一个ROC对象，如果该类型的ROC未注册将返回 error
func (this *ROCManager) RegObj(obj IObj) error {
	objtype := obj.GetROCObjType()
	roc := this.GetROC(objtype)
	if roc == nil {
		return ErrUnregisterROC
	}
	return roc.RegObj(obj)
}

// 获取一个类型ROC
func (this *ROCManager) GetROC(objtype ROCObjType) *ROC {
	vi, ok := this.rocs.Load(objtype)
	if !ok {
		return nil
	}
	return vi.(*ROC)
}

// 解码ROC调用路径
func (this *ROCManager) CallPathDecode(kstr string) (ROCObjType, string) {
	return kstrDecode(kstr)
}

// kstr的格式必须为 ROC 远程对象调用那样定义的格式
// (对象类型)([对象的键])
func (this *ROCManager) getObj(objType ROCObjType, objID string) (IObj, bool) {
	roc := this.GetROC(objType)
	if roc == nil {
		return nil, false
	}
	return roc.GetObj(objID)
}

// kstr的格式必须为 ROC 远程对象调用那样定义的格式
// (对象类型)([对象的键])
func (this *ROCManager) GetObj(objType ROCObjType, objID string) (IObj, bool) {
	return this.getObj(objType, objID)
}

// 执行远程发来的ROC调用请求
func (this *ROCManager) Call(callstr string, arg []byte) ([]byte, error) {
	path := NewROCPath(callstr)
	obj, ok := this.getObj(path.GetObjType(), path.GetObjID())
	if !ok || obj == nil {
		path.Reset()
		return nil, ErrUnknowObj
	}
	path.Reset()
	return obj.OnROCCall(path, arg)
}
