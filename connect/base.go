package connect

import (
	"fmt"
	"time"

	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/msg"
)

type BaseConnect struct {
	*log.Logger
	// 连接实体
	IConnection

	// 唯一编号
	Tempid string
	// 结束时间 为0表示不结束
	terminate_time int64
	// 主动断开连接标记
	terminate_force bool
	// 验证是否成功，没有成功不允许处理后面的消息
	verify_ok bool
	// 当前连接上的计算量
	jobnum uint32
	// 是否是正常的断开连接
	IsNormalDisconnect bool
	// 连接创建的时间
	CreateTime int64
	// 连接的延迟信息
	ping Ping
}

func (this *BaseConnect) Init() {
	this.CreateTime = int64(time.Now().Unix())
}

func (this *BaseConnect) SetLogger(l *log.Logger) {
	this.Logger = l.Clone()
}

// 设置过期时间
func (this *BaseConnect) SetTerminateTime(value int64) {
	this.terminate_time = value
}

// 判断是否已强制终止
func (this *BaseConnect) IsTerminateTimeout(curtime int64) bool {
	if this.terminate_time > 0 && this.terminate_time < curtime {
		return true
	}
	return false
}

// 判断是否已到达终止时间
func (this *BaseConnect) IsTerminateForce() bool {
	return this.terminate_force
}

// 判断是否已终止
func (this *BaseConnect) IsTerminate(curtime int64) bool {
	if this.IsTerminateForce() || this.IsTerminateTimeout(curtime) {
		return true
	}
	return false
}

// 强制终止该连接
func (this *BaseConnect) Terminate() {
	this.Debug("[BaseConnect.Terminate] 连接停止 Tempid[%s]", this.Tempid)
	this.terminate_force = true
}

// 是否通过了验证
func (this *BaseConnect) IsVertify() bool {
	return this.verify_ok
}

// 设置验证状态
func (this *BaseConnect) SetVertify(value bool) {
	this.verify_ok = value
}

// 获取服务器连接当前负载
func (this *BaseConnect) GetJobNum() uint32 {
	return this.jobnum
}

// 设置该服务器连接
func (this *BaseConnect) SetJobNum(jnum uint32) {
	this.jobnum = jnum
}

// 异步发送一条消息，不带发送完成回调
func (this *BaseConnect) SendCmd(v msg.MsgStruct) error {
	if !this.IConnection.IsAlive() {
		this.Warn("[BaseConnect.SendCmd] 连接已被关闭，取消发送 Msg[%s]",
			v.GetMsgName())
		return fmt.Errorf("link has been closed")
	}
	msg := msg.GetByObj(v)
	if msg == nil {
		this.Error("[BaseConnect.SendCmd] msg==nil")
		return fmt.Errorf("can't get message binary")
	}
	return this.IConnection.SendMessageBinary(msg)
}

// 异步发送一条消息，带发送完成回调
func (this *BaseConnect) SendCmdWithCallback(v msg.MsgStruct,
	cb func(interface{}), cbarg interface{}) error {
	if !this.IConnection.IsAlive() {
		this.Warn("[BaseConnect.SendCmdWithCallback] 连接已被关闭，取消发送 Msg[%s]",
			v.GetMsgName())
		return fmt.Errorf("link has been closed")
	}
	msg := msg.GetByObj(v)
	if msg == nil {
		this.Error("[BaseConnect.SendCmd] msg==nil")
		return fmt.Errorf("can't get message binary")
	}
	msg.RegSendFinish(cb, cbarg)
	return this.IConnection.SendMessageBinary(msg)
}

func (this *BaseConnect) SendBytes(
	cmdid uint16, protodata []byte) error {
	return this.IConnection.SendBytes(cmdid, protodata)
}

func (this *BaseConnect) Shutdown() {
	this.IConnection.Shutdown()
}

func (this *BaseConnect) RemoteAddr() string {
	return this.IConnection.RemoteAddr()
}

func (this *BaseConnect) GetIConnection() IConnection {
	return this.IConnection
}

func (this *Client) GetPing() *Ping {
	return &this.ping
}
