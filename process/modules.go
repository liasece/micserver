// Package process 所有注册于本进程的 Module 都会维护在本单例中。
package process

import (
	"sync"

	"github.com/liasece/micserver/base"
)

var (
	_gModules sync.Map
)

// AddModule 增加一个本进程的 Module
func AddModule(module base.IModule) {
	_gModules.Store(module.GetModuleID(), module)
}

// HasModule 判断目标 Module 是否在本进程中
func HasModule(moduleID string) bool {
	_, ok := _gModules.Load(moduleID)
	return ok
}
