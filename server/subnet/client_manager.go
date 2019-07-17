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
	"errors"
	"fmt"
	"github.com/liasece/micserver/comm"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/tcpconn"
	"github.com/liasece/micserver/util"
	"net"
	"time"
)

func (this *GBTCPClientManager) tryConnectServerThread(serverid uint32,
	clienter IServerHandler) {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := util.GetPanicInfo(recover()); err != nil {
			log.Error("[tryConnectServerThread] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
		}
	}()
	serverinfo := this.subnetManager.TopConfigs.GetTopConfigByID(serverid)
	if serverinfo.Serverid != serverid {
		log.Error("[GBTCPClientManager.tryConnectServerThread] "+
			"tryConnectServerThread 错误的 ServerID[%d]", serverid)
		return
	}

	for true {
		this.serverexitchanmutex.Lock()
		c := this.serverexitchan[serverid]
		this.serverexitchanmutex.Unlock()
		select {
		case <-c:
			this.serverexitchanmutex.Lock()
			serverinfo := this.subnetManager.TopConfigs.
				GetTopConfigByID(serverid)
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
	oldclient := this.GetTCPClient(uint64(serverid))
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
	client := this.AddTCPClient(Conn,
		uint64(serverid))
	clienter.OnCreateTCPConnect(client)

	// 发起登录

	// 构造登陆消息
	sendmsg := &comm.SLoginCommand{}
	// sendmsg.Servertype = this.subnetManager.moudleConf.
	// 	Myserverinfo.Servertype
	// sendmsg.Serverip = this.subnetManager.moudleConf.Myserverinfo.Serverip
	// sendmsg.Servername = this.subnetManager.moudleConf.Myservername
	// sendmsg.Serverport = this.subnetManager.moudleConf.
	// 	Myserverinfo.Serverport
	// sendmsg.ServerNumber = this.subnetManager.moudleConf.
	// 	GetPropUint("servernumber")
	// sendmsg.Version = this.subnetManager.moudleConf.Version
	// if this.subnetManager.moudleConf.Myserverinfo.Serverid > 0 {
	// 	sendmsg.Serverid = this.subnetManager.moudleConf.
	// 		Myserverinfo.Serverid
	// }
	// 发送登陆请求
	client.SendCmd(sendmsg)

	log.Debug("[GBTCPClientManager.ConnectServer] "+
		"开始连接服务器 ServerID[%d] IP[%s] Port[%d]", serverid,
		serverip, serverport)

	// 监听处理消息
	go this.handleClientConnection(client, clienter)
	return client, nil
}

func (this *GBTCPClientManager) onClientDisconnected(client *tcpconn.ServerConn,
	clienter IServerHandler) {
	clienter.OnRemoveTCPConnect(client)
	this.RemoveTCPClient(client.Tempid)
	superserverid := this.subnetManager.moudleConf.
		GetPropUint("superserverid")
	if client.Tempid == uint64(superserverid) {
		this.superexitchan <- true
		log.Warn("[onClientDisconnected] " +
			"super断开连接,准备重新连接super")
	} else {
		this.serverexitchanmutex.Lock()
		defer this.serverexitchanmutex.Unlock()
		if this.
			serverexitchan[uint32(client.Serverinfo.Serverid)] != nil {
			this.serverexitchan[uint32(client.Serverinfo.Serverid)] <- true
			log.Warn("[onClientDisconnected] " +
				"服务服务器断开连接,准备重新连接")
		} else {
			log.Debug("[onClientDisconnected] "+
				"服务器重连管道已关闭,取消重连 ServerID[%d]",
				client.Serverinfo.Serverid)
		}
	}
}

func (this *GBTCPClientManager) connectRelyServers(serverinfos []comm.SServerInfo,
	clienter IServerHandler) {
	// 依赖服务器启动成功了
	for _, serverinfo := range serverinfos {
		log.Debug("[GBTCPClientManager.connectRelyServers] "+
			"收到依赖服务器信息成功"+
			" ServerID[%d] IP[%s] Port[%d] HTTPPort[%d] Name[%s]",
			serverinfo.Serverid, serverinfo.Serverip,
			serverinfo.Serverport, serverinfo.Httpport,
			serverinfo.Servername)
		this.subnetManager.TopConfigs.AddTopConfig(serverinfo)
		if this.GetTCPClient(uint64(serverinfo.Serverid)) == nil {
			log.Debug("[GBTCPClientManager.connectRelyServers] "+
				"尝试保持与 %d 的连接", serverinfo.Serverid)
			// 连接依赖服务器
			this.TryConnectServer(
				serverinfo.Serverid, clienter)
		} else {
			log.Debug("[GBTCPClientManager.connectRelyServers] "+
				"ServerID[%d] 已连接，不再重复连接", serverinfo.Serverid)
		}
	}
}
