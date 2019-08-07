package subnet

import (
	"github.com/liasece/micserver/servercomm"
	"github.com/liasece/micserver/tcpconn"
)

type SubnetCallback struct {
	regHandleServerMsg func(conn *tcpconn.ServerConn, msg *servercomm.SForwardToServer)
}

func (this *SubnetCallback) RegHandleServerMsg(
	cb func(conn *tcpconn.ServerConn, msg *servercomm.SForwardToServer)) {
	this.regHandleServerMsg = cb
}
