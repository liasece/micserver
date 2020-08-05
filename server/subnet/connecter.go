/**
 * \file SubnetManager.go
 * \version
 * \author wzy
 * \date  2018年01月15日 18:22:43
 * \brief conn连接管理器
 *
 */

package subnet

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/liasece/micserver/conf"
	"github.com/liasece/micserver/connect"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/process"
	"github.com/liasece/micserver/servercomm"
	"github.com/liasece/micserver/util/sysutil"
)

// tryConnectServerThread 保持与目标服务器的连接
func (manager *Manager) tryConnectServerThread(id string, addr string) {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if stackInfo, err := sysutil.GetPanicInfo(recover()); err != nil {
			manager.Error("[tryConnectServerThread] Panic", log.ErrorField(err), log.String("Stack", stackInfo))
		}
	}()

	for true {
		manager.connectMutex.Lock()
		c := manager.serverexitchan[id]
		manager.connectMutex.Unlock()
		select {
		case <-c:
			manager.Syslog("[Manager.tryConnectServerThread] Connecting", log.String("ModuleID", id), log.String("IPPort", addr))
			err := manager.ConnectServer(id, addr)
			if err != nil && err.Error() != "duplicate connection" {
				time.Sleep(2 * time.Second) // 2秒重连一次
				c <- true
			} else {
				time.Sleep(1 * time.Second) // 1秒重连一次
			}
		}
	}
}

// TryConnectServer 这种连接不会跟着super一起停机
func (manager *Manager) TryConnectServer(id string, addr string) {
	manager.connectMutex.Lock()
	defer manager.connectMutex.Unlock()
	if manager.serverexitchan == nil {
		manager.serverexitchan = make(map[string]chan bool)
	}
	if _, finded := manager.serverexitchan[id]; !finded {
		manager.serverexitchan[id] = make(chan bool, 100)
	} else {
		manager.Syslog("[Manager.TryConnectServer] The daemon thread has been started and will not be started again", log.String("ModuleID", id))
		return
	}
	manager.serverexitchan[id] <- true
	go manager.tryConnectServerThread(id, addr)
}

// ConnectServer 连接服务器
func (manager *Manager) ConnectServer(id string, addr string) error {
	manager.connectMutex.Lock()
	defer manager.connectMutex.Unlock()
	oldconn := manager.GetServer(id)
	// 重复连接
	if oldconn != nil {
		manager.Syslog("[Manager.ConnectServer] Duplicate connection", log.String("ModuleID", id))
		return errors.New("duplicate connection")
	}
	if chanServer := process.GetServerChan(id); chanServer != nil {
		newMsgChan := make(chan *msg.MessageBinary, 1000)
		chanServer <- &process.ChanServerHandshake{
			ModuleInfo:    manager.myServerInfo,
			ServerMsgChan: nil,
			ClientMsgChan: newMsgChan,
			Seq:           0,
		}
	} else {
		tcpaddr, err := net.ResolveTCPAddr("tcp4", addr)
		if err != nil {
			manager.Syslog("[Manager.ConnectServer] Server connection address creation failed", log.String("ServerIPPort", addr), log.ErrorField(err))
			return err
		}
		netconn, err := net.DialTCP("tcp", nil, tcpaddr)
		if err != nil {
			manager.Error("[Manager.ConnectServer] Server connection failure", log.String("ServerIPPort", addr), log.ErrorField(err))
			return err
		}
		manager.doConnectTCPServer(netconn, id)
	}

	return nil
}

// doConnectTCPServer 使用一个TCP连接实际连接一个服务器
func (manager *Manager) doConnectTCPServer(netconn net.Conn, id string) {
	manager.Syslog("Start logging into the TCP server", log.String("ServerID", id))
	conn := manager.NewTCPServer(connect.ServerSCTypeClient, netconn, id, manager.onConnectRecv, manager.onConnectClose)
	conn.Logger = manager.Logger
	manager.OnCreateNewServer(conn)
	// 发起登录
	manager.onClientConnected(conn)
}

// doConnectChanServer 使用一个本地连接实际连接一个服务器
func (manager *Manager) doConnectChanServer(
	sendchan, recvchan chan *msg.MessageBinary, id string) {
	manager.Syslog("Start logging into the Chan server", log.String("ServerID", id))
	conn := manager.NewChanServer(connect.ServerSCTypeClient, sendchan, recvchan, id, manager.onConnectRecv, manager.onConnectClose)
	conn.Logger = manager.Logger
	manager.OnCreateNewServer(conn)
	// 发起登录
	manager.onClientConnected(conn)
}

// onClientConnected 当连接到本服务器时
func (manager *Manager) onClientConnected(conn *connect.Server) {
	// 开始请求登陆
	// 构造登陆消息
	sendmsg := &servercomm.SLoginCommand{}
	sendmsg.ModuleID = manager.myServerInfo.ModuleID
	sendmsg.ModuleAddr = manager.moudleConf.GetString(conf.SubnetTCPAddr)
	sendmsg.ConnectPriority = conn.ConnectPriority
	// 发送登陆请求
	conn.SendCmd(sendmsg)
	manager.Syslog("Request login", log.String("ServerTempID", conn.GetTempID()))
}

// onClientDisconnected 当与本服务器的连接断开时
func (manager *Manager) onClientDisconnected(conn *connect.Server) {
	manager.onConnectClose(conn)
	manager.RemoveServer(conn.GetTempID())

	if !conn.IsNormalDisconnect &&
		conn.GetSCType() == connect.ServerSCTypeClient {
		manager.connectMutex.Lock()
		defer manager.connectMutex.Unlock()
		if manager.serverexitchan[fmt.Sprint(conn.ModuleInfo.ModuleID)] != nil {
			manager.serverexitchan[fmt.Sprint(conn.ModuleInfo.ModuleID)] <- true
			manager.Warn("[onClientDisconnected] Service server disconnected and ready to reconnect", log.String("ModuleID", conn.ModuleInfo.ModuleID))
		} else {
			manager.Syslog("[onClientDisconnected] The server reconnection pipeline has been closed and the reconnection cancelled", log.String("ModuleID", conn.ModuleInfo.ModuleID))
		}
	}
}
