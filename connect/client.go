package connect

import (
	"fmt"
	"net"
	"time"

	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/session"
)

type Client struct {
	BaseConnect

	// 会话信息 可在不同服务器之间同步的
	session.Session

	// 接收消息通道
	readch chan *msg.MessageBinary
	// 回调
	connHook ConnectHook
}

// 客户端连接发送消息缓冲不宜过大， 10*64KiB*100000连接=64GiB
const ClientConnSendChanSize = 256

// 发送缓冲大小，用于将多个小消息拼接发送的缓冲大小
const ClientConnSendBufferSize = msg.MessageMaxSize * 2

// 客户端连接发送消息缓冲不宜过大， 10*64KiB*100000连接=64GiB
const ClientConnRecvChanSize = 256

// 发送缓冲大小，用于将多个小消息拼接发送的缓冲大小
const ClientConnRecvBufferSize = msg.MessageMaxSize * 2

// Initial a new client
// netconn: 连接的net.Conn对象
func (this *Client) InitTCP(netconn net.Conn, connHook ConnectHook) {
	this.BaseConnect.Init()
	this.IConnection = NewTCP(netconn, this.Logger,
		ClientConnSendChanSize, ClientConnSendBufferSize,
		ClientConnRecvChanSize, ClientConnRecvBufferSize)
	if this.Logger != nil {
		this.Logger.SetTopic(fmt.Sprintf("Client:%s(%s)",
			this.IConnection.RemoteAddr(), this.GetConnectID()))
	}
	this.readch = this.IConnection.GetRecvMessageChannel()
	this.connHook = connHook
	go this.recvMsgThread()
}

func (this *Client) DialTCP(addr string, connHook ConnectHook) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	this.InitTCP(conn, connHook)
	this.StartReadData()
	return nil
}

func (this *Client) StartReadData() {
	this.IConnection.StartRecv()
}

func (this *Client) onRecvMessage(msg *msg.MessageBinary) {
	if this.connHook != nil {
		this.connHook.OnRecvMessage(this, msg)
	}
}

func (this *Client) onClose() {
	if this.connHook != nil {
		this.connHook.OnClose(this)
	}
}

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

func (this *Client) SendBytes(
	cmdid uint16, protodata []byte) error {
	return this.IConnection.SendBytes(cmdid, protodata)
}

func (this *Client) SetLogger(logger *log.Logger) {
	this.Logger = logger.Clone()
}
