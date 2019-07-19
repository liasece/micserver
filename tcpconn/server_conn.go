package tcpconn

import (
	//	"os"
	// "msg/log"
	"github.com/liasece/micserver/comm"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/msg"
	"math/rand"
	"net"
	// "sync"
	// "time"
)

type TServerSCType uint32

const (
	ServerSCTypeNone TServerSCType = iota
	ServerSCTypeTask
	ServerSCTypeClient
)

// 服务器连接发送消息缓冲要考虑到服务器处理消息的能力
const ServerConnSendMsgBufferSize = 500000

// 发送等待缓冲区大小
const ServerConnMaxWaitSendMsgBufferSize = 512 * 1024 * 1024

// 发送缓冲大小，用于将多个小消息拼接发送的缓冲大小
const ServerConnSendBufferSize = 64 * 1024 * 1024

type ServerConn struct {
	Conn            TCPConn
	Tempid          string // 唯一编号
	terminate_time  uint64 // 结束时间 为0表示不结束
	terminate_force bool   // 主动断开连接
	verify_ok       bool   // 验证是否成功，没有成功不允许处理后面的消息
	isAlive         bool
	jobnum          uint32 // 当前连接上的计算量，userserver表示用户数，
	// matchserver表示正在匹配的队列数，
	// battleserver表示在战斗的房间数
	IsNormalDisconnect bool // 是否是正常的断开连接
	ConnectPriority    int64

	Serverinfo comm.SServerInfo // 该连接对方服务器信息

	serverSCType TServerSCType // 用于区分该连接是服务器 client task 连接
}

// 获取一个新的服务器连接
// sctype: 连接的 客户端/服务器 类型
// conn: 连接的net.Conn对象
func NewServerConn(sctype TServerSCType, conn net.Conn) *ServerConn {
	tcpconn := new(ServerConn)
	tcpconn.SetAlive(true)
	tcpconn.SetSC(sctype)
	tcpconn.Conn.Init(conn, ServerConnSendMsgBufferSize,
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

// 获取本服务器连接的net.Conn对象
func (this *ServerConn) GetConn() net.Conn {
	return this.Conn.GetConn()
}

// 强制终止该连接
func (this *ServerConn) Terminate() {
	log.Debug("[ServerConn.Terminate] 连接停止 Tempid[%s]", this.Tempid)
	this.terminate_force = true
}

// 异步发送一条消息，不带发送完成回调
func (this *ServerConn) SendCmd(v msg.MsgStruct) error {
	return this.Conn.SendCmd(v, 0)
}

// 异步发送一条消息，带发送完成回调
func (this *ServerConn) SendCmdWithCallback(v msg.MsgStruct,
	callback func(interface{}), cbarg interface{}) error {
	return this.Conn.SendCmdWithCallback(v, callback, cbarg, 0)
}

// 异步发送一条消息，带发送完成回调
func (this *ServerConn) IsAlive() bool {
	return this.isAlive
}

// 异步发送一条消息，带发送完成回调
func (this *ServerConn) SetAlive(value bool) {
	this.isAlive = value
}

func (this *ServerConn) SetSC(sctype TServerSCType) {
	this.serverSCType = sctype
}

func (this *ServerConn) GetSCType() TServerSCType {
	return this.serverSCType
}
