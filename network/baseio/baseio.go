/*
micserver 中网络连接的基础IO接口实现，支持自定义网络协议
*/
package baseio

import (
	"io"
)

// 网络层协议接口
type Protocal interface {
	// 读数据
	DoRead(reader io.ReadWriter, state interface{},
		data []byte) (int, interface{}, error)
	// 写数据
	DoWrite(writer io.ReadWriter, state interface{},
		toData []byte) (int, interface{}, error)
}

// 协议处理工作
type Worker struct {
	rw       io.ReadWriter
	protocal Protocal
	// 外部协议提供的发送接收中间状态，在使用 TCP 之上的协议时，
	// 可以通过该成员提供协议中间状态
	protocolState interface{}
}

// 初始化网络基础IO
func (this *Worker) Init(rw io.ReadWriter) {
	this.rw = rw
}

// 设置网络层协议工作
func (this *Worker) HookProtocal(p Protocal) {
	this.protocal = p
}

// 从io中读数据，如果自定义了网络层协议，则使用自定义协议读取，否则默认TCP协议
func (this *Worker) Read(toData []byte) (int, error) {
	if this.protocal != nil {
		n, state, err := this.protocal.DoRead(
			this.rw, this.protocolState, toData)
		this.protocolState = state
		return n, err
	} else {
		return this.rw.Read(toData)
	}
}

// 向io中写数据，如果自定义了网络层协议，则使用自定义协议写入，否则默认TCP协议
func (this *Worker) Write(data []byte) (int, error) {
	if this.protocal != nil {
		n, state, err := this.protocal.DoWrite(this.rw,
			this.protocolState, data)
		this.protocolState = state
		return n, err
	} else {
		return this.rw.Write(data)
	}
}
