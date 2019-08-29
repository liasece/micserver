package tcpconn

import (
	"io"
)

type handler struct {
	regSendTCPBytes func(writer io.ReadWriter, state interface{},
		data []byte) (int, interface{}, error)
	regReadTCPBytes func(reader io.ReadWriter, state interface{},
		toData []byte) (int, interface{}, error)
}

func (this *handler) RegSendTCPBytes(
	cb func(writer io.ReadWriter, state interface{},
		data []byte) (int, interface{}, error)) {
	this.regSendTCPBytes = cb
}

func (this *handler) RegReadTCPBytes(
	cb func(reader io.ReadWriter, state interface{},
		toData []byte) (int, interface{}, error)) {
	this.regReadTCPBytes = cb
}
