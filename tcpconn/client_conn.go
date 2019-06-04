package tcpconn

import (
	//	"os"
	"base"
	"base/logger"
	// "math/rand"
	"net"
	// "servercomm"
	// "sync"
	"fmt"
	"io/ioutil"
	"time"
)

type ClientConn struct {
	Conn TCPConn

	Openid string
	UUID   uint64
	Quizid uint64

	Roomid uint64

	Userserverid  uint32 // userserver 的serverid
	Roomserverid  uint32 // RoomServer 的serverid
	Matchserverid uint32 // MatchServer 的serverid

	Tempid          uint64 // 唯一编号
	terminate_time  uint64 // 结束时间 为0表示不结束
	terminate_force bool   // 主动断开连接
	verify_ok       bool   // 验证是否成功，没有成功不允许处理后面的消息

	CreateTime uint64 // 连接创建的时间

	Encryption base.TEncryptionType

	ping Ping

	loghead string
}

const MaxMsgSize = 64 * 1024
const JoinSendMsgSize = 32 * 1024

// 客户端连接发送消息缓冲不宜过大， 10*64KiB*100000连接=64GiB
const ClientConnSendMsgBufferSize = 256

// 发送等待缓冲区大小
const ClientConnMaxWaitSendMsgBufferSize = 2 * 1024 * 1024

// 发送缓冲大小，用于将多个小消息拼接发送的缓冲大小
const ClientConnSendBufferSize = MaxMsgSize + JoinSendMsgSize

// 获取一个新的服务器连接
// sctype: 连接的 客户端/服务器 类型
// conn: 连接的net.Conn对象
func NewClientConn(conn net.Conn) *ClientConn {
	tcpconn := new(ClientConn)
	tcpconn.Conn.Init(conn, ClientConnSendMsgBufferSize,
		ClientConnSendBufferSize, ClientConnMaxWaitSendMsgBufferSize)
	tcpconn.CreateTime = uint64(time.Now().Unix())
	return tcpconn
}

// 返回连接是否仍可用
func (this *ClientConn) Check() bool {
	curtime := uint64(time.Now().Unix())
	// 检查本服务器时候还存活
	if this.IsTerminateForce() {
		// 本服务器关闭
		this.Debug("[ClientConn.Check] 服务器强制断开连接")
		// 强制移除客户端连接
		return false
	}
	// 检查客户端连接是否验证超时
	if this.IsTerminateTimeout(curtime) {
		// 客户端超时未通过验证
		if !this.IsVertify() {
			this.Debug("[ClientConn.Check] 长时间未通过验证，断开连接")
		} else {
			this.Debug("[ClientConn.Check] 长时间未活动，断开连接")
		}
		return false
	}
	return true
}

func (this *ClientConn) GetPing() *Ping {
	return &this.ping
}

// 读数据
func (this *ClientConn) Read() (msg []byte, cmdlen int, err error) {
	msg, err4 := ioutil.ReadAll(this.Conn.GetConn())
	this.Debug("[ClientConn.Read] Read N[%d] ", len(msg))

	return msg, len(msg), err4
}

// 是否通过了验证
func (this *ClientConn) IsVertify() bool {
	return this.verify_ok
}

// 设置验证状态
func (this *ClientConn) SetVertify(value bool) {
	this.verify_ok = value
}

// 设置过期时间
func (this *ClientConn) SetTerminateTime(value uint64) {
	this.terminate_time = value
}

// 判断是否已强制终止
func (this *ClientConn) IsTerminateTimeout(curtime uint64) bool {
	if this.terminate_time > 0 && this.terminate_time < curtime {
		return true
	}
	return false
}

// 判断是否已到达终止时间
func (this *ClientConn) IsTerminateForce() bool {
	return this.terminate_force
}

// 判断是否已终止
func (this *ClientConn) IsTerminate(curtime uint64) bool {
	if this.IsTerminateForce() || this.IsTerminateTimeout(curtime) {
		return true
	}
	return false
}

// 获取本服务器连接的net.Conn对象
func (this *ClientConn) GetConn() net.Conn {
	return this.Conn.GetConn()
}

// 强制终止该连接
func (this *ClientConn) Terminate() {
	this.terminate_force = true
}

// 异步发送一条消息，不带发送完成回调
func (this *ClientConn) SendCmd(v base.MsgStruct) {
	this.Debug("[SendCmd] 发送 MsgID[%d] MsgName[%s] DataLen[%d]",
		v.GetMsgId(), v.GetMsgName(), v.GetSize())
	this.Conn.SendCmd(v, this.Encryption)
}

// 异步发送一条消息，带发送完成回调
func (this *ClientConn) SendCmdWithCallback(v base.MsgStruct,
	callback func(interface{}), cbarg interface{}) {
	this.Debug("[SendCmdWithCallback] 发送 MsgID[%d] MsgName[%s] DataLen[%d]",
		v.GetMsgId(), v.GetMsgName(), v.GetSize())
	this.Conn.SendCmdWithCallback(v, callback, cbarg, this.Encryption)
}

func (this *ClientConn) SendBytes(
	cmdid uint16, protodata []byte) error {
	return this.Conn.SendBytes(cmdid, protodata, this.Encryption)
}

func (this *ClientConn) GetLogHead() string {
	// if this.loghead == "" {
	this.loghead = fmt.Sprintf("[ClientConn] TmpID[%d] IPPort[%s] "+
		"OID[%s] UID[%d] USID[%d] RSID[%d] ",
		this.Tempid, this.Conn.Conn.RemoteAddr().String(), this.Openid,
		this.UUID, this.Userserverid,
		this.Roomserverid)
	// }
	return this.loghead
}

func (this *ClientConn) Debug(fmt string, args ...interface{}) {
	fmt = this.GetLogHead() + fmt
	logger.Debug(fmt, args...)
}

func (this *ClientConn) Warn(fmt string, args ...interface{}) {
	fmt = this.GetLogHead() + fmt
	logger.Warn(fmt, args...)
}

func (this *ClientConn) Info(fmt string, args ...interface{}) {
	fmt = this.GetLogHead() + fmt
	logger.Info(fmt, args...)
}

func (this *ClientConn) Error(fmt string, args ...interface{}) {
	fmt = this.GetLogHead() + fmt
	logger.Error(fmt, args...)
}

func (this *ClientConn) Fatal(fmt string, args ...interface{}) {
	fmt = this.GetLogHead() + fmt
	logger.Fatal(fmt, args...)
}
