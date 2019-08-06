package app

import (
	"github.com/liasece/micserver/conf"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/module"
	"github.com/liasece/micserver/util"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"
	"time"
)

type App struct {
	*log.Logger
	Configer *conf.TopConfig
	modules  []module.IModule

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

// 监听系统消息
func (this *App) SignalListen() {
	c := make(chan os.Signal, 10)
	signal.Notify(c, syscall.SIGINT,
		syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1,
		syscall.SIGUSR2)
	for {
		s := <-c
		this.Debug("[SubNetManager.SignalListen] "+
			"Get signal Signal[%d]", s)
		// manager.OnSignal(s)
		switch s {
		case syscall.SIGUSR1:
			go this.startTestCpuProfile()
		case syscall.SIGUSR2:
		case syscall.SIGTERM:
		case syscall.SIGINT:
			// //收到信号后的处理
			this.isStoped = true
		case syscall.SIGQUIT:
			// 捕捉到就退不出了
			buf := make([]byte, 1<<20)
			stacklen := runtime.Stack(buf, true)
			this.Debug("[SubNetManager.SignalListen] "+
				"Received SIGQUIT, \n: Stack[%s]", buf[:stacklen])
		}
	}
}

func (this *App) Run() {
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
