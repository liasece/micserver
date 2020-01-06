/*
基础模块，实现了基本的模块接口，利于上层使用 micserver ，在业务中自定义模块时，
可以直接继承该基础模块，继承于基础模块（使用匿名成员实现）的上层业务模块，
即实现了 IModule 接口。例如：
	type FooModule struct {
		module.BaseModule
		...
	}
需要注意的是有些方法在重载时，需要在重载中调用父类的该方法，且调用顺序有要求：
	如果你需要在例如 AfterInitModule() 中增加逻辑，请使用如下顺序：
		func (this *FooModule) AfterInitModule() {
			// 先调用父类方法
			this.BaseModule.AfterInitModule()
			// 其他逻辑
			...
		}
	如果你需要在例如 KillModule() 中增加逻辑，请使用如下顺序：
		func (this *FooModule) KillModule() {
			// 其他逻辑
			...
			// 最后调用父类方法
			this.BaseModule.KillModule()
		}
*/
package module

import (
	"time"

	"github.com/liasece/micserver/base"
	"github.com/liasece/micserver/conf"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/roc"
	"github.com/liasece/micserver/server"
	"github.com/liasece/micserver/util"
	"github.com/liasece/micserver/util/hash"
	"github.com/liasece/micserver/util/monitor"
	"github.com/liasece/micserver/util/timer"
	"github.com/liasece/micserver/util/uid"
)

// 一个模块应具备的接口
type IModule interface {
	base.IModule
	InitModule(conf.ModuleConfig)
	BindSubnet(map[string]string)
	AfterInitModule()
	TopRunner()
	KillModule()
	IsStopped() bool
	GetConfiger() *conf.ModuleConfig
	ROCCallNR(callpath *roc.ROCPath, callarg []byte) error
	ROCCallBlock(callpath *roc.ROCPath, callarg []byte) ([]byte, error)
}

// 基础模块
type BaseModule struct {
	*log.Logger
	timer.TimerManager
	server.Server

	// 模块配置
	configer *conf.ModuleConfig
	// 模块的负载
	load monitor.Load

	moduleID     string
	moduleIDHash uint32
	moduleType   string
	moduleNum    int

	hasKilledModule bool
	hasStopped      bool
	lastCheckLoad   int64
}

// 初始化模块
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

// 在初始化完成后调用
func (this *BaseModule) AfterInitModule() {
	this.Syslog("[BaseModule.AfterInitModule] 模块 [%s] 初始化完成",
		this.GetModuleID())
}

// 获取模块的配置
func (this *BaseModule) GetConfiger() *conf.ModuleConfig {
	return this.configer
}

// 获取模块的ID，模块的ID有模块类型和模块编号确定，如
// 	moduleid = fmt.Sprintf("%s%d", typ, num)
func (this *BaseModule) GetModuleID() string {
	return this.moduleID
}

// 设置模块的ID，需要谨慎使用，不可在模块运行起来后设置！
func (this *BaseModule) SetModuleID(id string) {
	this.moduleID = id
	this.moduleType = util.GetModuleIDType(id)
	this.moduleNum = util.GetModuleIDNum(id)
	this.moduleIDHash = hash.GetHash([]byte(id))
}

// 获取模块类型
func (this *BaseModule) GetModuleType() string {
	return this.moduleType
}

// 获取模块编号
func (this *BaseModule) GetModuleNum() int {
	return this.moduleNum
}

// 获取模块ID哈希值
func (this *BaseModule) GetModuleIDHash() uint32 {
	return this.moduleIDHash
}

// 在该模块环境下生成一个UUID，这个UUID保证在本模块中是唯一的
func (this *BaseModule) GenUniqueID() (string, error) {
	return uid.GenUniqueID(uint16(this.GetModuleIDHash()))
}

// 当模块被中止时调用
func (this *BaseModule) KillModule() {
	this.Syslog("[BaseModule] Killing module...")
	this.Server.Stop()
	this.hasKilledModule = true
	this.KillRegister()

	// 退出完成
	this.hasStopped = true
}

// 判断模块是否已中止
func (this *BaseModule) IsStopped() bool {
	return this.hasStopped
}

// 开始运行一个模块
func (this *BaseModule) TopRunner() {
	this.RegTimer(time.Minute, 0, false, func(t time.Duration) bool {
		this.Syslog("[BaseModule] Timer 1 Minute...")
		return true
	})
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
