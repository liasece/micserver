package base

import (
	"github.com/liasece/micserver/servercomm"
)

type ServerHook interface {
	OnServerMessage(msg *servercomm.ServerMessage)
	OnClientMessage(msg *servercomm.ClientMessage)
}
