package micserver

import (
	"math/rand"
	"time"

	"github.com/liasece/micserver/app"
	"github.com/liasece/micserver/conf"
)

// SetupApp func
func SetupApp(configpath string) (*app.App, error) {
	// 初始化随机数种子
	rand.Seed(time.Now().UnixNano())
	cfg, err := conf.LoadConfig(configpath)
	if err != nil {
		return nil, err
	}
	res := &app.App{}
	res.Setup(cfg)
	return res, nil
}
