package tcpconn

type TCPConnError string

func (this TCPConnError) Error() string {
	return string(this)
}

const (
	ErrSendNilData TCPConnError = "send nil data"
	ErrCloseed     TCPConnError = "conn has been closed"
	ErrBufferFull  TCPConnError = "buffer full"
)
