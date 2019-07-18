package subnet

import (
	"github.com/liasece/micserver/tcpconn"
)

type SubnetHandler struct {
}

func (this *SubnetHandler) OnCreateTCPConnect(conn *tcpconn.ServerConn) {
}

func (this *SubnetHandler) OnRemoveTCPConnect(conn *tcpconn.ServerConn) {
}
