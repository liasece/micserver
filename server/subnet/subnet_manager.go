package subnet

import (
	"base"
	"base/logger"
	// "bytes"
	// "encoding/binary"
	// "encoding/json"
	"base/tcpconn"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"servercomm"
	// "strconv"
	// "strings"
	// "encoding/hex"
	"base/def"
	"base/subnet/serconfs"
	"base/util"
	"strings"
	"sync"
	"syscall"
	"time"
)

type ConnectMsgQueueStruct struct {
	task *tcpconn.ServerConn
	msg  *base.MessageBinary
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
}

var subnetmanager_s *SubnetManger

func init() {
	subnetmanager_s = &SubnetManger{}

	subnetmanager_s.clientManager = &GBTCPClientManager{}
	subnetmanager_s.clientManager.connPool.
		Init(tcpconn.ServerSCTypeClient, 1)
	subnetmanager_s.clientManager.superexitchan = make(chan bool, 1)

	subnetmanager_s.taskManager = &GBTCPTaskManager{}
	subnetmanager_s.taskManager.connPool.
		Init(tcpconn.ServerSCTypeTask, 2)
}

func GetSubnetManager() *SubnetManger {
	return subnetmanager_s
}

func GetGBTCPClientManager() *GBTCPClientManager {
	return GetSubnetManager().clientManager
}

func GetGBTCPTaskManager() *GBTCPTaskManager {
	return GetSubnetManager().taskManager
}

