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
	"github.com/liasece/micserver/connect"
	"github.com/liasece/micserver/servercomm"
	"github.com/liasece/micserver/util/sysutil"
	"net"
	"time"
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
			this.Debug("[SubnetManager.tryConnectServerThread] "+
				"正在连接 ServerID[%s] IPPort[%s]",
				id, addr)
			_, err := this.ConnectServer(id, addr)
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
		this.Debug("[SubnetManager.TryConnectServer] "+
			"ServerID[%s] 守护线程已启动，不再重复启动",
			id)
		return
	}
	this.serverexitchan[id] <- true
	go this.tryConnectServerThread(id, addr)
}

// 连接服务器
func (this *SubnetManager) ConnectServer(id string,
	addr string) (*connect.Server, error) {
	this.connectMutex.Lock()
	defer this.connectMutex.Unlock()
	oldconn := this.GetServer(id)
	// 重复连接
	if oldconn != nil {
		this.Debug("[SubnetManager.ConnectServer] "+
			"ServerID[%s] 重复的连接", id)
		return nil, errors.New("重复连接")
	}
	this.Debug("[SubnetManager.ConnectServer] "+
		"服务器连接创建地址开始 ServerID[%s] ServerIPPort[%s]",
		id, addr)
	tcpaddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		this.Debug("[SubnetManager.ConnectServer] "+
			"服务器连接创建地址失败 ServerIPPort[%s] Err[%s]",
			addr, err.Error())
		return nil, err
	}
	Conn, err := net.DialTCP("tcp", nil, tcpaddr)
	if err != nil {
		this.Debug("[SubnetManager.ConnectServer] "+
			"服务器连接失败 ServerIPPort[%s] Err[%s]",
			addr, err.Error())
		return nil, err
	}
	conn := this.NewTCPServer(connect.ServerSCTypeClient, Conn, id,
		this.onConnectRecv, this.onConnectClose)
	conn.Logger = this.Logger
	// 发起登录

	// 构造登陆消息
	sendmsg := &servercomm.SLoginCommand{}
	sendmsg.ServerID = this.moudleConf.ID
	sendmsg.ServerAddr = this.moudleConf.GetModuleSetting("subnettcpaddr")
	sendmsg.ConnectPriority = conn.ConnectPriority
	// 发送登陆请求
	conn.SendCmd(sendmsg)

	this.Debug("[SubnetManager.ConnectServer] "+
		"开始连接服务器 ServerID[%s] IPPort[%s]", id,
		addr)

	this.OnCreateTCPConnect(conn)
	return conn, nil
}

func (this *SubnetManager) onClientDisconnected(conn *connect.Server) {
	this.onConnectClose(conn)
	this.RemoveServer(conn.Tempid)

	if !conn.IsNormalDisconnect &&
		conn.GetSCType() == connect.ServerSCTypeClient {
		this.connectMutex.Lock()
		defer this.connectMutex.Unlock()
		if this.serverexitchan[fmt.Sprint(conn.ServerInfo.ServerID)] != nil {
			this.serverexitchan[fmt.Sprint(conn.ServerInfo.ServerID)] <- true
			this.Warn("[onClientDisconnected] "+
				"服务服务器断开连接,准备重新连接 ServerID[%s]",
				conn.ServerInfo.ServerID)
		} else {
			this.Debug("[onClientDisconnected] "+
				"服务器重连管道已关闭,取消重连 ServerID[%s]",
				fmt.Sprint(conn.ServerInfo.ServerID))
		}
	}
}
