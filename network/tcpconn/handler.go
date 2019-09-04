package tcpconn

import (
	"io"
)

type handler struct {
	fdoSendBytes func(writer io.ReadWriter, state interface{},
		data []byte) (int, interface{}, error)
	fdoReadBytes func(reader io.ReadWriter, state interface{},
		toData []byte) (int, interface{}, error)
}

func (this *handler) RegDoSendBytes(
	cb func(writer io.ReadWriter, state interface{},
		data []byte) (int, interface{}, error)) {
	this.fdoSendBytes = cb
}

func (this *handler) RegDoReadBytes(
	cb func(reader io.ReadWriter, state interface{},
		toData []byte) (int, interface{}, error)) {
	this.fdoReadBytes = cb
}
