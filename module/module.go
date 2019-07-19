package module

import (
	"github.com/liasece/micserver/conf"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/server/gate"
	"github.com/liasece/micserver/server/subnet"
)

type IModule interface {
	GetModuleID() string
	InitModule(conf.ModuleConfig)
	TopRunner()
	KillModule()
	InitSubnet(map[string]string)
}

type BaseModule struct {
	*log.Logger
	ModuleID string
	Configer *conf.ModuleConfig

	subnetManager *subnet.SubnetManager
	gateBase      *gate.GateBase
}

func (this *BaseModule) InitModule(configer conf.ModuleConfig) {
	this.Configer = &configer
	// 初始化logger
	if this.Configer.HasModuleSetting("logpath") {
		this.Logger = log.NewLogger(this.Configer.Settings)
		this.SetLogName(this.ModuleID)
	} else {
		this.Logger = log.GetDefaultLogger()
	}
	// 申请内存
	if this.subnetManager == nil {
		this.subnetManager = &subnet.SubnetManager{}
	}
	this.subnetManager.Logger = this.Logger
	// 初始化服务器网络管理器
	this.subnetManager.InitManager(this.Configer)

	this.Debug("module initting...")
	// gateway初始化
	if gateaddr := this.Configer.GetModuleSetting("gatetcpaddr"); gateaddr != "" {
		this.gateBase = &gate.GateBase{}
		this.gateBase.BindOuterTCP(gateaddr)
	}
}

func (this *BaseModule) InitSubnet(subnetAddrMap map[string]string) {
	for k, addr := range subnetAddrMap {
		if k != this.GetModuleID() {
			this.subnetManager.TryConnectServer(k, addr)
		}
	}
}

func (this *BaseModule) GetModuleID() string {
	return this.ModuleID
}

func (this *BaseModule) TopRunner() {
}

func (this *BaseModule) KillModule() {
	this.Debug("KillModule...")
}
