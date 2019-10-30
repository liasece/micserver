package process

import (
	"sync"

	"github.com/liasece/micserver/base"
)

var (
	_gModules     []base.IModule
	_gModuleMutex sync.Mutex
)

func AddModule(app base.IModule) {
	_gModuleMutex.Lock()
	defer _gModuleMutex.Unlock()

	if _gModules == nil {
		_gModules = make([]base.IModule, 0)
	}
	_gModules = append(_gModules, app)
}

func HasModule(moduleID string) bool {
	_gModuleMutex.Lock()
	defer _gModuleMutex.Unlock()

	for _, v := range _gModules {
		if v.GetModuleID() == moduleID {
			return true
		}
	}
	return false
}

func GetModules() []base.IModule {
	_gModuleMutex.Lock()
	defer _gModuleMutex.Unlock()

	return append(make([]base.IModule, 0), _gModules...)
}
