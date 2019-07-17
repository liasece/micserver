package micserver

import (
	"github.com/liasece/micserver/app"
	"github.com/liasece/micserver/conf"
	"github.com/liasece/micserver/module"
	"math/rand"
	"time"
)

func CreateApp(configpath string, modules []module.IModule) (*app.App, error) {
	// 初始化随机数种子
	rand.Seed(time.Now().UnixNano())
	configer, err := conf.LoadConfig(configpath)
	if err != nil {
		return nil, err
	}
	res := &app.App{}
	res.Init(configer, modules)
	return res, nil
}
