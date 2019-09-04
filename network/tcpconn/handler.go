package tcpconn

import (
	"io"
)

type handler struct {
	fdoSendTCPBytes func(writer io.ReadWriter, state interface{},
		data []byte) (int, interface{}, error)
	fdoReadTCPBytes func(reader io.ReadWriter, state interface{},
		toData []byte) (int, interface{}, error)
}

func (this *handler) RegDoSendTCPBytes(
	cb func(writer io.ReadWriter, state interface{},
		data []byte) (int, interface{}, error)) {
	this.fdoSendTCPBytes = cb
}

func (this *handler) RegDoReadTCPBytes(
	cb func(reader io.ReadWriter, state interface{},
		toData []byte) (int, interface{}, error)) {
	this.fdoReadTCPBytes = cb
}
