package tcpconn

import (
	"github.com/liasece/micserver/msg"
	// "math/rand"
	"net"
	// "servercomm"
	// "sync"
	"fmt"
	"time"
)

type ClientConn struct {
	TCPConn
	// 唯一编号
	Tempid string
	// 结束时间 为0表示不结束
	terminate_time int64
	// 主动断开连接
	terminate_force bool
	// 验证是否成功，没有成功不允许处理后面的消息
	verify_ok bool
	// 连接创建的时间
	CreateTime int64

	Encryption msg.TEncryptionType

	ping Ping

	loghead string
}

// 客户端连接发送消息缓冲不宜过大， 10*64KiB*100000连接=64GiB
const ClientConnSendChanSize = 256

// 发送缓冲大小，用于将多个小消息拼接发送的缓冲大小
const ClientConnSendBufferSize = msg.MessageMaxSize * 2

// 获取一个新的服务器连接
// sctype: 连接的 客户端/服务器 类型
// conn: 连接的net.Conn对象
func NewClientConn(conn net.Conn) *ClientConn {
	tcpconn := new(ClientConn)
	tcpconn.Init(conn, ClientConnSendChanSize, ClientConnSendBufferSize)
	tcpconn.CreateTime = int64(time.Now().Unix())
	return tcpconn
}

func ClientDial(addr string) (*ClientConn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return NewClientConn(conn), err
}

// 返回连接是否仍可用
func (this *ClientConn) Check() bool {
	curtime := int64(time.Now().Unix())
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
	msg, err4 := this.ReadAll()
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
func (this *ClientConn) SetTerminateTime(value int64) {
	this.terminate_time = value
}

// 判断是否已强制终止
func (this *ClientConn) IsTerminateTimeout(curtime int64) bool {
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
func (this *ClientConn) IsTerminate(curtime int64) bool {
	if this.IsTerminateForce() || this.IsTerminateTimeout(curtime) {
		return true
	}
	return false
}

// 强制终止该连接
func (this *ClientConn) Terminate() {
	this.terminate_force = true
}

// 异步发送一条消息，不带发送完成回调
func (this *ClientConn) SendCmd(v msg.MsgStruct) {
	this.Debug("[SendCmd] 发送 MsgID[%d] MsgName[%s] DataLen[%d]",
		v.GetMsgId(), v.GetMsgName(), v.GetSize())
	this.TCPConn.SendCmd(v, this.Encryption)
}

// 异步发送一条消息，带发送完成回调
func (this *ClientConn) SendCmdWithCallback(v msg.MsgStruct,
	callback func(interface{}), cbarg interface{}) {
	this.Debug("[SendCmdWithCallback] 发送 MsgID[%d] MsgName[%s] DataLen[%d]",
		v.GetMsgId(), v.GetMsgName(), v.GetSize())
	this.TCPConn.SendCmdWithCallback(v, callback, cbarg, this.Encryption)
}

func (this *ClientConn) SendBytes(
	cmdid uint16, protodata []byte) error {
	return this.TCPConn.SendBytes(cmdid, protodata, this.Encryption)
}

func (this *ClientConn) GetLogHead() string {
	this.loghead = fmt.Sprintf("[ClientConn] TmpID[%s] IPPort[%s] ",
		this.Tempid, this.Conn.RemoteAddr().String())
	return this.loghead
}

func (this *ClientConn) Debug(fmt string, args ...interface{}) {
	fmt = this.GetLogHead() + fmt
	this.Logger.Debug(fmt, args...)
}

func (this *ClientConn) Warn(fmt string, args ...interface{}) {
	fmt = this.GetLogHead() + fmt
	this.Logger.Warn(fmt, args...)
}

func (this *ClientConn) Info(fmt string, args ...interface{}) {
	fmt = this.GetLogHead() + fmt
	this.Logger.Info(fmt, args...)
}

func (this *ClientConn) Error(fmt string, args ...interface{}) {
	fmt = this.GetLogHead() + fmt
	this.Logger.Error(fmt, args...)
}

func (this *ClientConn) Fatal(fmt string, args ...interface{}) {
	fmt = this.GetLogHead() + fmt
	this.Logger.Fatal(fmt, args...)
}
