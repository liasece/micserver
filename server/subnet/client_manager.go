/**
 * \file GBTCPClientManager.go
 * \version
 * \author wzy
 * \date  2018年01月15日 18:22:43
 * \brief client连接管理器
 *
 */

package subnet

import (
	"base"
	"github.com/liasece/micserver/def"
	// "base/functime"
	"errors"
	"fmt"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/tcpconn"
	"github.com/liasece/micserver/util"
	// "io"
	"net"
	// "runtime"
	"github.com/liasece/micserver/servercomm"
	// "sync"
	"time"
)

// 连接到Super服务器
func (this *GBTCPClientManager) connectSuper(clienter IServerHandler) {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := util.GetPanicInfo(recover()); err != nil {
			log.Error("[connectSuper] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
		}
	}()

	// 获取 SuperServer 信息
	superserverip := base.GetGBServerConfigM().
		GetProp("superserverip")
	superserverid := base.GetGBServerConfigM().
		GetPropUint("superserverid")
	superserverport := base.GetGBServerConfigM().
		GetPropUint("superserverport")

	this.superexitchan <- true
	for !base.GetGBServerConfigM().TerminateServer {
		select {
		case <-this.superexitchan:
			client, err := this.ConnectServer(superserverid,
				superserverip, superserverport, clienter)
			if err != nil {
				time.Sleep(2 * time.Second) // 2秒重连一次
				this.superexitchan <- true
			} else {
				client.Serverinfo.Servertype = def.TypeSuperServer
			}
		}
	}
}

// 这种连接不会跟着super一起停机
func (this *GBTCPClientManager) TryConnectSuperServer(
	clienter IServerHandler) error {
	superserverip := base.GetGBServerConfigM().
		GetProp("superserverip")
	superserverid := base.GetGBServerConfigM().
		GetPropUint("superserverid")
	superserverport := base.GetGBServerConfigM().
		GetPropUint("superserverport")
	if len(superserverip) == 0 || superserverid == 0 ||
		superserverport == 0 {
		return errors.New("error server ip/id/port")
	}

	time.Sleep(1 * time.Second) // 防止拥挤
	go this.connectSuper(clienter)
	// 暂停等待获取本机信息
	for {
		if base.GetGBServerConfigM().Myserverinfo.Serverid > 0 {
			// 登陆Super完成
			// 服务器名字可能会发生改变
			// 重新加载一次配置文件
			base.GetGBServerConfigM().ReloadConfig()
			break
		}
		// 休眠100毫秒
		time.Sleep(200 * time.Millisecond)
	}
	return nil
}

func (this *GBTCPClientManager) tryConnectServerThread(serverid uint32,
	clienter IServerHandler) {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := util.GetPanicInfo(recover()); err != nil {
			log.Error("[tryConnectServerThread] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
		}
	}()
	serverinfo := GetSubnetManager().ServerConfigs.GetServerConfigByID(serverid)
	if serverinfo.Serverid != serverid {
		log.Error("[GBTCPClientManager.tryConnectServerThread] "+
			"tryConnectServerThread 错误的 ServerID[%d]", serverid)
		return
	}

	for !base.GetGBServerConfigM().TerminateServer {
		this.serverexitchanmutex.Lock()
		c := this.serverexitchan[serverid]
		this.serverexitchanmutex.Unlock()
		select {
		case <-c:
			this.serverexitchanmutex.Lock()
			serverinfo := GetSubnetManager().ServerConfigs.
				GetServerConfigByID(serverid)
			if serverinfo.Serverid != serverid {
				log.Debug("[GBTCPClientManager.tryConnectServerThread] "+
					"本地已删除该服务器信息，取消连接 ServerID[%d] Info[%s]",
					serverid, serverinfo.GetJson())
				close(this.serverexitchan[serverid])
				this.serverexitchan[serverid] = nil
				this.serverexitchanmutex.Unlock()
				return
			}
			this.serverexitchanmutex.Unlock()
			serverip := serverinfo.Serverip
			log.Debug("[GBTCPClientManager.tryConnectServerThread] "+
				"正在连接 ServerID[%d] IPPort[%s:%d]",
				serverid, serverip, serverinfo.Serverport)
			_, err := this.ConnectServer(serverid,
				serverip, serverinfo.Serverport, clienter)
			if err != nil {
				time.Sleep(2 * time.Second) // 2秒重连一次
				c <- true
			} else {
				time.Sleep(1 * time.Second) // 1秒重连一次
			}
		}
	}
}

// 这种连接不会跟着super一起停机
func (this *GBTCPClientManager) TryConnectServer(serverid uint32,
	clienter IServerHandler) {
	this.serverexitchanmutex.Lock()
	defer this.serverexitchanmutex.Unlock()
	if this.serverexitchan == nil {
		this.serverexitchan = make(map[uint32]chan bool)
	}
	if _, finded := this.serverexitchan[serverid]; !finded {
		this.serverexitchan[serverid] = make(chan bool, 100)
	} else if this.serverexitchan[serverid] == nil {
		this.serverexitchan[serverid] = make(chan bool, 100)
	} else {
		log.Debug("[GBTCPClientManager.TryConnectServer] "+
			"ServerID[%d] 守护线程已启动，不再重复启动",
			serverid)
		return
	}
	this.serverexitchan[serverid] <- true
	go this.tryConnectServerThread(serverid, clienter)
}

