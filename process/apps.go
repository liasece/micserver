// Package process 是一个 micserver 中的单例，用于获取该进程中存在的所有 App 以及 Module，
// 以便实现例如 roc catch 的优化，有一些数据不必每个 Module 都持有一份。
package process

import (
	"sync"

	"github.com/liasece/micserver/base"
)

var (
	_gApps     []base.IApp
	_gAppMutex sync.Mutex
)

// AddApp 增加一个App
func AddApp(app base.IApp) {
	_gAppMutex.Lock()
	defer _gAppMutex.Unlock()

	if _gApps == nil {
		_gApps = make([]base.IApp, 0)
	}
	_gApps = append(_gApps, app)
}

// GetApps 获取当前进程的 App 列表
func GetApps() []base.IApp {
	_gAppMutex.Lock()
	defer _gAppMutex.Unlock()

	return append(make([]base.IApp, 0), _gApps...)
}
