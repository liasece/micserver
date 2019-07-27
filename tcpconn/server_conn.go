package tcpconn

import (
	//	"os"
	// "msg/log"
	"github.com/liasece/micserver/comm"
	"github.com/liasece/micserver/msg"
	"math/rand"
	"net"
	// "sync"
	// "time"
)

type TServerSCType uint32

const (
	ServerSCTypeNone   TServerSCType = 1
	ServerSCTypeTask   TServerSCType = 2
	ServerSCTypeClient TServerSCType = 3
)

// 服务器连接发送消息缓冲要考虑到服务器处理消息的能力
const ServerConnSendMsgBufferSize = 500000

// 发送等待缓冲区大小
const ServerConnMaxWaitSendMsgBufferSize = 512 * 1024 * 1024

// 发送缓冲大小，用于将多个小消息拼接发送的缓冲大小
const ServerConnSendBufferSize = 64 * 1024 * 1024

type ServerConn struct {
	TCPConn
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
	Serverinfo comm.SServerInfo
	// 用于区分该连接是服务器 client task 连接
	serverSCType TServerSCType
}

// 获取一个新的服务器连接
// sctype: 连接的 客户端/服务器 类型
// conn: 连接的net.Conn对象
func NewServerConn(sctype TServerSCType, conn net.Conn) *ServerConn {
	tcpconn := new(ServerConn)
	tcpconn.SetSC(sctype)
	tcpconn.Init(conn, ServerConnSendMsgBufferSize,
		ServerConnSendBufferSize, ServerConnMaxWaitSendMsgBufferSize)
	tcpconn.ConnectPriority = rand.Int63()
	return tcpconn
}

// 获取服务器连接当前负载
func (this *ServerConn) GetJobNum() uint32 {
	return this.jobnum
}

// 设置该服务器连接
func (this *ServerConn) SetJobNum(jnum uint32) {
	this.jobnum = jnum
}

// 是否通过了验证
func (this *ServerConn) IsVertify() bool {
	return this.verify_ok
}

// 设置验证状态
func (this *ServerConn) SetVertify(value bool) {
	this.verify_ok = value
}

// 设置过期时间
func (this *ServerConn) SetTerminateTime(value uint64) {
	this.terminate_time = value
}

// 判断是否已强制终止
func (this *ServerConn) IsTerminateTimeout(curtime uint64) bool {
	if this.terminate_time > 0 && this.terminate_time < curtime {
		return true
	}
	return false
}

// 判断是否已到达终止时间
func (this *ServerConn) IsTerminateForce() bool {
	return this.terminate_force
}

// 判断是否已终止
func (this *ServerConn) IsTerminate(curtime uint64) bool {
	if this.IsTerminateForce() || this.IsTerminateTimeout(curtime) {
		return true
	}
	return false
}

// 强制终止该连接
func (this *ServerConn) Terminate() {
	this.Debug("[ServerConn.Terminate] 连接停止 Tempid[%s]", this.Tempid)
	this.terminate_force = true
}

// 异步发送一条消息，不带发送完成回调
func (this *ServerConn) SendCmd(v msg.MsgStruct) error {
	return this.TCPConn.SendCmd(v, 0)
}

// 异步发送一条消息，带发送完成回调
func (this *ServerConn) SendCmdWithCallback(v msg.MsgStruct,
	callback func(interface{}), cbarg interface{}) error {
	return this.TCPConn.SendCmdWithCallback(v, callback, cbarg, 0)
}

func (this *ServerConn) SetSC(sctype TServerSCType) {
	this.serverSCType = sctype
}

func (this *ServerConn) GetSCType() TServerSCType {
	return this.serverSCType
}
