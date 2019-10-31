package base

import (
	"github.com/liasece/micserver/servercomm"
)

type ServerHook interface {
	OnForwardToServer(msg *servercomm.SForwardToServer)
	OnForwardFromGate(msg *servercomm.SForwardFromGate)
}
