package module

import (
	"github.com/liasece/micserver/conf"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/server/subnet"
)

type Module struct {
	subnetManager *subnet.SubnetManager
	Configer      *conf.ServerConfig
}

func (this *Module) New(version string) *Module {
	log.AutoConfig("/home/jansen/logs/main.log", "Main", true)
	return &Module{
		subnetManager: &subnet.SubnetManager{},
		Configer:      &conf.ServerConfig{},
	}
}

func (this *Module) Init() {
	// 初始化服务器网络管理器
	this.subnetManager.InitManager()
	log.Debug("hello world!")
	// 显示网卡信息
	// localip := util.GetIPv4ByInterface(ifname)
	// log.Debug("网卡名称:%s,ip:%s", ifname, localip)
}

func (this *Module) Run() {
	// this.subnetManager.StartMain()
}
