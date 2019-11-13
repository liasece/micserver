package module

import (
	"time"

	"github.com/liasece/micserver/base"
	"github.com/liasece/micserver/conf"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/server"
	"github.com/liasece/micserver/util"
	"github.com/liasece/micserver/util/monitor"
	"github.com/liasece/micserver/util/timer"
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
	timer.TimerManager
	server.Server

	// 模块配置
	configer *conf.ModuleConfig
	// 模块的负载
	load monitor.Load

	moduleID        string
	hasKilledModule bool
	hasStopped      bool
	lastCheckLoad   int64
}

func (this *BaseModule) InitModule(configer conf.ModuleConfig) {
	this.configer = &configer
	// 初始化logger
	if this.configer.ExistInModule(conf.LogWholePath) {
		this.Logger = log.NewLogger(this.configer.GetBool(conf.IsDaemon),
			this.configer.GetString(conf.LogWholePath))
		this.SetLogName(this.moduleID)
	} else {
		this.Logger = log.GetDefaultLogger().Clone()
		this.Logger.SetLogName(this.moduleID)
	}
	this.Syslog("[BaseModule.InitModule] module initting...")
	this.Server.SetLogger(this.Logger)
	this.Server.Init(this.moduleID)
	this.Server.InitSubnet(this.configer)

	// gateway初始化
	if gateaddr := this.configer.GetString(conf.GateTCPAddr); gateaddr != "" {
		this.Server.InitGate(gateaddr)
	}

	this.RegTimer(time.Second*5, 0, false, this.watchLoadToLog)
}

func (this *BaseModule) AfterInitModule() {
	this.Syslog("[BaseModule.AfterInitModule] 模块 [%s] 初始化完成",
		this.GetModuleID())
}

func (this *BaseModule) GetConfiger() *conf.ModuleConfig {
	return this.configer
}

func (this *BaseModule) GetModuleID() string {
	return this.moduleID
}

func (this *BaseModule) SetModuleID(id string) {
	this.moduleID = id
}

func (this *BaseModule) KillModule() {
	this.Syslog("[BaseModule] Killing module...")
	this.Server.Stop()
	this.hasKilledModule = true
	this.KillRegister()

	// 退出完成
	this.hasStopped = true
}

func (this *BaseModule) IsStopped() bool {
	return this.hasStopped
}

func (this *BaseModule) TopRunner() {
	this.RegTimer(time.Minute, 0, false, func(t time.Duration) bool {
		this.Syslog("[BaseModule] Timer 1 Minute...")
		return true
	})
}

func (this *BaseModule) GetModuleType() string {
	return util.GetModuleIDType(this.moduleID)
}

func (this *BaseModule) watchLoadToLog(dt time.Duration) bool {
	load := this.load.GetLoad()
	incValue := load - this.lastCheckLoad
	if incValue > 0 {
		this.Info("[BaseModule] Within %d sec load:[%d]",
			int64(dt.Seconds()), incValue)
	}
	this.lastCheckLoad = load
	return true
}
