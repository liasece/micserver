/*
connect 实现了 micserver 中 模块间连接/客户端Gateway连接 的逻辑，
包括了所有连接需要用到的方法，连接池管理方法。
*/
package connect

import (
	"fmt"
	"time"

	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/network/baseio"
	"github.com/liasece/micserver/util/uid"
)

// 一个连接都会具有的基础对象，整合了一个连接的基础属性
type BaseConnect struct {
	*log.Logger
	// 连接实体
	IConnection IConnection

	// 唯一编号
	tempID string
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
	createTime int64
	// 连接的延迟信息
	ping Ping
}

// 初始化这个基础连接
func (this *BaseConnect) Init() {
	this.createTime = int64(time.Now().Unix())
	tmpid, err := uid.GenUniqueID(uint16(time.Now().UnixNano()))
	if err == nil {
		this.setTempID(tmpid)
	} else {
		this.Error("[BaseConnect.Init] 生成UUID出错 Error[%s]",
			err.Error())
	}
}

// 设置该连接的 Logger ，便于Log信息整理
func (this *BaseConnect) SetLogger(l *log.Logger) {
	this.Logger = l.Clone()
}

// 获取该连接对象构造完成的时间，由 BaseConnect.Init 初始化
func (this *BaseConnect) GetCreateTime() int64 {
	return this.createTime
}

// 获取连接的唯一ID
func (this *BaseConnect) GetTempID() string {
	return this.tempID
}

// 设置连接的唯一ID，不可提供给外部更改，因为其更改需要保证连接池等周边系统一并更改
func (this *BaseConnect) setTempID(id string) {
	this.tempID = id
}

// 设置连接过期时间，如果一个连接过期了，在下一个发送或接收行为会将连接置为断开状态
func (this *BaseConnect) SetTerminateTime(value int64) {
	this.terminate_time = value
}

// 通过过期时间等判断是否已强制终止
func (this *BaseConnect) IsTerminateTimeout(curtime int64) bool {
	if this.terminate_time > 0 && this.terminate_time < curtime {
		return true
	}
	return false
}

// 强制终止该连接
func (this *BaseConnect) Terminate() {
	this.Debug("[BaseConnect.Terminate] 连接停止 tempID[%s]", this.tempID)
	this.terminate_force = true
}

// 判断连接是否已被主动强制终止，通过 BaseConnect.Terminate 设置
func (this *BaseConnect) IsTerminateForce() bool {
	return this.terminate_force
}

// 判断连接是否已终止，包括了主动终止以及超时终止
func (this *BaseConnect) IsTerminate(curtime int64) bool {
	if this.IsTerminateForce() || this.IsTerminateTimeout(curtime) {
		return true
	}
	return false
}

// 设置该连接的消息编解码器
func (this *BaseConnect) SetMsgCodec(codec msg.IMsgCodec) {
	this.IConnection.SetMsgCodec(codec)
}

// 开始接收消息，在调用该消息前，无法在 BaseConnect.GetRecvMessageChannel
// 中接收到消息
func (this *BaseConnect) StartRecv() {
	this.IConnection.StartRecv()
}

// 设置该连接的协议，如果 p == nil ，该连接的网络协议视为普通的 TCP 协议，
// 通过提供非空的 p ，可以将该连接实现为使用 websocket 等其他网络协议
func (this *BaseConnect) HookProtocal(p baseio.Protocal) {
	this.IConnection.HookProtocal(p)
}

// 获取该连接的消息处理 channel ，可以通过该 channel 接收到该连接收到的消息，
// 接收到的消息已经经过了 BaseConnect.HookProtocal 处理特殊网络协议，
// 经过 SetMsgCodec 处理特殊消息编解码格式，
// 从该 channel 中得到的 *msg.MessageBinary 的 protodata 已是消息本身的内容。
func (this *BaseConnect) GetRecvMessageChannel() chan *msg.MessageBinary {
	return this.IConnection.GetRecvMessageChannel()
}

// 该连接是否通过了验证，如果这是一个 Module 间的连接，需要经过 server.Server
// 的登陆逻辑处理，才会成为一个经过验证的连接。如果是一个客户端连接，
// 默认该连接不是一个经过验证的连接。
func (this *BaseConnect) IsVertify() bool {
	return this.verify_ok
}

// 设置该连接的验证状态
func (this *BaseConnect) SetVertify(value bool) {
	this.verify_ok = value
}

// 获取该连接的负载
func (this *BaseConnect) GetJobNum() uint32 {
	return this.jobnum
}

// 设置该连接的负载
func (this *BaseConnect) SetJobNum(jnum uint32) {
	this.jobnum = jnum
}

// 异步发送一条消息
func (this *BaseConnect) SendCmd(v msg.MsgStruct) error {
	if !this.IConnection.IsAlive() {
		this.Warn("[BaseConnect.SendCmd] 连接已被关闭，取消发送 Msg[%s]",
			v.GetMsgName())
		return fmt.Errorf("link has been closed")
	}
	msg := this.IConnection.GetMsgCodec().EncodeObj(v)
	if msg == nil {
		this.Error("[BaseConnect.SendCmd] msg==nil")
		return fmt.Errorf("can't get message binary")
	}
	return this.IConnection.SendMessageBinary(msg)
}

// 异步发送一条消息，带发送完成回调，在消息真正通过 network 发送成功之后，
// 会调用 cb 回调。
func (this *BaseConnect) SendCmdWithCallback(v msg.MsgStruct,
	cb func(interface{}), cbarg interface{}) error {
	if !this.IConnection.IsAlive() {
		this.Warn("[BaseConnect.SendCmdWithCallback] 连接已被关闭，取消发送 Msg[%s]",
			v.GetMsgName())
		return fmt.Errorf("link has been closed")
	}
	msg := this.IConnection.GetMsgCodec().EncodeObj(v)
	if msg == nil {
		this.Error("[BaseConnect.SendCmd] msg==nil")
		return fmt.Errorf("can't get message binary")
	}
	msg.RegSendFinish(cb, cbarg)
	return this.IConnection.SendMessageBinary(msg)
}

// 异步发送一条消息，使用 cmdid 及 protodata 来发送，如果不使用 msg.MsgStruct
// 作为消息发送，你可以利用该方法，将消息编码之后发送。
func (this *BaseConnect) SendBytes(
	cmdid uint16, protodata []byte) error {
	return this.IConnection.SendBytes(cmdid, protodata)
}

// 断开该连接的底层连接
func (this *BaseConnect) Shutdown() {
	this.IConnection.Shutdown()
}

// 该连接的远程地址
func (this *BaseConnect) RemoteAddr() string {
	return this.IConnection.RemoteAddr()
}

// 获取该连接的底层连接接口
func (this *BaseConnect) GetIConnection() IConnection {
	return this.IConnection
}

// 获取该连接的 Ping 信息
func (this *Client) GetPing() *Ping {
	return &this.ping
}
