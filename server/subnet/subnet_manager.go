package subnet

import (
	"fmt"
	"github.com/liasece/micserver/comm"
	"github.com/liasece/micserver/conf"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/server/subnet/serconfs"
	"github.com/liasece/micserver/tcpconn"
	"github.com/liasece/micserver/util"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	"syscall"
	"time"
)

type ConnectMsgQueueStruct struct {
	task *tcpconn.ServerConn
	msg  *msg.MessageBinary
}

func CheckServerType(servertype uint32) bool {
	if servertype <= 0 || servertype > 10 {
		return false
	}
	return true
}

// websocket连接管理器
type SubnetManger struct {
	serverhandler IServerHandler
	servermanger  ServerManager

	clientManager *GBTCPClientManager
	taskManager   *GBTCPTaskManager

	ServerConfigs serconfs.ServerConfigsManager // 所有服务器信息

	moudleConf *conf.ServerConfig
}

func (this *SubnetManger) InitManager() {

	this.clientManager = &GBTCPClientManager{}
	this.clientManager.subnetManager = this
	this.clientManager.connPool.
		Init(tcpconn.ServerSCTypeClient, 1)
	this.clientManager.superexitchan = make(chan bool, 1)

	this.taskManager = &GBTCPTaskManager{}
	this.taskManager.subnetManager = this
	this.taskManager.connPool.
		Init(tcpconn.ServerSCTypeTask, 2)
}

func (this *SubnetManger) GetLatestVersionServerConfigByType(servertype uint32) uint64 {
	latestVersion := uint64(0)
	this.ServerConfigs.RangeServerConfig(
		func(value comm.SServerInfo) bool {
			if value.Servertype == servertype &&
				value.Version > latestVersion {
				latestVersion = value.Version
			}
			return true
		})
	return latestVersion
}

func (this *SubnetManger) GetServerManger() ServerManager {
	return this.servermanger
}

func (this *SubnetManger) GetServerHandler() IServerHandler {
	return this.serverhandler
}

type IServerHandler interface {
	OnInit()
	OnFinal()
	OnCreateTCPConnect(serverconn *tcpconn.ServerConn)
	OnRemoveTCPConnect(serverconn *tcpconn.ServerConn)

	TCPMsgParse(serverconn *tcpconn.ServerConn,
		msgbin *msg.MessageBinary)
	GetTCPMsgParseChan(serverconn *tcpconn.ServerConn,
		maxchan int, msg *msg.MessageBinary) int
}

type ServerManager interface {
	OnSignal(os.Signal)
	NotifyOtherMyInfo()
}

func (this *SubnetManger) StartMain(server IServerHandler, manager ServerManager) {
	log.Debug("[SubNetManager.StartMain] " +
		"Main is start------")
	// 初始化参数
	this.servermanger = manager
	this.serverhandler = server

	// 绑定本地服务端口
	// 必须等待本地服务器端口绑定完成之后，再进行其他操作
	// 这是服务器加入内部网络的基础
	err := this.BindMyTCPServer(server)
	if err != nil {
		log.Error("[StartMain] BindMyTCPServer Err[%s]", err)
		return
	}

	// 监听系统Signal
	go this.SignalListen(manager)

	// 初始化服务器
	server.OnInit()

	// 保持程序运行
	for !this.moudleConf.TerminateServer {
		time.Sleep(1 * time.Second)
	}

	// 当程序即将结束时
	server.OnFinal()
	log.Debug("[SubNetManager.StartMain] " +
		"All server is over add save datas")

	// 等日志打完
	time.Sleep(1 * time.Second)
}

var serverLoginMutex sync.Mutex