// 连接服务器
func (this *GBTCPClientManager) ConnectServer(serverid uint32,
	serverip string, serverport uint32,
	clienter IServerHandler) (*tcpconn.ServerConn, error) {
	this.connectMutex.Lock()
	defer this.connectMutex.Unlock()
	oldclient := GetGBTCPClientManager().GetTCPClient(uint64(serverid))
	// 重复连接
	if oldclient != nil {
		log.Error("[GBTCPClientManager.ConnectServer] "+
			"ServerID[%d] 重复的连接", serverid)
		return nil, errors.New("重复连接")
	}
	log.Debug("[GBTCPClientManager.ConnectServer] "+
		"服务器连接创建地址开始 ServerID[%d] ServerIP[%s] Port[%d]",
		serverid, serverip, serverport)
	serverinfo := fmt.Sprintf("%s:%d", serverip, serverport)
	tcpaddr, err := net.ResolveTCPAddr("tcp4", serverinfo)
	if err != nil {
		log.Debug("[GBTCPClientManager.ConnectServer] "+
			"服务器连接创建地址失败 ServerIP[%s] Port[%d] Err[%s]",
			serverip, serverport, err.Error())
		return nil, err
	}
	Conn, err := net.DialTCP("tcp", nil, tcpaddr)
	if err != nil {
		log.Debug("[GBTCPClientManager.ConnectServer] "+
			"服务器连接失败 ServerIP[%s] Port[%d] Err[%s]",
			serverip, serverport, err.Error())
		return nil, err
	}
	this.serverhandler = clienter
	client := GetGBTCPClientManager().AddTCPClient(Conn,
		uint64(serverid))
	clienter.OnCreateTCPConnect(client)

	// 发起登录

	// 构造登陆消息
	sendmsg := &servercomm.SLoginCommand{}
	sendmsg.Servertype = base.GetGBServerConfigM().
		Myserverinfo.Servertype
	sendmsg.Serverip = base.GetGBServerConfigM().Myserverinfo.Serverip
	sendmsg.Servername = base.GetGBServerConfigM().Myservername
	sendmsg.Serverport = base.GetGBServerConfigM().
		Myserverinfo.Serverport
	sendmsg.ServerNumber = base.GetGBServerConfigM().
		GetPropUint("servernumber")
	sendmsg.Version = base.GetGBServerConfigM().Version
	if base.GetGBServerConfigM().Myserverinfo.Serverid > 0 {
		sendmsg.Serverid = base.GetGBServerConfigM().
			Myserverinfo.Serverid
	}
	// 发送登陆请求
	client.SendCmd(sendmsg)

	log.Debug("[GBTCPClientManager.ConnectServer] "+
		"开始连接服务器 ServerID[%d] IP[%s] Port[%d]", serverid,
		serverip, serverport)

	// 监听处理消息
	go handleClientConnection(client, clienter)
	return client, nil
}

func onClientDisconnected(client *tcpconn.ServerConn,
	clienter IServerHandler) {
	clienter.OnRemoveTCPConnect(client)
	GetGBTCPClientManager().RemoveTCPClient(client.Tempid)
	superserverid := base.GetGBServerConfigM().
		GetPropUint("superserverid")
	if client.Tempid == uint64(superserverid) {
		GetGBTCPClientManager().superexitchan <- true
		log.Warn("[onClientDisconnected] " +
			"super断开连接,准备重新连接super")
	} else {
		GetGBTCPClientManager().serverexitchanmutex.Lock()
		defer GetGBTCPClientManager().serverexitchanmutex.Unlock()
		if GetGBTCPClientManager().
			serverexitchan[uint32(client.Serverinfo.Serverid)] != nil {
			GetGBTCPClientManager().
				serverexitchan[uint32(client.Serverinfo.Serverid)] <- true
			log.Warn("[onClientDisconnected] " +
				"服务服务器断开连接,准备重新连接")
		} else {
			log.Debug("[onClientDisconnected] "+
				"服务器重连管道已关闭,取消重连 ServerID[%d]",
				client.Serverinfo.Serverid)
		}
	}
}

func connectRelyServers(serverinfos []servercomm.SServerInfo,
	clienter IServerHandler) {
	// 依赖服务器启动成功了
	for _, serverinfo := range serverinfos {
		log.Debug("[GBTCPClientManager.connectRelyServers] "+
			"收到依赖服务器信息成功"+
			" ServerID[%d] IP[%s] Port[%d] HTTPPort[%d] Name[%s]",
			serverinfo.Serverid, serverinfo.Serverip,
			serverinfo.Serverport, serverinfo.Httpport,
			serverinfo.Servername)
		GetSubnetManager().ServerConfigs.AddServerConfig(serverinfo)
		if GetGBTCPClientManager().
			GetTCPClient(uint64(serverinfo.Serverid)) == nil {
			log.Debug("[GBTCPClientManager.connectRelyServers] "+
				"尝试保持与 %d 的连接", serverinfo.Serverid)
			// 连接依赖服务器
			GetGBTCPClientManager().TryConnectServer(
				serverinfo.Serverid, clienter)
		} else {
			log.Debug("[GBTCPClientManager.connectRelyServers] "+
				"ServerID[%d] 已连接，不再重复连接", serverinfo.Serverid)
		}
	}
}
