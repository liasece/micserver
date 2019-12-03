package micserver

import (
	"github.com/liasece/micserver/app"
	"github.com/liasece/micserver/conf"
	"math/rand"
	"time"
)

func SetupApp(configpath string) (*app.App, error) {
	// 初始化随机数种子
	rand.Seed(time.Now().UnixNano())
	configer, err := conf.LoadConfig(configpath)
	if err != nil {
		return nil, err
	}
	res := &app.App{}
	res.Setup(configer)
	return res, nil
}