func (this *SubnetManger) OnServerLogin(tcptask *tcpconn.ServerConn,
	tarinfo *comm.SLoginCommand) {
	serverLoginMutex.Lock()
	defer serverLoginMutex.Unlock()

	// 来源服务器请求登陆本服务器
	oldtcptask := this.taskManager.GetTCPTask(uint64(tarinfo.Serverid))
	if oldtcptask != nil {
		// 已经连接成功过了，非法连接
		log.Error("[SubNetManager.OnServerLogin] "+
			"收到了重复的Server连接请求 Msg[%s]",
			tarinfo.GetJson())
		retmsg := &comm.SLoginRetCommand{}
		retmsg.Loginfailed = 1
		tcptask.SendCmd(retmsg)
		tcptask.Terminate()
		return
	}
	if !CheckServerType(tarinfo.Servertype) {
		// 未知的服务器类型
		log.Error("[SubNetManager.OnServerLogin] "+
			"收到未知的服务器类型 Msg[%s]",
			tarinfo.GetJson())
		retmsg := &comm.SLoginRetCommand{}
		retmsg.Loginfailed = 1
		tcptask.SendCmd(retmsg)
		tcptask.Terminate()
		return
	}
	var serverconfig comm.SServerInfo

	// 来源服务器地址
	remoteaddr := tcptask.GetConn().RemoteAddr().String()
	// 检查来源服务器IP地址与消息中的地址一致
	sameip := false
	s := strings.Split(remoteaddr, ":")
	if len(s) == 2 && s[0] == tarinfo.Serverip {
		// IP 地址正确
		sameip = true
	}
	// 获取来源服务器ID在本地的配置
	serverconfig = this.ServerConfigs.
		GetServerConfigByInfo(tarinfo)
	// 如果成功获取到了一个serverid
	if serverconfig.Serverid != 0 {
		// 检查该配置是否已被占用
		oldtcptask = this.taskManager.GetTCPTask(uint64(tarinfo.Serverid))
		if oldtcptask != nil {
			// 已经连接成功过了，非法连接
			log.Error("[SubNetManager.OnServerLogin] "+
				"收到了重复的Server连接请求 2 Msg[%s]",
				tarinfo.GetJson())
			retmsg := &comm.SLoginRetCommand{}
			retmsg.Loginfailed = 1
			tcptask.SendCmd(retmsg)
			tcptask.Terminate()
			return
		}
	}
	// 检查是否获取信息成功
	if serverconfig.Serverid == 0 || !sameip ||
		serverconfig.Servertype != tarinfo.Servertype {
		// 如果获取信息不成功
		log.Error("[SubNetManager.OnServerLogin] "+
			"连接分配异常 未知服务器连接 "+
			"Addr[%s] Msg[%s]",
			remoteaddr, tarinfo.GetJson())
		retmsg := &comm.SLoginRetCommand{}
		retmsg.Loginfailed = 1
		tcptask.SendCmd(retmsg)
		tcptask.Terminate()

		// if this.moudleConf.Myserverinfo.Servertype !=
		// 	def.TypeSuperServer {
		// 	requestServerInfo := &comm.SRequestServerInfo{}
		// 	this.clientManager.BroadcastByType(def.TypeSuperServer,
		// 		requestServerInfo)
		// }
		return
	}

	serverconfig.Version = tarinfo.Version

	// 来源服务器检查完毕
	// 完善来源服务器在本服务器的信息
	tcptask.SetVertify(true)
	tcptask.SetTerminateTime(0) // 清除终止时间状态
	tcptask.Serverinfo = serverconfig
	this.taskManager.ChangeTCPTaskTempid(tcptask,
		uint64(tcptask.Serverinfo.Serverid))
	log.Debug("[SubNetManager.OnServerLogin] "+
		"客户端连接验证成功 "+
		" SerID[%d] IP[%s] Type[%d] Port[%d] Name[%s]",
		serverconfig.Serverid, serverconfig.Serverip,
		serverconfig.Servertype, serverconfig.Serverport,
		serverconfig.Servername)
	// 向来源服务器回复登陆成功消息
	retmsg := &comm.SLoginRetCommand{}
	retmsg.Loginfailed = 0
	retmsg.Clientinfo = serverconfig
	retmsg.Taskinfo = this.moudleConf.Myserverinfo
	retmsg.Redisinfo = this.moudleConf.RedisConfig
	tcptask.SendCmd(retmsg)

	// 通知其他服务器，次服务器登陆完成
	// 如果我是SuperServer
	// 向来源服务器发送本地已有的所有服务器信息
	this.taskManager.NotifyAllServerInfo(tcptask)
	// 把来源服务器信息广播给其它所有服务器
	notifymsg := &comm.SStartMyNotifyCommand{}
	notifymsg.Serverinfo = serverconfig
	this.taskManager.BroadcastAll(notifymsg)
}

