package process

import (
	"sync"

	"github.com/liasece/micserver/base"
)

var (
	_gModules sync.Map
)

func AddModule(module base.IModule) {
	_gModules.Store(module.GetModuleID(), module)
}

func HasModule(moduleID string) bool {
	res := false
	_gModules.Range(func(ki, vi interface{}) bool {
		k := ki.(string)
		if k == moduleID {
			res = true
			return false
		}
		return true
	})
	return res
}
