package process

import (
	"sync"

	"github.com/liasece/micserver/base"
)

var (
	_gApps     []base.IApp
	_gAppMutex sync.Mutex
)

func AddApp(app base.IApp) {
	_gAppMutex.Lock()
	defer _gAppMutex.Unlock()

	if _gApps == nil {
		_gApps = make([]base.IApp, 0)
	}
	_gApps = append(_gApps, app)
}

func GetApps() []base.IApp {
	_gAppMutex.Lock()
	defer _gAppMutex.Unlock()

	return append(make([]base.IApp, 0), _gApps...)
}