// 绑定我的HTTP服务 会阻塞
func (this *SubnetManger) BindMyHttpServer() {
	usehttps := this.moudleConf.GetProp("usehttps")
	if this.moudleConf.Myserverinfo.Httpport > 0 {
		httpportstr := fmt.Sprintf(":%d",
			this.moudleConf.Myserverinfo.Httpport)
		if usehttps == "true" {
			log.Debug("[SubNetManager.BindMyHttpServer] "+
				"服务器https绑定成功 HTTPPort[%s] use https",
				httpportstr)
			httpscertfile := this.moudleConf.
				GetProp("httpscertfile")
			httpskeyfile := this.moudleConf.
				GetProp("httpskeyfile")
			err := http.ListenAndServeTLS(httpportstr, httpscertfile,
				httpskeyfile, nil)
			if err != nil {
				panic("ListenAndServe: " + err.Error())
				// return err
			}
		} else {
			log.Debug("[SubNetManager.BindMyHttpServer] "+
				"服务器http绑定成功,%s", httpportstr)
			err := http.ListenAndServe(httpportstr, nil)
			if err != nil {
				panic("ListenAndServe: " + err.Error())
				// return err
			}
		}
	}
}

// 绑定我的HTTPS服务 会阻塞
func (this *SubnetManger) BindMyHttpsServer() error {
	if this.moudleConf.Myserverinfo.Httpsport > 0 {
		httpportstr := fmt.Sprintf(":%d",
			this.moudleConf.Myserverinfo.Httpsport)
		log.Debug("[SubNetManager.BindMyHttpsServer] "+
			"服务器https绑定成功 HTTPPort[%s] use https",
			httpportstr)
		httpscertfile := this.moudleConf.
			GetProp("httpscertfile")
		httpskeyfile := this.moudleConf.
			GetProp("httpskeyfile")
		err := http.ListenAndServeTLS(httpportstr, httpscertfile,
			httpskeyfile, nil)
		if err != nil {
			panic("ListenAndServe: " + err.Error())
			return err
		}
	}
	return nil
}

// cpu性能测试
func (this *SubnetManger) startTestCpuProfile() {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := util.GetPanicInfo(recover()); err != nil {
			log.Error("[startTestCpuProfile] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
		}
	}()
	log.Debug("[SubNetManager.startTestCpuProfile] " +
		"[性能分析] StartTestCpuProfile start")
	filename := this.moudleConf.GetProp("profile_filename")
	testtime := this.moudleConf.GetPropInt("profile_time")
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
func (this *SubnetManger) SignalListen(manager ServerManager) {
	c := make(chan os.Signal, 10)
	signal.Notify(c, syscall.SIGINT,
		syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1,
		syscall.SIGUSR2)
	for {
		s := <-c
		log.Debug("[SubNetManager.SignalListen] "+
			"Get signal Signal[%d]", s)
		manager.OnSignal(s)
		switch s {
		case syscall.SIGUSR1:
			go this.startTestCpuProfile()
		case syscall.SIGUSR2:
		case syscall.SIGTERM:
		case syscall.SIGINT:
			// //收到信号后的处理
			// if this.moudleConf.Myserverinfo.Servertype !=
			// 	def.TypeSuperServer {
			// 	// 通知我的主服务器，退出连接了
			// 	sendmsg := &comm.SLogoutCommand{}
			// 	this.clientManager.BroadcastByType(def.TypeSuperServer,
			// 		sendmsg)
			// 	// 发送给所有连接到我的服务器，我要关闭了，别再尝试连接我了
			// 	this.taskManager.BroadcastAll(sendmsg)
			// }
			this.moudleConf.TerminateServer = true
		case syscall.SIGQUIT:
			// 捕捉到就退不出了
			buf := make([]byte, 1<<20)
			stacklen := runtime.Stack(buf, true)
			log.Debug("[SubNetManager.SignalListen] "+
				"Received SIGQUIT, \n: Stack[%s]", buf[:stacklen])
		}
	}
}
