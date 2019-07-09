package app

import (
	"github.com/liasece/micserver/conf"
	"github.com/liasece/micserver/gate"
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
	serverhandler IServerHandler
	servermanger  ServerManager

	Configer *conf.ServerConfig
	modules  []*module.Module
	gatebase *gate.GateBase

	isStoped bool
}

func (this *App) New(version string) *App {
	log.AutoConfig("/home/jansen/logs/main.log", "Main", true)
	return &App{
		Configer: &conf.ServerConfig{},
		modules:  make([]*module.Module, 0),
	}
}

func (this *App) Init(includegate bool) {
	if includegate {
		this.gatebase = &gate.GateBase{}
		this.gatebase.Init(this.Configer)
	}

	log.Debug("hello world!")
}

func (this *App) GetServerManger() ServerManager {
	return this.servermanger
}

func (this *App) GetServerHandler() IServerHandler {
	return this.serverhandler
}

type IServerHandler interface {
	OnInit()
	OnFinal()
	// OnCreateTCPConnect(serverconn *tcpconn.ServerConn)
	// OnRemoveTCPConnect(serverconn *tcpconn.ServerConn)

	// TCPMsgParse(serverconn *tcpconn.ServerConn,
	// 	msgbin *msg.MessageBinary)
	// GetTCPMsgParseChan(serverconn *tcpconn.ServerConn,
	// 	maxchan int, msg *msg.MessageBinary) int
}

type ServerManager interface {
	OnSignal(os.Signal)
	NotifyOtherMyInfo()
}

// cpu性能测试
func (this *App) startTestCpuProfile() {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := util.GetPanicInfo(recover()); err != nil {
			log.Error("[startTestCpuProfile] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
		}
	}()
	log.Debug("[SubNetManager.startTestCpuProfile] " +
		"[性能分析] StartTestCpuProfile start")
	filename := this.Configer.GetProp("profile_filename")
	testtime := this.Configer.GetPropInt("profile_time")
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	err = pprof.StartCPUProfile(f)
	if err != nil {
		log.Error("[startTestCpuProfile] pprof.StartCPUProfile Err[%s]",
			err.Error())
		return
	}
	for i := 0; i < int(testtime); i++ {
		time.Sleep(1 * time.Second)
	}
	pprof.StopCPUProfile()
	f.Close()
	log.Debug("[SubNetManager.startTestCpuProfile] " +
		"[性能分析] StartTestCpuProfile end")
}

// 监听系统消息
func (this *App) SignalListen(manager ServerManager) {
	c := make(chan os.Signal, 10)
	signal.Notify(c, syscall.SIGINT,
		syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1,
		syscall.SIGUSR2)
	for {
		s := <-c
		log.Debug("[SubNetManager.SignalListen] "+
			"Get signal Signal[%d]", s)
		// manager.OnSignal(s)
		switch s {
		case syscall.SIGUSR1:
			go this.startTestCpuProfile()
		case syscall.SIGUSR2:
		case syscall.SIGTERM:
		case syscall.SIGINT:
			// //收到信号后的处理
			// if this.Configer.Myserverinfo.Servertype !=
			// 	def.TypeSuperServer {
			// 	// 通知我的主服务器，退出连接了
			// 	sendmsg := &comm.SLogoutCommand{}
			// 	this.clientManager.BroadcastByType(def.TypeSuperServer,
			// 		sendmsg)
			// 	// 发送给所有连接到我的服务器，我要关闭了，别再尝试连接我了
			// 	this.taskManager.BroadcastAll(sendmsg)
			// }
			// this.Configer.TerminateServer = true
			this.isStoped = true
		case syscall.SIGQUIT:
			// 捕捉到就退不出了
			buf := make([]byte, 1<<20)
			stacklen := runtime.Stack(buf, true)
			log.Debug("[SubNetManager.SignalListen] "+
				"Received SIGQUIT, \n: Stack[%s]", buf[:stacklen])
		}
	}
}

func (this *App) Run() {
	log.Debug("[App.Run] ----- Main is start ----- ")
	// 初始化参数
	// this.servermanger = manager
	// this.serverhandler = server

	// 绑定本地服务端口
	// 必须等待本地服务器端口绑定完成之后，再进行其他操作
	// 这是服务器加入内部网络的基础
	// err := this.BindMyTCPServer(server)
	// if err != nil {
	// 	log.Error("[StartMain] BindMyTCPServer Err[%s]", err)
	// 	return
	// }
	if this.gatebase != nil {
		tcpport := this.Configer.GetPropUint("tcpouterport")
		tcpport = 8888
		this.gatebase.BindOuterTCP(tcpport)
	}

	// 监听系统Signal
	go this.SignalListen(this.servermanger)

	// 初始化服务器
	// server.OnInit()

	// 保持程序运行
	for !this.isStoped {
		time.Sleep(1 * time.Second)
	}

	// 当程序即将结束时
	// server.OnFinal()
	log.Debug("[App.Run] All server is over add save datas")

	// 等日志打完
	time.Sleep(1 * time.Second)
	log.Debug("[App.Run] ----- Main is quit ----- ")
}
