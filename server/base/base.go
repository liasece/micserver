package base

import (
	"github.com/liasece/micserver/servercomm"
)

type ServerHook interface {
	OnModuleMessage(msg *servercomm.ModuleMessage)
	OnClientMessage(msg *servercomm.ClientMessage)
}
