package module

import (
	"github.com/liasece/micserver/conf"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/server/subnet"
	// "github.com/liasece/micserver/util"
)

type App struct {
	subnetManager *subnet.SubnetManger
	Configer      *conf.ServerConfig
}

func (this *App) InitApp() {
	// 初始化服务器网络管理器
	this.subnetManager.InitManager()
}

func (this *App) NewApp(version string) *App {
	log.AutoConfig("/home/jansen/logs/main.log", "Main", true)
	return &App{
		subnetManager: &subnet.SubnetManger{},
		Configer:      &conf.ServerConfig{},
	}
}

func (this *App) Init() {
	this.InitApp()
	log.Debug("hello world!")
	// 显示网卡信息
	// localip := util.GetIPv4ByInterface(ifname)
	// log.Debug("网卡名称:%s,ip:%s", ifname, localip)
}
