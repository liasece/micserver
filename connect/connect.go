package connect

import (
	"github.com/liasece/micserver/network/tcpconn"
	"net"
)

func NewTCP(conn net.Conn, sendChanSize int, sendBufferSize int,
	recvChanSize int, recvBufferSize int) IConnection {
	tcp := &tcpconn.TCPConn{}
	tcp.Init(conn,
		sendChanSize, sendBufferSize,
		recvChanSize, recvBufferSize)
	return tcp
}
