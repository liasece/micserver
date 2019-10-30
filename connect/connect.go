package connect

import (
	"net"

	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/network/tcpconn"
)

func NewTCP(conn net.Conn, l *log.Logger, sendChanSize int, sendBufferSize int,
	recvChanSize int, recvBufferSize int) IConnection {
	tcp := &tcpconn.TCPConn{}
	tcp.SetLogger(l)
	tcp.Init(conn,
		sendChanSize, sendBufferSize,
		recvChanSize, recvBufferSize)
	return tcp
}
