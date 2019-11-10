package process

import (
	"sync"

	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/servercomm"
)

type ChanServerHandshake struct {
	ModuleInfo    *servercomm.ModuleInfo
	ClientMsgChan chan *msg.MessageBinary
	ServerMsgChan chan *msg.MessageBinary
	Seq           int
}

var (
	_gServerChan sync.Map
)

func AddServerChan(id string, serverChan chan *ChanServerHandshake) {
	_gServerChan.Store(id, serverChan)
}

func DeleteServerChan(id string) {
	_gServerChan.Delete(id)
}

func GetServerChan(id string) chan *ChanServerHandshake {
	if vi, ok := _gServerChan.Load(id); ok {
		if v, ok := vi.(chan *ChanServerHandshake); ok {
			return v
		}
	}
	return nil
}
