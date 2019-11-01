package connect

import (
	"net"

	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/network/chanconn"
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

func NewChan(sendChan chan *msg.MessageBinary, recvChan chan *msg.MessageBinary,
	l *log.Logger) IConnection {
	chanconn := &chanconn.ChanConn{}
	chanconn.SetLogger(l)
	chanconn.Init(sendChan, recvChan)
	return chanconn
}
