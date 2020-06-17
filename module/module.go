// Package module 基础模块，实现了基本的模块接口，利于上层使用 micserver ，在业务中自定义模块时，
// 可以直接继承该基础模块，继承于基础模块（使用匿名成员实现）的上层业务模块，
// 即实现了 IModule 接口。例如：
// 	type FooModule struct {
// 		module.BaseModule
// 		...
// 	}
// 需要注意的是有些方法在重载时，需要在重载中调用父类的该方法，且调用顺序有要求：
// 	如果你需要在例如 AfterInitModule() 中增加逻辑，请使用如下顺序：
// 		func (bm *FooModule) AfterInitModule() {
// 			// 先调用父类方法
// 			bm.BaseModule.AfterInitModule()
// 			// 其他逻辑
// 			...
// 		}
// 	如果你需要在例如 KillModule() 中增加逻辑，请使用如下顺序：
// 		func (bm *FooModule) KillModule() {
// 			// 其他逻辑
// 			...
// 			// 最后调用父类方法
// 			bm.BaseModule.KillModule()
// 		}
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

// IModule 一个模块应具备的接口
type IModule interface {
	base.IModule
	InitModule(conf.ModuleConfig)
	BindSubnet(map[string]string)
	AfterInitModule()
	TopRunner()
	KillModule()
	IsStopped() bool
	GetConfiger() *conf.ModuleConfig
	ROCCallNR(callpath *roc.Path, callarg []byte) error
	ROCCallBlock(callpath *roc.Path, callarg []byte) ([]byte, error)
}

// BaseModule 基础模块
type BaseModule struct {
	*log.Logger
	TimerManager timer.Manager
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

// InitModule 初始化模块
func (bm *BaseModule) InitModule(configer conf.ModuleConfig) {
	bm.configer = &configer
	// 初始化logger
	if bm.configer.ExistInModule(conf.LogWholePath) {
		bm.Logger = log.NewLogger(&log.Options{
			NoConsole: bm.configer.GetBool(conf.IsDaemon),
			FilePaths: []string{bm.configer.GetString(conf.LogWholePath)},
		})
		bm.SetLogName(bm.moduleID)
	} else {
		bm.Logger = log.GetDefaultLogger().Clone()
		bm.Logger.SetLogName(bm.moduleID)
	}
	bm.Syslog("[BaseModule.InitModule] module initting...")
	bm.Server.SetLogger(bm.Logger)
	bm.Server.Init(bm.moduleID)
	bm.Server.InitSubnet(bm.configer)

	// gateway初始化
	if gateaddr := bm.configer.GetString(conf.GateTCPAddr); gateaddr != "" {
		bm.Server.InitGate(gateaddr)
	}

	bm.TimerManager.RegTimer(time.Second*5, 0, false, bm.watchLoadToLog)
}

// AfterInitModule 在初始化完成后调用
func (bm *BaseModule) AfterInitModule() {
	bm.Syslog("[BaseModule.AfterInitModule] 模块 [%s] 初始化完成",
		bm.GetModuleID())
}

// GetConfiger 获取模块的配置
func (bm *BaseModule) GetConfiger() *conf.ModuleConfig {
	return bm.configer
}

// GetModuleID 获取模块的ID，模块的ID有模块类型和模块编号确定，如
// moduleid = fmt.Sprintf("%s%d", typ, num)
func (bm *BaseModule) GetModuleID() string {
	return bm.moduleID
}

// SetModuleID 设置模块的ID，需要谨慎使用，不可在模块运行起来后设置！
func (bm *BaseModule) SetModuleID(id string) {
	bm.moduleID = id
	bm.moduleType = util.GetModuleIDType(id)
	bm.moduleNum = util.GetModuleIDNum(id)
	bm.moduleIDHash = hash.GetHash([]byte(id))
}

// GetModuleType 获取模块类型
func (bm *BaseModule) GetModuleType() string {
	return bm.moduleType
}

// GetModuleNum 获取模块编号
func (bm *BaseModule) GetModuleNum() int {
	return bm.moduleNum
}

// GetModuleIDHash 获取模块ID哈希值
func (bm *BaseModule) GetModuleIDHash() uint32 {
	return bm.moduleIDHash
}

// GenUniqueID 在该模块环境下生成一个UUID，这个UUID保证在本模块中是唯一的
func (bm *BaseModule) GenUniqueID() (string, error) {
	return uid.GenUniqueID(uint16(bm.GetModuleIDHash()))
}

// KillModule 当模块被中止时调用
func (bm *BaseModule) KillModule() {
	bm.Syslog("[BaseModule] Killing module...")
	bm.Server.Stop()
	bm.hasKilledModule = true
	bm.TimerManager.KillRegister()

	// 退出完成
	bm.hasStopped = true
}

// IsStopped 判断模块是否已中止
func (bm *BaseModule) IsStopped() bool {
	return bm.hasStopped
}

// TopRunner 开始运行一个模块
func (bm *BaseModule) TopRunner() {
	bm.TimerManager.RegTimer(time.Minute, 0, false, func(t time.Duration) bool {
		bm.Syslog("[BaseModule] Timer 1 Minute...")
		return true
	})
}

func (bm *BaseModule) watchLoadToLog(dt time.Duration) bool {
	load := bm.load.GetLoad()
	incValue := load - bm.lastCheckLoad
	if incValue > 0 {
		bm.Info("[BaseModule] Within %d sec load:[%d]",
			int64(dt.Seconds()), incValue)
	}
	bm.lastCheckLoad = load
	return true
}
