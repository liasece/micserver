package micserver

import (
	"github.com/liasece/micserver/app"
	"math/rand"
	"time"
)

var Version string = "0.0.1"

var defaultApp app.App

func CreateApp() *app.App {
	// 初始化随机数种子
	rand.Seed(time.Now().UnixNano())

	return defaultApp.New(Version)
}
