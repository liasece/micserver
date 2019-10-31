package baseio

import (
	"io"
)

type Protocal interface {
	DoRead(reader io.ReadWriter, state interface{},
		data []byte) (int, interface{}, error)
	DoWrite(writer io.ReadWriter, state interface{},
		toData []byte) (int, interface{}, error)
}

type Worker struct {
	rw       io.ReadWriter
	protocal Protocal
	// 外部协议提供的发送接收中间状态，在使用 TCP 之上的协议时，
	// 可以通过该成员提供协议中间状态
	protocolState interface{}
}

func (this *Worker) Init(rw io.ReadWriter) {
	this.rw = rw
}

func (this *Worker) HookProtocal(p Protocal) {
	this.protocal = p
}

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
