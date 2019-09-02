package app

import (
	"github.com/liasece/micserver/conf"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/module"
	"github.com/liasece/micserver/util"
	"os"
	"runtime/pprof"
	"time"
)

/**
 * App 是 MicServer 中在 "Module" 上一层的概念，使用 MicServer 的
 * 第一步就是实例化出一个 App 对象，并且向其中插入你的 Modules 。
 * 建议一个代码上下文中仅存在一个 App 对象，如果你的需求让你觉得你有
 * 必要实例化多个 App 在同一个可执行文件中，那么你应该考虑增加一个
 * Module 而不是 App 。
 */
type App struct {
	/*
	 * 匿名成员 Logger 帮助你在此 App 中通过 App.Debug 等方式输出你的
	 * 日志信息到指定的输出上。
	 * 你可以在任何模块代码中新增 Logger 引用，并且从该 Logger 中 Clone()
	 * 出来，定制化 log 的名字/主题 等，Logger 的底层实现已经帮你处理好了
	 * Clone() 出来的 Logger 指向同一个输出。
	 */
	*log.Logger
	Configer *conf.TopConfig
	modules  []module.IModule

	/**
	 * App 是否已停止，如果为 true ，App将会在下一个循环周期执行清理工作，
	 * 退出阻塞循环。
	 */
	isStoped bool
}

func (this *App) Init(configer *conf.TopConfig, modules []module.IModule) {
	this.Configer = configer
	if this.Configer.AppConfig.HasSetting("logpath") {
		this.Configer.AppConfig.AppSettings["logfilename"] = "app.log"
		this.Logger = log.NewLogger(this.Configer.AppConfig.AppSettings)
		log.SetDefaultLogger(this.Logger)
		this.Logger.SetLogName("app")
	} else {
		this.Logger = log.GetDefaultLogger()
	}

	this.modules = modules
	// create all module
	for _, m := range this.modules {
		this.Debug("[App.Init] init module: %s", m.GetModuleID())
		m.InitModule(*this.Configer.AppConfig.GetModuleConfig(m.GetModuleID()))
		m.AfterInitModule()
		go m.TopRunner()
	}

	subnetTCPAddrMap := this.Configer.AppConfig.GetSubnetTCPAddrMap()
	for _, m := range this.modules {
		m.InitSubnet(subnetTCPAddrMap)
	}

	this.Debug("[App.Init] App 初始化成功！")
	this.Debug("[App.Init] App 初始化 Module 数量：%d", len(this.modules))
}

// cpu性能测试
func (this *App) startTestCpuProfile() {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := util.GetPanicInfo(recover()); err != nil {
			this.Error("[startTestCpuProfile] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
		}
	}()
	this.Debug("[SubNetManager.startTestCpuProfile] " +
		"[性能分析] StartTestCpuProfile start")
	filename := this.Configer.GetProp("profile_filename")
	testtime := this.Configer.GetPropInt("profile_time")
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	err = pprof.StartCPUProfile(f)
	if err != nil {
		this.Error("[startTestCpuProfile] pprof.StartCPUProfile Err[%s]",
			err.Error())
		return
	}
	for i := 0; i < int(testtime); i++ {
		time.Sleep(1 * time.Second)
	}
	pprof.StopCPUProfile()
	f.Close()
	this.Debug("[SubNetManager.startTestCpuProfile] " +
		"[性能分析] StartTestCpuProfile end")
}

// 阻塞运行
func (this *App) RunAndBlock() {
	this.Debug("[App.Run] ----- Main has started ----- ")

	// 监听系统Signal
	go this.SignalListen()

	// 保持程序运行
	for !this.isStoped {
		time.Sleep(1 * time.Second)
	}

	for _, v := range this.modules {
		v.KillModule()
	}

	for _, v := range this.modules {
		if !v.IsStopped() {
			// 等待模块退出完成
			time.Sleep(300 * time.Millisecond)
		}
	}

	// 当程序即将结束时
	// server.OnFinal()
	this.Debug("[App.Run] All server is over add save datas")

	this.Debug("[App.Run] ----- Main has stopped ----- ")
	// 等日志打完
	// time.Sleep(1 * time.Second)
	this.Logger.CloseLogger()
}

// 默认阻塞运行
func (this *App) Run() {
	this.RunAndBlock()
}
