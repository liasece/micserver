package roc

import (
	"fmt"
	"strings"
	"sync"
)

type ROCManager struct {
	rocs sync.Map
}

func (this *ROCManager) NewObjectType(objtype string) {
	newroc := &ROC{}
	_, isnew := this.rocs.LoadOrStore(objtype, newroc)
	if isnew {
		newroc.Init()
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

// kstr的格式必须为 ROC 远程对象调用那样定义的格式，如：
// 对象类型[对象的键]
func (this *ROCManager) kstrDecode(kstr string) (string, string) {
	t := ""
	key := ""
	inkey := false
	for _, k := range kstr {
		if k == '[' {
			inkey = true
		} else if k == ']' {
			inkey = false
		} else {
			if key == "" && !inkey {
				t = t + fmt.Sprintf("%c", k)
			} else if t != "" && inkey {
				key = key + fmt.Sprintf("%c", k)
			} else {
				return "", ""
			}
		}
	}
	return t, key
}

// kstr的格式必须为 ROC 远程对象调用那样定义的格式
// (对象类型)([对象的键])
func (this *ROCManager) getObj(kstr string) (IObj, bool) {
	t, k := this.kstrDecode(kstr)
	roc := this.GetROC(t)
	if roc == nil {
		return nil, false
	}
	return roc.GetObj(k)
}

// kstr的格式必须为 ROC 远程对象调用那样定义的格式
// (对象类型)([对象的键])
func (this *ROCManager) GetObj(kstr string) (IObj, bool) {
	return this.getObj(kstr)
}

func (this *ROCManager) Call(callstr string, arg []byte) {
	strs := strings.Split(callstr, ".")
	if len(strs) < 1 {
		return
	}
	obj, ok := this.getObj(strs[0])
	if !ok || obj == nil {
		return
	}
	obj.ROCCall(strs[1:], arg)
}
