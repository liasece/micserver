// Package baseio micserver 中网络连接的基础IO接口实现，支持自定义网络协议
package baseio

import (
	"io"
)

// Protocal 网络层协议接口
type Protocal interface {
	// 读数据
	DoRead(reader io.ReadWriter, state interface{},
		data []byte) (int, interface{}, error)
	// 写数据
	DoWrite(writer io.ReadWriter, state interface{},
		toData []byte) (int, interface{}, error)
}

// Worker 协议处理工作
type Worker struct {
	rw       io.ReadWriter
	protocal Protocal
	// 外部协议提供的发送接收中间状态，在使用 TCP 之上的协议时，
	// 可以通过该成员提供协议中间状态
	protocolState interface{}
}

// Init 初始化网络基础IO
func (w *Worker) Init(rw io.ReadWriter) {
	w.rw = rw
}

// HookProtocal 设置网络层协议工作
func (w *Worker) HookProtocal(p Protocal) {
	w.protocal = p
}

// Read 从io中读数据，如果自定义了网络层协议，则使用自定义协议读取，否则默认TCP协议
func (w *Worker) Read(toData []byte) (int, error) {
	if w.protocal != nil {
		n, state, err := w.protocal.DoRead(
			w.rw, w.protocolState, toData)
		w.protocolState = state
		return n, err
	}
	return w.rw.Read(toData)
}

// Write 向io中写数据，如果自定义了网络层协议，则使用自定义协议写入，否则默认TCP协议
func (w *Worker) Write(data []byte) (int, error) {
	if w.protocal != nil {
		n, state, err := w.protocal.DoWrite(w.rw,
			w.protocolState, data)
		w.protocolState = state
		return n, err
	}

	return w.rw.Write(data)
}
