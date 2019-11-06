package roc

import (
	"fmt"
	"strings"
	"sync"
)

type ROCManager struct {
	rocs      sync.Map
	eventHook IROCObjEventHook
}

func (this *ROCManager) NewObjectType(objtype string) {
	newroc := &ROC{}
	_, isLoad := this.rocs.LoadOrStore(objtype, newroc)
	if !isLoad {
		newroc.Init()
		newroc.HookObjEvent(this)
	}
}

func (this *ROCManager) HookObjEvent(hook IROCObjEventHook) {
	this.eventHook = hook
}

func (this *ROCManager) OnROCObjAdd(obj IObj) {
	if this.eventHook != nil {
		this.eventHook.OnROCObjAdd(obj)
	}
}

func (this *ROCManager) OnROCObjDel(obj IObj) {
	if this.eventHook != nil {
		this.eventHook.OnROCObjDel(obj)
	}
}

func (this *ROCManager) RegObj(obj IObj) error {
	objtype := obj.GetObjType()
	roc := this.GetROC(objtype)
	if roc == nil {
		return ErrUnregisterRoc
	}
	return roc.RegObj(obj)
}

func (this *ROCManager) GetROC(objtype string) *ROC {
	vi, ok := this.rocs.Load(objtype)
	if !ok {
		return nil
	}
	return vi.(*ROC)
}

func (this *ROCManager) CallPathDecode(kstr string) (string, string) {
	return kstrDecode(kstr)
}

// kstr的格式必须为 ROC 远程对象调用那样定义的格式
// (对象类型)([对象的键])
func (this *ROCManager) getObj(objType string, objID string) (IObj, bool) {
	roc := this.GetROC(objType)
	if roc == nil {
		return nil, false
	}
	return roc.GetObj(objID)
}

// kstr的格式必须为 ROC 远程对象调用那样定义的格式
// (对象类型)([对象的键])
func (this *ROCManager) GetObj(objType string, objID string) (IObj, bool) {
	return this.getObj(objType, objID)
}

func (this *ROCManager) Call(callstr string, arg []byte) ([]byte, error) {
	strs := strings.Split(callstr, ".")
	if len(strs) < 1 {
		return nil, fmt.Errorf("callstr split on . is length == 0")
	}
	path := NewROCPath(strs)
	obj, ok := this.getObj(path.GetObjType(), path.GetObjID())
	if !ok || obj == nil {
		path.Reset()
		return nil, fmt.Errorf("has't this object:%+v,%+v", path, path.Move())
	}
	path.Reset()
	return obj.ROCCall(path, arg)
}
