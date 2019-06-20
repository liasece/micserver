package micserver

import (
	"github.com/liasece/micserver/module"
	"math/rand"
)

var Version string = "0.0.1"

var defaultApp module.App

func CreateApp() *module.App {
	// 初始化随机数种子
	rand.Seed(time.Now().UnixNano())

	return defaultApp.NewApp(Version)
}
