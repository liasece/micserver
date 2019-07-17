package subnet

import (
	"fmt"
	"github.com/liasece/micserver/comm"
	"github.com/liasece/micserver/conf"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/server/subnet/serconfs"
	"github.com/liasece/micserver/tcpconn"
	"net/http"
	"os"
	"strings"
	"sync"
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
type SubnetManager struct {
	serverhandler IServerHandler
	servermanger  ServerManager

	clientManager *GBTCPClientManager
	taskManager   *GBTCPTaskManager

	TopConfigs serconfs.TopConfigsManager // 所有服务器信息

	moudleConf *conf.TopConfig
}

func (this *SubnetManager) InitManager() {

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

func (this *SubnetManager) GetClientManager() *GBTCPClientManager {
	return this.clientManager
}

func (this *SubnetManager) GetTaskManager() *GBTCPTaskManager {
	return this.taskManager
}

func (this *SubnetManager) GetLatestVersionTopConfigByType(servertype uint32) uint64 {
	latestVersion := uint64(0)
	this.TopConfigs.RangeTopConfig(
		func(value comm.SServerInfo) bool {
			if value.Servertype == servertype &&
				value.Version > latestVersion {
				latestVersion = value.Version
			}
			return true
		})
	return latestVersion
}

func (this *SubnetManager) GetServerManger() ServerManager {
	return this.servermanger
}

func (this *SubnetManager) GetServerHandler() IServerHandler {
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

var serverLoginMutex sync.Mutex

func (this *SubnetManager) OnServerLogin(tcptask *tcpconn.ServerConn,
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
	var TopConfig comm.SServerInfo

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
	TopConfig = this.TopConfigs.
		GetTopConfigByInfo(tarinfo)
	// 如果成功获取到了一个serverid
	if TopConfig.Serverid != 0 {
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
	if TopConfig.Serverid == 0 || !sameip ||
		TopConfig.Servertype != tarinfo.Servertype {
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

	TopConfig.Version = tarinfo.Version

	// 来源服务器检查完毕
	// 完善来源服务器在本服务器的信息
	tcptask.SetVertify(true)
	tcptask.SetTerminateTime(0) // 清除终止时间状态
	tcptask.Serverinfo = TopConfig
	this.taskManager.ChangeTCPTaskTempid(tcptask,
		uint64(tcptask.Serverinfo.Serverid))
	log.Debug("[SubNetManager.OnServerLogin] "+
		"客户端连接验证成功 "+
		" SerID[%d] IP[%s] Type[%d] Port[%d] Name[%s]",
		TopConfig.Serverid, TopConfig.Serverip,
		TopConfig.Servertype, TopConfig.Serverport,
		TopConfig.Servername)
	// 向来源服务器回复登陆成功消息
	retmsg := &comm.SLoginRetCommand{}
	retmsg.Loginfailed = 0
	retmsg.Clientinfo = TopConfig
	// retmsg.Taskinfo = this.moudleConf.Myserverinfo
	// retmsg.Redisinfo = this.moudleConf.RedisConfig
	tcptask.SendCmd(retmsg)

	// 通知其他服务器，次服务器登陆完成
	// 如果我是SuperServer
	// 向来源服务器发送本地已有的所有服务器信息
	this.taskManager.NotifyAllServerInfo(tcptask)
	// 把来源服务器信息广播给其它所有服务器
	notifymsg := &comm.SStartMyNotifyCommand{}
	notifymsg.Serverinfo = TopConfig
	this.taskManager.BroadcastAll(notifymsg)
}

// 绑定我的HTTP服务 会阻塞
func (this *SubnetManager) BindMyHttpServer() {
	usehttps := this.moudleConf.GetProp("usehttps")
	httpportstr := fmt.Sprintf(":%d", 8080)
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

// // 绑定我的HTTPS服务 会阻塞
// func (this *SubnetManager) BindMyHttpsServer() error {
// 	if this.moudleConf.Myserverinfo.Httpsport > 0 {
// 		httpportstr := fmt.Sprintf(":%d", 8080)
// 		log.Debug("[SubNetManager.BindMyHttpsServer] "+
// 			"服务器https绑定成功 HTTPPort[%s] use https",
// 			httpportstr)
// 		httpscertfile := this.moudleConf.
// 			GetProp("httpscertfile")
// 		httpskeyfile := this.moudleConf.
// 			GetProp("httpskeyfile")
// 		err := http.ListenAndServeTLS(httpportstr, httpscertfile,
// 			httpskeyfile, nil)
// 		if err != nil {
// 			panic("ListenAndServe: " + err.Error())
// 			return err
// 		}
// 	}
// 	return nil
// }