func GetLatestVersionServerConfigByType(servertype uint32) uint64 {
	latestVersion := uint64(0)
	GetSubnetManager().ServerConfigs.RangeServerConfig(
		func(value servercomm.SServerInfo) bool {
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
		msgbin *base.MessageBinary)
	GetTCPMsgParseChan(serverconn *tcpconn.ServerConn,
		maxchan int, msg *base.MessageBinary) int
}

type ServerManager interface {
	OnSignal(os.Signal)
	NotifyOtherMyInfo()
}

func StartMain(server IServerHandler, manager ServerManager) {
	logger.Debug("[SubNetManager.StartMain] " +
		"Main is start------")
	// 初始化参数
	GetSubnetManager().servermanger = manager
	GetSubnetManager().serverhandler = server

	// 绑定本地服务端口
	// 必须等待本地服务器端口绑定完成之后，再进行其他操作
	// 这是服务器加入内部网络的基础
	err := BindMyTCPServer(server)
	if err != nil {
		logger.Error("[StartMain] BindMyTCPServer Err[%s]", err)
		return
	}

	// 监听系统Signal
	go SignalListen(manager)

	// 初始化服务器
	server.OnInit()

	// 通知主服务器
	if base.GetGBServerConfigM().Myserverinfo.Servertype != def.TypeSuperServer {
		// 如果当前不是SuperServer
		// 等待之后再通知SuperServer我启动完成了
		time.Sleep(1 * time.Second)
		startokmsg := &servercomm.SSeverStartOKCommand{}
		startokmsg.Serverid = base.GetGBServerConfigM().Myserverinfo.Serverid
		GetGBTCPClientManager().BroadcastByType(def.TypeSuperServer, startokmsg)
	}

	// 保持程序运行
	for !base.GetGBServerConfigM().TerminateServer {
		time.Sleep(1 * time.Second)
	}

	// 当程序即将结束时
	server.OnFinal()
	logger.Debug("[SubNetManager.StartMain] " +
		"All server is over add save datas")

	// 等日志打完
	time.Sleep(1 * time.Second)
}

var serverLoginMutex sync.Mutex

func OnServerLogin(tcptask *tcpconn.ServerConn,
	tarinfo *servercomm.SLoginCommand) {
	serverLoginMutex.Lock()
	defer serverLoginMutex.Unlock()

	// 来源服务器请求登陆本服务器
	oldtcptask := GetGBTCPTaskManager().GetTCPTask(uint64(tarinfo.Serverid))
	if oldtcptask != nil {
		// 已经连接成功过了，非法连接
		logger.Error("[SubNetManager.OnServerLogin] "+
			"收到了重复的Server连接请求 Msg[%s]",
			tarinfo.GetJson())
		retmsg := &servercomm.SLoginRetCommand{}
		retmsg.Loginfailed = 1
		tcptask.SendCmd(retmsg)
		tcptask.Terminate()
		return
	}
	if !CheckServerType(tarinfo.Servertype) {
		// 未知的服务器类型
		logger.Error("[SubNetManager.OnServerLogin] "+
			"收到未知的服务器类型 Msg[%s]",
			tarinfo.GetJson())
		retmsg := &servercomm.SLoginRetCommand{}
		retmsg.Loginfailed = 1
		tcptask.SendCmd(retmsg)
		tcptask.Terminate()
		return
	}
	var serverconfig servercomm.SServerInfo

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
	serverconfig = GetSubnetManager().ServerConfigs.
		GetServerConfigByInfo(tarinfo)
	// 如果成功获取到了一个serverid
	if serverconfig.Serverid != 0 {
		// 检查该配置是否已被占用
		oldtcptask = GetGBTCPTaskManager().GetTCPTask(uint64(tarinfo.Serverid))
		if oldtcptask != nil {
			// 已经连接成功过了，非法连接
			logger.Error("[SubNetManager.OnServerLogin] "+
				"收到了重复的Server连接请求 2 Msg[%s]",
				tarinfo.GetJson())
			retmsg := &servercomm.SLoginRetCommand{}
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
		logger.Error("[SubNetManager.OnServerLogin] "+
			"连接分配异常 未知服务器连接 "+
			"Addr[%s] Msg[%s]",
			remoteaddr, tarinfo.GetJson())
		retmsg := &servercomm.SLoginRetCommand{}
		retmsg.Loginfailed = 1
		tcptask.SendCmd(retmsg)
		tcptask.Terminate()

		if base.GetGBServerConfigM().Myserverinfo.Servertype !=
			def.TypeSuperServer {
			requestServerInfo := &servercomm.SRequestServerInfo{}
			GetGBTCPClientManager().BroadcastByType(def.TypeSuperServer,
				requestServerInfo)
		}
		return
	}

	serverconfig.Version = tarinfo.Version

	// 来源服务器检查完毕
	// 完善来源服务器在本服务器的信息
	tcptask.SetVertify(true)
	tcptask.SetTerminateTime(0) // 清除终止时间状态
	tcptask.Serverinfo = serverconfig
	GetGBTCPTaskManager().ChangeTCPTaskTempid(tcptask,
		uint64(tcptask.Serverinfo.Serverid))
	logger.Debug("[SubNetManager.OnServerLogin] "+
		"客户端连接验证成功 "+
		" SerID[%d] IP[%s] Type[%d] Port[%d] Name[%s]",
		serverconfig.Serverid, serverconfig.Serverip,
		serverconfig.Servertype, serverconfig.Serverport,
		serverconfig.Servername)
	// 向来源服务器回复登陆成功消息
	retmsg := &servercomm.SLoginRetCommand{}
	retmsg.Loginfailed = 0
	retmsg.Clientinfo = serverconfig
	retmsg.Taskinfo = base.GetGBServerConfigM().Myserverinfo
	retmsg.Redisinfo = base.GetGBServerConfigM().RedisConfig
	tcptask.SendCmd(retmsg)

	// 通知其他服务器，次服务器登陆完成
	// 如果我是SuperServer
	// 向来源服务器发送本地已有的所有服务器信息
	GetGBTCPTaskManager().NotifyAllServerInfo(tcptask)
	// 把来源服务器信息广播给其它所有服务器
	notifymsg := &servercomm.SStartMyNotifyCommand{}
	notifymsg.Serverinfo = serverconfig
	GetGBTCPTaskManager().BroadcastAll(notifymsg)
}

// 绑定我的HTTP服务 会阻塞
func BindMyHttpServer() {
	usehttps := base.GetGBServerConfigM().GetProp("usehttps")
	if base.GetGBServerConfigM().Myserverinfo.Httpport > 0 {
		httpportstr := fmt.Sprintf(":%d",
			base.GetGBServerConfigM().Myserverinfo.Httpport)
		if usehttps == "true" {
			logger.Debug("[SubNetManager.BindMyHttpServer] "+
				"服务器https绑定成功 HTTPPort[%s] use https",
				httpportstr)
			httpscertfile := base.GetGBServerConfigM().
				GetProp("httpscertfile")
			httpskeyfile := base.GetGBServerConfigM().
				GetProp("httpskeyfile")
			err := http.ListenAndServeTLS(httpportstr, httpscertfile,
				httpskeyfile, nil)
			if err != nil {
				panic("ListenAndServe: " + err.Error())
				// return err
			}
		} else {
			logger.Debug("[SubNetManager.BindMyHttpServer] "+
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
func BindMyHttpsServer() error {
	if base.GetGBServerConfigM().Myserverinfo.Httpsport > 0 {
		httpportstr := fmt.Sprintf(":%d",
			base.GetGBServerConfigM().Myserverinfo.Httpsport)
		logger.Debug("[SubNetManager.BindMyHttpsServer] "+
			"服务器https绑定成功 HTTPPort[%s] use https",
			httpportstr)
		httpscertfile := base.GetGBServerConfigM().
			GetProp("httpscertfile")
		httpskeyfile := base.GetGBServerConfigM().
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
func startTestCpuProfile() {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := util.GetPanicInfo(recover()); err != nil {
			logger.Error("[startTestCpuProfile] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
		}
	}()
	logger.Debug("[SubNetManager.startTestCpuProfile] " +
		"[性能分析] StartTestCpuProfile start")
	filename := base.GetGBServerConfigM().GetProp("profile_filename")
	testtime := base.GetGBServerConfigM().GetPropInt("profile_time")
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	err = pprof.StartCPUProfile(f)
	if err != nil {
		logger.Error("[startTestCpuProfile] pprof.StartCPUProfile Err[%s]",
			err.Error())
		return
	}
	for i := 0; i < int(testtime); i++ {
		time.Sleep(1 * time.Second)
	}
	pprof.StopCPUProfile()
	f.Close()
	logger.Debug("[SubNetManager.startTestCpuProfile] " +
		"[性能分析] StartTestCpuProfile end")
}

// 监听系统消息
func SignalListen(manager ServerManager) {
	c := make(chan os.Signal, 10)
	signal.Notify(c, syscall.SIGINT,
		syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1,
		syscall.SIGUSR2)
	for {
		s := <-c
		logger.Debug("[SubNetManager.SignalListen] "+
			"Get signal Signal[%d]", s)
		manager.OnSignal(s)
		switch s {
		case syscall.SIGUSR1:
			go startTestCpuProfile()
		case syscall.SIGUSR2:
		case syscall.SIGTERM:
		case syscall.SIGINT:
			//收到信号后的处理
			if base.GetGBServerConfigM().Myserverinfo.Servertype !=
				def.TypeSuperServer {
				// 通知我的主服务器，退出连接了
				sendmsg := &servercomm.SLogoutCommand{}
				GetGBTCPClientManager().BroadcastByType(def.TypeSuperServer,
					sendmsg)
				// 发送给所有连接到我的服务器，我要关闭了，别再尝试连接我了
				GetGBTCPTaskManager().BroadcastAll(sendmsg)
			}
			base.GetGBServerConfigM().TerminateServer = true
		case syscall.SIGQUIT:
			// 捕捉到就退不出了
			buf := make([]byte, 1<<20)
			stacklen := runtime.Stack(buf, true)
			logger.Debug("[SubNetManager.SignalListen] "+
				"Received SIGQUIT, \n: Stack[%s]", buf[:stacklen])
		}
	}
}
