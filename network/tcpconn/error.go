package tcpconn

import (
	"errors"
)

// tcp 连接的错误
var (
	ErrSendNilData = errors.New("send nil data")
	ErrCloseed     = errors.New("conn has been closed")
	ErrBufferFull  = errors.New("buffer full")
)
