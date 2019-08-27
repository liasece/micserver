package subnet

import (
	"github.com/liasece/micserver/servercomm"
)

type SubnetCallback struct {
	regForwardToServer func(msg *servercomm.SForwardToServer)
	regForwardFromGate func(msg *servercomm.SForwardFromGate)
	regForwardToClient func(msg *servercomm.SForwardToClient)
	regUpdateSession   func(msg *servercomm.SUpdateSession)
}

func (this *SubnetCallback) RegForwardToServer(
	cb func(msg *servercomm.SForwardToServer)) {
	this.regForwardToServer = cb
}

func (this *SubnetCallback) RegForwardFromGate(
	cb func(msg *servercomm.SForwardFromGate)) {
	this.regForwardFromGate = cb
}

func (this *SubnetCallback) RegForwardToClient(
	cb func(msg *servercomm.SForwardToClient)) {
	this.regForwardToClient = cb
}

func (this *SubnetCallback) RegUpdateSession(
	cb func(msg *servercomm.SUpdateSession)) {
	this.regUpdateSession = cb
}
