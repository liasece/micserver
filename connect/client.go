package connect

import (
	"fmt"
	"net"
	"time"

	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/session"
)

// 一个客户端连接，一般由 Gateway 创建
type Client struct {
	BaseConnect

	// 会话信息 可在不同服务器之间同步的
	*session.Session

	// 接收消息通道
	readch chan *msg.MessageBinary
	// 回调
	connHook ConnectHook
}

// Initial a new client
// netconn: 连接的net.Conn对象
func (this *Client) InitTCP(netconn net.Conn, connHook ConnectHook) {
	this.BaseConnect.Init()
	this.Session = &session.Session{}
	this.IConnection = NewTCP(netconn, this.Logger,
		ClientConnSendChanSize, ClientConnSendBufferSize,
		ClientConnRecvChanSize, ClientConnRecvBufferSize)
	if this.Logger != nil {
		this.Logger.SetTopic(fmt.Sprintf("Client:%s(%s)",
			this.IConnection.RemoteAddr(), this.GetTempID()))
	}
	// 客户端连接的连接ID就是该连接的TmpID
	this.Session.SetConnectID(this.GetTempID())
	this.readch = this.GetRecvMessageChannel()
	this.connHook = connHook
	go this.recvMsgThread()
}

// 为该客户端建立一个 TCP 连接
func (this *Client) DialTCP(addr string, connHook ConnectHook) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	this.InitTCP(conn, connHook)
	this.StartRecv()
	return nil
}

// 当收到一个消息时调用
func (this *Client) onRecvMessage(msg *msg.MessageBinary) {
	if this.connHook != nil {
		this.connHook.OnRecvConnectMessage(this, msg)
	}
}

// 当客户端连接关闭时调用
func (this *Client) onClose() {
	if this.connHook != nil {
		this.connHook.OnConnectClose(this)
	}
}

// 接收消息线程
func (this *Client) recvMsgThread() {
	defer func() {
		this.onClose()
	}()

	for {
		select {
		case m, ok := <-this.readch:
			if !ok || m == nil {
				return
			}
			this.onRecvMessage(m)
		}
	}
}

// 返回连接是否仍可用
func (this *Client) Check() bool {
	curtime := time.Now().Unix()
	// 检查本服务器时候还存活
	if this.IsTerminateForce() {
		// 本服务器关闭
		this.Debug("[Client.Check] 服务器强制断开连接")
		// 强制移除客户端连接
		return false
	}
	// 检查客户端连接是否验证超时
	if this.IsTerminateTimeout(curtime) {
		// 客户端超时未通过验证
		if !this.Session.IsVertify() {
			this.Debug("[Client.Check] 长时间未通过验证，断开连接")
		} else {
			this.Debug("[Client.Check] 长时间未活动，断开连接")
		}
		return false
	}
	return true
}
