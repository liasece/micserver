package session

import (
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/servercomm"
	"sync"
)

type ISInner_SendServerMsg interface {
	SInner_SendServerMsg(gate string, msg msg.MsgStruct)
	GetModuleType() string
}

type Session struct {
	m sync.Map
}

func (this *Session) get(key string) string {
	if vi, ok := this.m.Load(key); ok {
		if vi == nil {
			return ""
		}
		return vi.(string)
	}
	return ""
}

func (this *Session) set(key string, value string) {
	this.m.Store(key, value)
}

func (this *Session) GetUUID() string {
	return this.get("UUID")
}

func (this *Session) SetUUID(value string) {
	this.set("UUID", value)
}

func (this *Session) GetConnectID() string {
	return this.get("ConnectID")
}

func (this *Session) SetConnectID(value string) {
	this.set("ConnectID", value)
}

func (this *Session) GetBind(moduleType string) string {
	return this.get("bind_" + moduleType)
}

func (this *Session) SetBind(moduleType string, value string) {
	this.set("bind_"+moduleType, value)
}

func (this *Session) HasKey(key string) bool {
	_, ok := this.m.Load(key)
	return ok
}

func (this *Session) IsVertify() bool {
	if !this.HasKey("UUID") {
		return false
	}
	if this.GetUUID() == "" {
		return false
	}
	return true
}

func (this *Session) SyncToServer(mod ISInner_SendServerMsg,
	targetServer string) {
	smsg := &servercomm.SUpdateSession{
		Session:      this.ToMap(),
		ClientConnID: this.GetConnectID(),
	}
	mod.SInner_SendServerMsg(targetServer, smsg)
}

func (this *Session) ToMap() map[string]string {
	res := make(map[string]string)
	this.m.Range(func(ki, vi interface{}) bool {
		if ki != nil {
			if vi != nil {
				res[ki.(string)] = vi.(string)
			} else {
				res[ki.(string)] = ""
			}
		}
		return true
	})
	return res
}

func (this *Session) FromMap(m map[string]string) {
	for k, v := range m {
		this.set(k, v)
	}
}
