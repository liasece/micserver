package module

import (
	"time"

	"github.com/liasece/micserver/base"
	"github.com/liasece/micserver/conf"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/server"
	"github.com/liasece/micserver/util"
)

type IModule interface {
	base.IModule
	InitModule(conf.ModuleConfig)
	BindSubnet(map[string]string)
	AfterInitModule()
	TopRunner()
	KillModule()
	IsStopped() bool
	GetConfiger() *conf.ModuleConfig
}

type BaseModule struct {
	*log.Logger
	util.TimerManager
	server.Server

	moduleID string
	Configer *conf.ModuleConfig

	hasKilledModule bool
	hasStopped      bool
	// 模块的负载
	Load          util.Load
	lastCheckLoad int64
}

func (this *BaseModule) InitModule(configer conf.ModuleConfig) {
	this.Configer = &configer
	// 初始化logger
	if this.Configer.HasSetting("logpath") {
		this.Logger = log.NewLogger(this.Configer.GetModuleSettingMap())
		this.SetLogName(this.moduleID)
	} else {
		this.Logger = log.GetDefaultLogger().Clone()
		this.Logger.SetLogName(this.moduleID)
	}
	this.Debug("[BaseModule.InitModule] module initting...")
	this.Server.SetLogger(this.Logger)
	this.Server.Init(this.moduleID)
	this.Server.InitSubnet(this.Configer)

	// gateway初始化
	if gateaddr := this.Configer.GetModuleSetting("gatetcpaddr"); gateaddr != "" {
		this.Server.InitGate(gateaddr)
	}

	this.RegTimer(time.Second*5, 0, false, this.watchLoadToLog)
}

func (this *BaseModule) AfterInitModule() {
	this.Debug("[BaseModule.AfterInitModule] 模块 [%s] 初始化完成",
		this.GetModuleID())
}

func (this *BaseModule) GetConfiger() *conf.ModuleConfig {
	return this.Configer
}

func (this *BaseModule) GetModuleID() string {
	return this.moduleID
}

func (this *BaseModule) SetModuleID(id string) {
	this.moduleID = id
}

func (this *BaseModule) KillModule() {
	this.Debug("Killing module...")
	this.Server.Stop()
	this.hasKilledModule = true
	this.KillRegister()

	// 退出完成
	this.hasStopped = true
	this.Logger.CloseLogger()
}

func (this *BaseModule) IsStopped() bool {
	return this.hasStopped
}

func (this *BaseModule) TopRunner() {
	this.RegTimer(time.Minute, 0, false, func(t time.Duration) bool {
		this.Debug("Timer 1 Minute...")
		return true
	})
}

func (this *BaseModule) GetServerType() string {
	return util.GetServerIDType(this.moduleID)
}

func (this *BaseModule) watchLoadToLog(dt time.Duration) bool {
	load := this.Load.GetLoad()
	incValue := load - this.lastCheckLoad
	if incValue > 0 {
		this.Info("[BaseModule]  Within %d sec load:[%d]",
			int64(dt.Seconds()), incValue)
	}
	this.lastCheckLoad = load
	return true
}
