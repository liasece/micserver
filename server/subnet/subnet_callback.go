package subnet

import (
	"github.com/liasece/micserver/servercomm"
)

type SubnetCallback struct {
	fonForwardToServer func(msg *servercomm.SForwardToServer)
	fonForwardFromGate func(msg *servercomm.SForwardFromGate)
	fonForwardToClient func(msg *servercomm.SForwardToClient)
	fonUpdateSession   func(msg *servercomm.SUpdateSession)
}

func (this *SubnetCallback) RegOnForwardToServer(
	cb func(msg *servercomm.SForwardToServer)) {
	this.fonForwardToServer = cb
}

func (this *SubnetCallback) RegOnForwardFromGate(
	cb func(msg *servercomm.SForwardFromGate)) {
	this.fonForwardFromGate = cb
}

func (this *SubnetCallback) RegOnForwardToClient(
	cb func(msg *servercomm.SForwardToClient)) {
	this.fonForwardToClient = cb
}

func (this *SubnetCallback) RegOnUpdateSession(
	cb func(msg *servercomm.SUpdateSession)) {
	this.fonUpdateSession = cb
}
