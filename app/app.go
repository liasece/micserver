// Package app micserver 最基础的运行单位，app中包含了多个module，app在启动时会初始化所有module，
// 并且根据配置初始化module之间的连接。
// App 是 micserver 中在 "Module" 上一层的概念，使用 micserver 的
// 第一步就是实例化出一个 App 对象，并且向其中插入你的 Modules 。
// 建议一个代码上下文中仅存在一个 App 对象，如果你的需求让你觉得你有
// 必要实例化多个 App 在同一个可执行文件中，那么你应该考虑增加一个
// Module 而不是 App 。
package app

import (
	"os"
	"runtime/pprof"
	"sync"
	"time"

	"github.com/liasece/micserver/conf"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/module"
	"github.com/liasece/micserver/process"
	"github.com/liasece/micserver/util/sysutil"
)

// App basic runtime
type App struct {

	// 匿名成员 Logger 帮助你在此 App 中通过 App.Debug 等方式输出你的
	// 日志信息到指定的输出上。
	// 你可以在任何模块代码中新增 Logger 引用，并且从该 Logger 中 Clone()
	// 出来，定制化 log 的名字/主题 等，Logger 的底层实现已经帮你处理好了
	// Clone() 出来的 Logger 指向同一个输出。
	*log.Logger
	Config  *conf.TopConfig
	modules []module.IModule

	initOnce sync.Once

	// App 是否已停止，如果为 true ，App将会在下一个循环周期执行清理工作，
	// 退出阻塞循环。
	isStoped chan struct{}
}

// Setup 初始化App的设置
func (a *App) Setup(cfg *conf.TopConfig) {
	{
		// check parameters
		if cfg == nil {
			cfg = conf.Default()
		}
	}

	process.AddApp(a)
	a.isStoped = make(chan struct{})
	a.Config = cfg
	// 初始化Log
	if a.Config.AppConfig.Exist(conf.LogWholePath) {
		setting := a.Config.AppConfig
		a.Logger = log.NewLogger(nil, log.Options().NoConsole(setting.GetBool(conf.IsDaemon)).FilePaths(setting.GetString(conf.LogWholePath)))
		log.SetDefaultLogger(a.Logger)
		a.Logger.SetLogName("app")
	} else {
		a.Logger = log.GetDefaultLogger()
	}
	// 初始化Log等级
	if a.Config.AppConfig.Exist(conf.LogLevel) {
		err := a.Logger.SetLogLevelByStr(
			a.Config.AppConfig.GetString(conf.LogLevel))
		if err != nil {
			a.Error("Set log level error", log.ErrorField(err))
		}
	}
	a.Info("APP setup secess!!!")
}

// Init this app
func (a *App) Init(modules []module.IModule) (err error) {
	return a.tryInit(modules)
}

// tryInit if this app has not init
func (a *App) tryInit(modules []module.IModule) (err error) {
	a.initOnce.Do(func() {
		err = a.init(modules)
	})
	return
}

// init func
func (a *App) init(modules []module.IModule) error {
	a.modules = modules
	// create all module
	for _, m := range a.modules {
		process.AddModule(m)
		a.Syslog("[App.Init] Initing module", log.String("ModuleID", m.GetModuleID()), log.String("ModuleType", m.GetModuleType()), log.Int("ModuleNum", m.GetModuleNum()), log.Uint32("ModuleIDHash", m.GetModuleIDHash()))
		{
			err := m.InitModule(*a.Config.AppConfig.GetModuleConfig(m.GetModuleID()))
			if err != nil {
				return err
			}
		}
		m.AfterInitModule()
		a.Syslog("[App.Init] Init module finish", log.String("ModuleID", m.GetModuleID()), log.String("ModuleType", m.GetModuleType()), log.Int("ModuleNum", m.GetModuleNum()), log.Uint32("ModuleIDHash", m.GetModuleIDHash()))
		go m.TopRunner()
	}

	subnetTCPAddrMap := a.Config.AppConfig.GetSubnetTCPAddrMap()
	for _, m := range a.modules {
		m.BindSubnet(subnetTCPAddrMap)
	}

	a.Syslog("[App.Init] App init finish!", log.Int("ModuleNum", len(a.modules)))
	return nil
}

// startTestCPUProfile cpu性能测试
func (a *App) startTestCPUProfile() {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if stackInfo, err := sysutil.GetPanicInfo(recover()); err != nil {
			a.Error("[startTestCPUProfile] Panic", log.String("stackInfo", stackInfo), log.ErrorField(err))
		}
	}()
	a.Debug("[SubNetManager.startTestCPUProfile] Cpu profile start")
	filename := a.Config.GetProp("profile_filename")
	testtime := a.Config.GetPropInt64("profile_time")
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	err = pprof.StartCPUProfile(f)
	if err != nil {
		a.Error("[startTestCPUProfile] pprof.StartCPUProfile error", log.ErrorField(err))
		return
	}
	for i := 0; i < int(testtime); i++ {
		time.Sleep(1 * time.Second)
	}
	pprof.StopCPUProfile()
	f.Close()
	a.Debug("[SubNetManager.startTestCPUProfile] Cpu profile finish")
}

// RunAndBlock 运行并阻塞本App，直到程序主动退出
func (a *App) RunAndBlock(modules []module.IModule) {
	a.tryInit(modules)
	a.Syslog("[App.Run] ----- App has started ----- ")

	// 监听系统Signal
	go a.SignalListen()

	// 保持程序运行
	<-a.isStoped

	for _, v := range a.modules {
		v.KillModule()
	}

	for _, v := range a.modules {
		if !v.IsStopped() {
			// 等待模块退出完成
			time.Sleep(300 * time.Millisecond)
		}
	}

	// 当程序即将结束时
	a.Syslog("[App.RunAndBlock] All modules is stopped")

	a.Debug("[App.RunAndBlock] ----- App has stopped ----- ")
	// 等日志打完
	time.Sleep(500 * time.Millisecond)
}

// Run 运行，默认执行 RunAndBlock ，阻塞
func (a *App) Run(modules []module.IModule) {
	a.RunAndBlock(modules)
}

// Stop the app
func (a *App) Stop() {
	a.isStoped <- struct{}{}
}
