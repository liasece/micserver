package session

import (
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/servercomm"
)

type ISInner_SendServerMsg interface {
	SInner_SendServerMsg(gate string, msg msg.MsgStruct)
	GetServerType() string
}

type Session map[string]string

func (this *Session) get(key string) string {
	if v, ok := (*this)[key]; ok {
		return v
	}
	return ""
}

func (this *Session) set(key string, value string) {
	(*this)[key] = value
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

func (this *Session) GetBindServer(servertype string) string {
	return this.get("bindserver_" + servertype)
}

func (this *Session) SetBindServer(servertype string, value string) {
	this.set("bindserver_"+servertype, value)
}

func (this *Session) HasKey(key string) bool {
	_, ok := (*this)[key]
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
		Session:      *this,
		ClientConnID: this.GetConnectID(),
	}
	mod.SInner_SendServerMsg(targetServer, smsg)
}
