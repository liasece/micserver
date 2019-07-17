package module

import (
	"github.com/liasece/micserver/conf"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/server/subnet"
)

type IModule interface {
	GetModuleID() string
	InitModule(conf.ModuleConfig)
	TopRunner()
	KillModule()
}

type BaseModule struct {
	ModuleID string
	Configer *conf.ModuleConfig
	Logger   *log.Logger

	subnetManager *subnet.SubnetManager
}

func (this *BaseModule) InitModule(configer conf.ModuleConfig) {
	// 申请内存
	if this.subnetManager == nil {
		this.subnetManager = &subnet.SubnetManager{}
	}
	this.Configer = &configer
	// 初始化logger
	if this.Configer.HasSetting("logpath") {
		this.Logger = log.NewLogger(this.Configer.Settings)
	} else {
		this.Logger = log.GetDefaultLogger()
	}

	// 初始化服务器网络管理器
	this.subnetManager.InitManager()

	this.Logger.Debug("module initting...")
	// 显示网卡信息
}

func (this *BaseModule) GetModuleID() string {
	return this.ModuleID
}

func (this *BaseModule) TopRunner() {
	// this.subnetManager.StartMain()
}

func (this *BaseModule) KillModule() {
	this.Logger.Debug("KillModule...")
}
