package subnet

import (
	"github.com/liasece/micserver/servercomm"
)

type SubnetCallback struct {
	regHandleServerMsg func(msg *servercomm.SForwardToServer)
}

func (this *SubnetCallback) RegHandleServerMsg(
	cb func(msg *servercomm.SForwardToServer)) {
	this.regHandleServerMsg = cb
}
