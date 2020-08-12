package connect

import (
	"fmt"
	"net"
	"time"

	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/session"
)

// Client 一个客户端连接，一般由 Gateway 创建
type Client struct {
	BaseConnect

	// 会话信息 可在不同服务器之间同步的
	*session.Session

	// 接收消息通道
	readch chan *msg.MessageBinary
	// 回调
	connHook IConnectHook
}

// InitTCP Initial a new client
// netconn: 连接的net.Conn对象
func (c *Client) InitTCP(netconn net.Conn, connHook IConnectHook) {
	c.BaseConnect.Init()
	c.Session = &session.Session{}
	c.IConnection = NewTCP(netconn, c.Logger,
		ClientConnSendChanSize, ClientConnSendBufferSize,
		ClientConnRecvChanSize, ClientConnRecvBufferSize)
	if c.Logger != nil {
		c.Logger.SetTopic(fmt.Sprintf("Client:%s(%s)", c.IConnection.RemoteAddr(), c.GetTempID()))
	}
	// 客户端连接的连接ID就是该连接的TmpID
	c.Session.SetConnectID(c.GetTempID())
	c.readch = c.GetRecvMessageChannel()
	c.connHook = connHook
	go c.recvMsgThread()
}

// DialTCP 为该客户端建立一个 TCP 连接
func (c *Client) DialTCP(addr string, connHook IConnectHook) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	c.InitTCP(conn, connHook)
	c.StartRecv()
	return nil
}

// onRecvMessage 当收到一个消息时调用
func (c *Client) onRecvMessage(msg *msg.MessageBinary) {
	if c.connHook != nil {
		c.connHook.OnRecvConnectMessage(c, msg)
	}
}

// onClose 当客户端连接关闭时调用
func (c *Client) onClose() {
	if c.connHook != nil {
		c.connHook.OnConnectClose(c)
	}
}

// recvMsgThread 接收消息线程
func (c *Client) recvMsgThread() {
	defer func() {
		c.onClose()
	}()

	for {
		select {
		case m, ok := <-c.readch:
			if !ok || m == nil {
				return
			}
			c.onRecvMessage(m)
		}
	}
}

// Check 返回连接是否仍可用
func (c *Client) Check() bool {
	curtime := time.Now().Unix()
	// 检查本服务器时候还存活
	if c.IsTerminateForce() {
		// 本服务器关闭
		c.Debug("[Client.Check] We initiated the disconnection")
		// 强制移除客户端连接
		return false
	}
	// 检查客户端连接是否验证超时
	if c.IsTerminateTimeout(curtime) {
		// 客户端超时未通过验证
		if !c.Session.IsVerify() {
			c.Debug("[Client.Check] Prolonged failure to verify, disconnect")
		} else {
			c.Debug("[Client.Check] Long periods of inactivity, disconnect")
		}
		return false
	}
	return true
}
