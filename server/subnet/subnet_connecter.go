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
			manager.Error("[tryConnectServerThread] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
		}
	}()

	for true {
		manager.connectMutex.Lock()
		c := manager.serverexitchan[id]
		manager.connectMutex.Unlock()
		select {
		case <-c:
			manager.Syslog("[Manager.tryConnectServerThread] "+
				"正在连接 ModuleID[%s] IPPort[%s]",
				id, addr)
			err := manager.ConnectServer(id, addr)
			if err != nil && err.Error() != "重复连接" {
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
		manager.Syslog("[Manager.TryConnectServer] "+
			"ModuleID[%s] 守护线程已启动，不再重复启动",
			id)
		return
	}
	manager.serverexitchan[id] <- true
	go manager.tryConnectServerThread(id, addr)
}

// ConnectServer 连接服务器
func (manager *Manager) ConnectServer(id string,
	addr string) error {
	manager.connectMutex.Lock()
	defer manager.connectMutex.Unlock()
	oldconn := manager.GetServer(id)
	// 重复连接
	if oldconn != nil {
		manager.Syslog("[Manager.ConnectServer] "+
			"ModuleID[%s] 重复的连接", id)
		return errors.New("重复连接")
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
			manager.Syslog("[Manager.ConnectServer] "+
				"服务器连接创建地址失败 ServerIPPort[%s] Err[%s]",
				addr, err.Error())
			return err
		}
		netconn, err := net.DialTCP("tcp", nil, tcpaddr)
		if err != nil {
			manager.Error("[Manager.ConnectServer] "+
				"服务器连接失败 ServerIPPort[%s] Err[%s]",
				addr, err.Error())
			return err
		}
		manager.doConnectTCPServer(netconn, id)
	}

	// manager.Syslog("[Manager.ConnectServer] "+
	// 	"开始连接服务器 ModuleID[%s] IPPort[%s]", id,
	// 	addr)

	return nil
}

// doConnectTCPServer 使用一个TCP连接实际连接一个服务器
func (manager *Manager) doConnectTCPServer(netconn net.Conn, id string) {
	manager.Syslog("开始登陆TCP服务器 Server:%s", id)
	conn := manager.NewTCPServer(connect.ServerSCTypeClient, netconn, id,
		manager.onConnectRecv, manager.onConnectClose)
	conn.Logger = manager.Logger
	manager.OnCreateNewServer(conn)
	// 发起登录
	manager.onClientConnected(conn)
}

// doConnectChanServer 使用一个本地连接实际连接一个服务器
func (manager *Manager) doConnectChanServer(
	sendchan, recvchan chan *msg.MessageBinary, id string) {
	manager.Syslog("开始登陆Chan服务器 Server:%s", id)
	conn := manager.NewChanServer(connect.ServerSCTypeClient, sendchan, recvchan, id,
		manager.onConnectRecv, manager.onConnectClose)
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
	manager.Syslog("请求登陆 Server:%s", conn.GetTempID())
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
			manager.Warn("[onClientDisconnected] "+
				"服务服务器断开连接,准备重新连接 ModuleID[%s]",
				conn.ModuleInfo.ModuleID)
		} else {
			manager.Syslog("[onClientDisconnected] "+
				"服务器重连管道已关闭,取消重连 ModuleID[%s]",
				fmt.Sprint(conn.ModuleInfo.ModuleID))
		}
	}
}
