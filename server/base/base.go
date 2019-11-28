package base

import (
	"github.com/liasece/micserver/servercomm"
	"github.com/liasece/micserver/session"
)

type ServerHook interface {
	OnModuleMessage(msg *servercomm.ModuleMessage)
	OnClientMessage(se *session.Session, msg *servercomm.ClientMessage)
}
