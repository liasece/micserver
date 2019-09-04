package connect

import (
	"fmt"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/session"
	"net"
	"time"
)

type Client struct {
	*log.Logger
	// 会话信息 可在不同服务器之间同步的
	session.Session
	// 连接实体
	IConnection
	// 结束时间 为0表示不结束
	terminate_time int64
	// 主动断开连接
	terminate_force bool
	// 连接创建的时间
	CreateTime int64

	// 连接的延迟信息
	ping Ping

	// 接收消息通道
	readch chan *msg.MessageBinary
	// 回调
	onRecv  func(*Client, *msg.MessageBinary)
	onClose func(*Client)
}

// 客户端连接发送消息缓冲不宜过大， 10*64KiB*100000连接=64GiB
const ClientConnSendChanSize = 256

// 发送缓冲大小，用于将多个小消息拼接发送的缓冲大小
const ClientConnSendBufferSize = msg.MessageMaxSize * 2

// 客户端连接发送消息缓冲不宜过大， 10*64KiB*100000连接=64GiB
const ClientConnRecvChanSize = 256

// 发送缓冲大小，用于将多个小消息拼接发送的缓冲大小
const ClientConnRecvBufferSize = msg.MessageMaxSize * 2

// 获取一个新的服务器连接
// sctype: 连接的 客户端/服务器 类型
// netconn: 连接的net.Conn对象
func NewClient(netconn net.Conn,
	onRecv func(*Client, *msg.MessageBinary),
	onClose func(*Client)) *Client {
	// 新建一个客户端连接
	conn := new(Client)
	conn.IConnection = NewTCP(netconn,
		ClientConnSendChanSize, ClientConnSendBufferSize,
		ClientConnRecvChanSize, ClientConnRecvBufferSize)
	conn.CreateTime = int64(time.Now().Unix())
	conn.readch = conn.IConnection.GetRecvMessageChannel()
	conn.onRecv = onRecv
	conn.onClose = onClose
	go conn.recvMsgThread()
	return conn
}

func ClientDial(addr string,
	onRecv func(*Client, *msg.MessageBinary),
	onClose func(*Client)) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	client := NewClient(conn, onRecv, onClose)
	client.StartReadData()
	return client, err
}

func (this *Client) StartReadData() {
	this.IConnection.StartRecv()
}

func (this *Client) recvMsgThread() {
	defer func() {
		if this.onClose != nil {
			this.onClose(this)
		}
	}()

	for {
		select {
		case m, ok := <-this.readch:
			if !ok || m == nil {
				return
			}
			if this.onRecv != nil {
				this.onRecv(this, m)
			}
		}
	}
}

// 返回连接是否仍可用
func (this *Client) Check() bool {
	curtime := int64(time.Now().Unix())
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
		if !this.IsVertify() {
			this.Debug("[Client.Check] 长时间未通过验证，断开连接")
		} else {
			this.Debug("[Client.Check] 长时间未活动，断开连接")
		}
		return false
	}
	return true
}

func (this *Client) GetPing() *Ping {
	return &this.ping
}

// 设置过期时间
func (this *Client) SetTerminateTime(value int64) {
	this.terminate_time = value
}

// 判断是否已强制终止
func (this *Client) IsTerminateTimeout(curtime int64) bool {
	if this.terminate_time > 0 && this.terminate_time < curtime {
		return true
	}
	return false
}

// 判断是否已到达终止时间
func (this *Client) IsTerminateForce() bool {
	return this.terminate_force
}

// 判断是否已终止
func (this *Client) IsTerminate(curtime int64) bool {
	if this.IsTerminateForce() || this.IsTerminateTimeout(curtime) {
		return true
	}
	return false
}

// 强制终止该连接
func (this *Client) Terminate() {
	this.terminate_force = true
}

// 异步发送一条消息，不带发送完成回调
func (this *Client) SendCmd(v msg.MsgStruct) error {
	if !this.IConnection.IsAlive() {
		this.Warn("[Client.SendCmd] 连接已被关闭，取消发送 Msg[%s]",
			v.GetMsgName())
		return fmt.Errorf("link has been closed")
	}
	this.Debug("[SendCmd] 发送 MsgID[%d] MsgName[%s] DataLen[%d]",
		v.GetMsgId(), v.GetMsgName(), v.GetSize())

	msg := msg.MakeMessageByObj(v)
	if msg == nil {
		this.Error("[Client.SendCmd] msg==nil")
		return fmt.Errorf("can't get message binary")
	}
	return this.IConnection.SendMessageBinary(msg)
}

// 异步发送一条消息，带发送完成回调
func (this *Client) SendCmdWithCallback(v msg.MsgStruct,
	cb func(interface{}), cbarg interface{}) error {
	if !this.IConnection.IsAlive() {
		this.Warn("[Client.SendCmdWithCallback] 连接已被关闭，取消发送 Msg[%s]",
			v.GetMsgName())
		return fmt.Errorf("link has been closed")
	}
	this.Debug("[SendCmdWithCallback] 发送 MsgID[%d] MsgName[%s] DataLen[%d]",
		v.GetMsgId(), v.GetMsgName(), v.GetSize())
	msg := msg.MakeMessageByObj(v)
	if msg == nil {
		this.Error("[Client.SendCmd] msg==nil")
		return fmt.Errorf("can't get message binary")
	}
	msg.RegSendDone(cb, cbarg)
	return this.IConnection.SendMessageBinary(msg)
}

func (this *Client) SendBytes(
	cmdid uint16, protodata []byte) error {
	return this.IConnection.SendBytes(cmdid, protodata)
}

func (this *Client) SetLogger(logger *log.Logger) {
	this.Logger = logger.Clone()
	this.Logger.SetTopic(fmt.Sprintf("Client.CID(%s).IP(%s)",
		this.GetConnectID(), this.IConnection.RemoteAddr()))
}
