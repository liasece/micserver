package connect

import (
	"fmt"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/servercomm"
	"math/rand"
	"net"
)

type TServerSCType uint32

const (
	ServerSCTypeNone   TServerSCType = 1
	ServerSCTypeTask   TServerSCType = 2
	ServerSCTypeClient TServerSCType = 3
)

// 服务器连接发送消息缓冲要考虑到服务器处理消息的能力
const ServerSendChanSize = 100000

// 发送缓冲大小，用于将多个小消息拼接发送的缓冲大小
const ServerSendBufferSize = msg.MessageMaxSize * 10

// 服务器连接发送消息缓冲要考虑到服务器处理消息的能力
const ServerRecvChanSize = 100000

// 发送缓冲大小，用于将多个小消息拼接发送的缓冲大小
const ServerRecvBufferSize = msg.MessageMaxSize * 10

type Server struct {
	*log.Logger
	// 连接实体
	IConnection
	// 唯一编号
	Tempid string
	// 结束时间 为0表示不结束
	terminate_time uint64
	// 主动断开连接标记
	terminate_force bool
	// 验证是否成功，没有成功不允许处理后面的消息
	verify_ok bool
	// 当前连接上的计算量
	jobnum uint32
	// 是否是正常的断开连接
	IsNormalDisconnect bool
	// 建立连接优先级
	ConnectPriority int64
	// 该连接对方服务器信息
	Serverinfo *servercomm.SServerInfo
	// 用于区分该连接是服务器 client task 连接
	serverSCType TServerSCType
}

// 获取一个新的服务器连接
// sctype: 连接的 客户端/服务器 类型
// netconn: 连接的net.Conn对象
func NewServer(sctype TServerSCType, netconn net.Conn,
	onRecv func(*Server, *msg.MessageBinary),
	onClose func(*Server)) *Server {
	conn := new(Server)
	conn.Serverinfo = &servercomm.SServerInfo{}
	conn.SetSC(sctype)
	conn.ConnectPriority = rand.Int63()
	conn.IConnection = NewTCP(netconn,
		ServerSendChanSize, ServerSendBufferSize,
		ServerRecvChanSize, ServerRecvBufferSize)
	conn.IConnection.StartRecv()
	go conn.recvMsgThread(conn.IConnection.GetRecvMessageChannel(),
		onRecv, onClose)
	return conn
}

func (this *Server) recvMsgThread(c chan *msg.MessageBinary,
	onRecv func(*Server, *msg.MessageBinary),
	onClose func(*Server)) {
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

// 获取服务器连接当前负载
func (this *Server) GetJobNum() uint32 {
	return this.jobnum
}

// 设置该服务器连接
func (this *Server) SetJobNum(jnum uint32) {
	this.jobnum = jnum
}

// 是否通过了验证
func (this *Server) IsVertify() bool {
	return this.verify_ok
}

// 设置验证状态
func (this *Server) SetVertify(value bool) {
	this.verify_ok = value
}

// 设置过期时间
func (this *Server) SetTerminateTime(value uint64) {
	this.terminate_time = value
}

// 判断是否已强制终止
func (this *Server) IsTerminateTimeout(curtime uint64) bool {
	if this.terminate_time > 0 && this.terminate_time < curtime {
		return true
	}
	return false
}

// 判断是否已到达终止时间
func (this *Server) IsTerminateForce() bool {
	return this.terminate_force
}

// 判断是否已终止
func (this *Server) IsTerminate(curtime uint64) bool {
	if this.IsTerminateForce() || this.IsTerminateTimeout(curtime) {
		return true
	}
	return false
}

// 强制终止该连接
func (this *Server) Terminate() {
	this.Debug("[Server.Terminate] 连接停止 Tempid[%s]", this.Tempid)
	this.terminate_force = true
}

// 异步发送一条消息，不带发送完成回调
func (this *Server) SendCmd(v msg.MsgStruct) error {
	if !this.IConnection.IsAlive() {
		this.Warn("[Server.SendCmd] 连接已被关闭，取消发送 Msg[%s]",
			v.GetMsgName())
		return fmt.Errorf("link has been closed")
	}
	msg := msg.GetByObj(v)
	if msg == nil {
		this.Error("[Server.SendCmd] msg==nil")
		return fmt.Errorf("can't get message binary")
	}
	return this.IConnection.SendMessageBinary(msg)
}

// 异步发送一条消息，带发送完成回调
func (this *Server) SendCmdWithCallback(v msg.MsgStruct,
	cb func(interface{}), cbarg interface{}) error {
	if !this.IConnection.IsAlive() {
		this.Warn("[Server.SendCmdWithCallback] 连接已被关闭，取消发送 Msg[%s]",
			v.GetMsgName())
		return fmt.Errorf("link has been closed")
	}
	msg := msg.GetByObj(v)
	if msg == nil {
		this.Error("[Server.SendCmd] msg==nil")
		return fmt.Errorf("can't get message binary")
	}
	msg.RegSendFinish(cb, cbarg)
	return this.IConnection.SendMessageBinary(msg)
}

func (this *Server) SetSC(sctype TServerSCType) {
	this.serverSCType = sctype
}

func (this *Server) GetSCType() TServerSCType {
	return this.serverSCType
}

func (this *Server) Shutdown() {
	this.IConnection.Shutdown()
}

func (this *Server) RemoteAddr() string {
	return this.IConnection.RemoteAddr()
}

func (this *Server) GetIConnection() IConnection {
	return this.IConnection
}
