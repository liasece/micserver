package subnet

import (
	"github.com/liasece/micserver/connect"
	"github.com/liasece/micserver/msg"
)

type SubnetCallback struct {
	fonRecvMsg func(conn *connect.Server, msgbinary *msg.MessageBinary)
}

func (this *SubnetCallback) RegOnRecvMsg(
	cb func(conn *connect.Server, msgbinary *msg.MessageBinary)) {
	this.fonRecvMsg = cb
}
