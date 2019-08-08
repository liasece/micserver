package subnet

import (
	"github.com/liasece/micserver/servercomm"
)

type SubnetCallback struct {
	regHandleServerMsg func(msg *servercomm.SForwardToServer)
	regHandleGateMsg   func(msg *servercomm.SForwardFromGate)
}

func (this *SubnetCallback) RegHandleServerMsg(
	cb func(msg *servercomm.SForwardToServer)) {
	this.regHandleServerMsg = cb
}

func (this *SubnetCallback) RegHandleGateMsg(
	cb func(msg *servercomm.SForwardFromGate)) {
	this.regHandleGateMsg = cb
}
