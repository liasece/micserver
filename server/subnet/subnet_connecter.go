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

func (this *SubnetManager) tryConnectServerThread(id string, addr string) {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := sysutil.GetPanicInfo(recover()); err != nil {
			this.Error("[tryConnectServerThread] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
		}
	}()

	for true {
		this.connectMutex.Lock()
		c := this.serverexitchan[id]
		this.connectMutex.Unlock()
		select {
		case <-c:
			this.Syslog("[SubnetManager.tryConnectServerThread] "+
				"正在连接 ModuleID[%s] IPPort[%s]",
				id, addr)
			err := this.ConnectServer(id, addr)
			if err != nil && err.Error() != "重复连接" {
				time.Sleep(2 * time.Second) // 2秒重连一次
				c <- true
			} else {
				time.Sleep(1 * time.Second) // 1秒重连一次
			}
		}
	}
}

// 这种连接不会跟着super一起停机
func (this *SubnetManager) TryConnectServer(id string, addr string) {
	this.connectMutex.Lock()
	defer this.connectMutex.Unlock()
	if this.serverexitchan == nil {
		this.serverexitchan = make(map[string]chan bool)
	}
	if _, finded := this.serverexitchan[id]; !finded {
		this.serverexitchan[id] = make(chan bool, 100)
	} else {
		this.Syslog("[SubnetManager.TryConnectServer] "+
			"ModuleID[%s] 守护线程已启动，不再重复启动",
			id)
		return
	}
	this.serverexitchan[id] <- true
	go this.tryConnectServerThread(id, addr)
}

// 连接服务器
func (this *SubnetManager) ConnectServer(id string,
	addr string) error {
	this.connectMutex.Lock()
	defer this.connectMutex.Unlock()
	oldconn := this.GetServer(id)
	// 重复连接
	if oldconn != nil {
		this.Syslog("[SubnetManager.ConnectServer] "+
			"ModuleID[%s] 重复的连接", id)
		return errors.New("重复连接")
	}
	// this.Syslog("[SubnetManager.ConnectServer] "+
	// 	"服务器连接创建地址开始 ModuleID[%s] ServerIPPort[%s]",
	// 	id, addr)
	if chanServer := process.GetServerChan(id); chanServer != nil {
		newMsgChan := make(chan *msg.MessageBinary, 1000)
		chanServer <- &process.ChanServerHandshake{
			ModuleInfo:    this.myServerInfo,
			ServerMsgChan: nil,
			ClientMsgChan: newMsgChan,
			Seq:           0,
		}
	} else {
		tcpaddr, err := net.ResolveTCPAddr("tcp4", addr)
		if err != nil {
			this.Syslog("[SubnetManager.ConnectServer] "+
				"服务器连接创建地址失败 ServerIPPort[%s] Err[%s]",
				addr, err.Error())
			return err
		}
		netconn, err := net.DialTCP("tcp", nil, tcpaddr)
		if err != nil {
			this.Error("[SubnetManager.ConnectServer] "+
				"服务器连接失败 ServerIPPort[%s] Err[%s]",
				addr, err.Error())
			return err
		}
		this.doConnectTCPServer(netconn, id)
	}

	// this.Syslog("[SubnetManager.ConnectServer] "+
	// 	"开始连接服务器 ModuleID[%s] IPPort[%s]", id,
	// 	addr)

	return nil
}

// 使用一个TCP连接实际连接一个服务器
func (this *SubnetManager) doConnectTCPServer(netconn net.Conn, id string) {
	this.Syslog("开始登陆TCP服务器 Server:%s", id)
	conn := this.NewTCPServer(connect.ServerSCTypeClient, netconn, id,
		this.onConnectRecv, this.onConnectClose)
	conn.Logger = this.Logger
	this.OnCreateNewServer(conn)
	// 发起登录
	this.onClientConnected(conn)
}

// 使用一个本地连接实际连接一个服务器
func (this *SubnetManager) doConnectChanServer(
	sendchan, recvchan chan *msg.MessageBinary, id string) {
	this.Syslog("开始登陆Chan服务器 Server:%s", id)
	conn := this.NewChanServer(connect.ServerSCTypeClient, sendchan, recvchan, id,
		this.onConnectRecv, this.onConnectClose)
	conn.Logger = this.Logger
	this.OnCreateNewServer(conn)
	// 发起登录
	this.onClientConnected(conn)
}

func (this *SubnetManager) onClientConnected(conn *connect.Server) {
	// 开始请求登陆
	// 构造登陆消息
	sendmsg := &servercomm.SLoginCommand{}
	sendmsg.ModuleID = this.myServerInfo.ModuleID
	sendmsg.ModuleAddr = this.moudleConf.GetString(conf.SubnetTCPAddr)
	sendmsg.ConnectPriority = conn.ConnectPriority
	// 发送登陆请求
	conn.SendCmd(sendmsg)
	this.Syslog("请求登陆 Server:%s", conn.GetTempID())
}

func (this *SubnetManager) onClientDisconnected(conn *connect.Server) {
	this.onConnectClose(conn)
	this.RemoveServer(conn.GetTempID())

	if !conn.IsNormalDisconnect &&
		conn.GetSCType() == connect.ServerSCTypeClient {
		this.connectMutex.Lock()
		defer this.connectMutex.Unlock()
		if this.serverexitchan[fmt.Sprint(conn.ModuleInfo.ModuleID)] != nil {
			this.serverexitchan[fmt.Sprint(conn.ModuleInfo.ModuleID)] <- true
			this.Warn("[onClientDisconnected] "+
				"服务服务器断开连接,准备重新连接 ModuleID[%s]",
				conn.ModuleInfo.ModuleID)
		} else {
			this.Syslog("[onClientDisconnected] "+
				"服务器重连管道已关闭,取消重连 ModuleID[%s]",
				fmt.Sprint(conn.ModuleInfo.ModuleID))
		}
	}
}
