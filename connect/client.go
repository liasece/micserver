package connect

import (
	"fmt"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/session"
	"github.com/liasece/micserver/tcpconn"
	"net"
	"time"
)

type Client struct {
	// 会话信息 可在不同服务器之间同步的
	session.Session
	// 连接实体
	tcpconn.TCPConn
	// 结束时间 为0表示不结束
	terminate_time int64
	// 主动断开连接
	terminate_force bool
	// 连接创建的时间
	CreateTime int64

	// 连接的延迟信息
	ping Ping
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
	ch := conn.Init(netconn,
		ClientConnSendChanSize, ClientConnSendBufferSize,
		ClientConnRecvChanSize, ClientConnRecvBufferSize)
	conn.CreateTime = int64(time.Now().Unix())
	conn.Session = make(map[string]string)
	go conn.recvMsgThread(ch, onRecv, onClose)
	return conn
}

func ClientDial(addr string,
	onRecv func(*Client, *msg.MessageBinary),
	onClose func(*Client)) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return NewClient(conn, onRecv, onClose), err
}

func (this *Client) recvMsgThread(c chan *msg.MessageBinary,
	onRecv func(*Client, *msg.MessageBinary),
	onClose func(*Client)) {
	defer func() {
		if onClose != nil {
			onClose(this)
		}
	}()

	for {
		select {
		case m, ok := <-c:
			if !ok || m == nil {
				return
			}
			if onRecv != nil {
				onRecv(this, m)
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

// 读数据
func (this *Client) Read() (msg []byte, cmdlen int, err error) {
	msg, err4 := this.ReadAll()
	this.Debug("[Client.Read] Read N[%d] ", len(msg))

	return msg, len(msg), err4
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
func (this *Client) SendCmd(v msg.MsgStruct) {
	this.Debug("[SendCmd] 发送 MsgID[%d] MsgName[%s] DataLen[%d]",
		v.GetMsgId(), v.GetMsgName(), v.GetSize())
	this.TCPConn.SendCmd(v)
}

// 异步发送一条消息，带发送完成回调
func (this *Client) SendCmdWithCallback(v msg.MsgStruct,
	callback func(interface{}), cbarg interface{}) {
	this.Debug("[SendCmdWithCallback] 发送 MsgID[%d] MsgName[%s] DataLen[%d]",
		v.GetMsgId(), v.GetMsgName(), v.GetSize())
	this.TCPConn.SendCmdWithCallback(v, callback, cbarg)
}

func (this *Client) SendBytes(
	cmdid uint16, protodata []byte) error {
	return this.TCPConn.SendBytes(cmdid, protodata)
}

func (this *Client) SetLogger(logger *log.Logger) {
	this.Logger = logger.Clone()
	this.Logger.SetTopic(fmt.Sprintf("Client.CID(%s).IP(%s)",
		this.GetConnectID(), this.Conn.RemoteAddr().String()))
}