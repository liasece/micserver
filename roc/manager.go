package roc

import (
	"sync"
)

// Manager ROC 管理器，每个 Module 都包含了一个ROC管理器，管理了本 Module 已注册的所有ROC类型及ROC对象。
type Manager struct {
	rocs      sync.Map
	eventHook IROCObjEventHook
}

// NewROC 新建一种类型的ROC
func (rocManager *Manager) NewROC(objtype ObjType) *ROC {
	var res *ROC
	newroc := &ROC{}
	vi, isLoad := rocManager.rocs.LoadOrStore(objtype, newroc)
	if !isLoad {
		newroc.Init()
		newroc.HookObjEvent(rocManager)
		res = newroc
	} else {
		res = vi.(*ROC)
	}
	return res
}

// HookObjEvent 监听ROC事件
func (rocManager *Manager) HookObjEvent(hook IROCObjEventHook) {
	rocManager.eventHook = hook
}

// OnROCObjAdd 当增加一个ROC对象时调用
func (rocManager *Manager) OnROCObjAdd(obj IObj) {
	if rocManager.eventHook != nil {
		rocManager.eventHook.OnROCObjAdd(obj)
	}
}

// OnROCObjDel 当删除一个ROC对象时调用
func (rocManager *Manager) OnROCObjDel(obj IObj) {
	if rocManager.eventHook != nil {
		rocManager.eventHook.OnROCObjDel(obj)
	}
}

// RegObj 注册一个ROC对象，如果该类型的ROC未注册将返回 error
func (rocManager *Manager) RegObj(obj IObj) error {
	objtype := obj.GetROCObjType()
	roc := rocManager.GetROC(objtype)
	if roc == nil {
		return ErrUnregisterROC
	}
	return roc.RegObj(obj)
}

// GetROC 获取一个类型ROC
func (rocManager *Manager) GetROC(objtype ObjType) *ROC {
	vi, ok := rocManager.rocs.Load(objtype)
	if !ok {
		return nil
	}
	return vi.(*ROC)
}

// CallPathDecode 解码ROC调用路径
func (rocManager *Manager) CallPathDecode(kstr string) (ObjType, string) {
	return kstrDecode(kstr)
}

// kstr的格式必须为 ROC 远程对象调用那样定义的格式
// getObj (对象类型)([对象的键])
func (rocManager *Manager) getObj(objType ObjType, objID string) (IObj, bool) {
	roc := rocManager.GetROC(objType)
	if roc == nil {
		return nil, false
	}
	return roc.GetObj(objID)
}

// GetObj kstr的格式必须为 ROC 远程对象调用那样定义的格式
// (对象类型)([对象的键])
func (rocManager *Manager) GetObj(objType ObjType, objID string) (IObj, bool) {
	return rocManager.getObj(objType, objID)
}

// Call 执行远程发来的ROC调用请求
func (rocManager *Manager) Call(callstr string, arg []byte) ([]byte, error) {
	path := NewPath(callstr)
	obj, ok := rocManager.getObj(path.GetObjType(), path.GetObjID())
	if !ok || obj == nil {
		path.Reset()
		return nil, ErrUnknownObj
	}
	path.Reset()
	return obj.OnROCCall(path, arg)
}
