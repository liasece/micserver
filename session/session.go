package session

import (
	"fmt"
	"strings"
	"sync"

	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/servercomm"
	"github.com/liasece/micserver/util/conv"
)

type SessionKey string

const (
	SessionKeyBindHead  SessionKey = "_s0_bind_"
	SessionKeyConnectID SessionKey = "_s0_connectid"
	SessionKeyUUID      SessionKey = "_s0_uuid"
)

const (
	GateModuelType = "gate"
)

type ISInner_SendModuleMsg interface {
	GetModuleID() string
	SInner_SendModuleMsg(gate string, msg msg.MsgStruct)
}

type ISInner_SendClientMsg interface {
	SInner_SendClientMsg(gateid string, connectid string, msgid uint16,
		data []byte)
}

type Session struct {
	m sync.Map
}

func NewSessionFromMap(session map[string]string) *Session {
	res := &Session{}
	res.FromMap(session)
	return res
}

func getFromMap(session map[string]string, key SessionKey) string {
	if v, ok := session[string(key)]; ok {
		return v
	}
	return ""
}

func GetUUIDFromMap(session map[string]string) string {
	return getFromMap(session, SessionKeyUUID)
}

func GetBindFromMap(session map[string]string,
	moduleType string) string {
	return getFromMap(session, SessionKeyBindHead+SessionKey(moduleType))
}

func GetConnectIDFromMap(session map[string]string) string {
	return getFromMap(session, SessionKeyConnectID)
}

func (this *Session) get(key SessionKey) string {
	if vi, ok := this.m.Load(string(key)); ok {
		if vi == nil {
			return ""
		}
		return vi.(string)
	}
	return ""
}

func (this *Session) set(key SessionKey, value string) {
	this.m.Store(string(key), value)
}

func (this *Session) rangeBinded(
	f func(moduletype string, moduleid string) bool) {
	this.m.Range(func(ki, vi interface{}) bool {
		if ki == nil || vi == nil {
			return true
		}
		k := ki.(string)
		v := vi.(string)
		if strings.HasPrefix(k, string(SessionKeyBindHead)) {
			slice := strings.Split(k, string(SessionKeyBindHead))
			if len(slice) == 2 {
				// 头部匹配
				if !f(slice[1], v) {
					return false
				}
			}
		}
		return true
	})
}

// 获取该 Session 绑定的所有模块
// 返回值 键为模块类型，值为模块ID
func (this *Session) GetBindedList() map[string]string {
	res := make(map[string]string)
	this.rangeBinded(func(moduletype string, moduleid string) bool {
		res[moduletype] = moduleid
		return true
	})
	return res
}

func (this *Session) GetUUID() string {
	return this.get(SessionKeyUUID)
}

func (this *Session) SetUUID(value string) {
	this.set(SessionKeyUUID, value)
}

// 获取 session 数据的接口
func (this *Session) Get(key SessionKey) string {
	return this.get(key)
}

func (this *Session) GetBool(key SessionKey) bool {
	return conv.MustInterfaceToBool(this.get(key))
}

func (this *Session) GetInt64(key SessionKey) int64 {
	return conv.MustInterfaceToInt64(this.get(key))
}

func (this *Session) Set(key SessionKey, value string) {
	this.set(key, value)
}

func (this *Session) SetBool(key SessionKey, value bool) {
	if value {
		this.set(key, "true")
	} else {
		this.set(key, "false")
	}
}

func (this *Session) SetInt64(key SessionKey, value int64) {
	this.set(key, fmt.Sprint(value))
}

func (this *Session) GetConnectID() string {
	return this.get(SessionKeyConnectID)
}

func (this *Session) SetConnectID(value string) {
	this.set(SessionKeyConnectID, value)
}

func (this *Session) GetBind(moduleType string) string {
	return this.get(SessionKeyBindHead + SessionKey(moduleType))
}

func (this *Session) SetBind(moduleType string, value string) {
	this.set(SessionKeyBindHead+SessionKey(moduleType), value)
}

func (this *Session) HasBind(moduleType string) bool {
	return this.HasKey(SessionKeyBindHead + SessionKey(moduleType))
}

func (this *Session) HasKey(key SessionKey) bool {
	_, ok := this.m.Load(string(key))
	return ok
}

func (this *Session) IsVertify() bool {
	if !this.HasKey(SessionKeyUUID) {
		return false
	}
	if this.GetUUID() == "" {
		return false
	}
	return true
}

// 获取用于在服务器间同步的消息
func (this *Session) getServerSyncMsg() *servercomm.SUpdateSession {
	smsg := &servercomm.SUpdateSession{
		Session:      this.ToMap(),
		ClientConnID: this.GetConnectID(),
		SessionUUID:  this.GetUUID(),
	}
	return smsg
}

// 同步 Session 到 目标模块
func (this *Session) SyncToModule(mod ISInner_SendModuleMsg,
	targetServer string) {
	mod.SInner_SendModuleMsg(targetServer, this.getServerSyncMsg())
}

// 同步 Session 到 所有已绑定的模块
func (this *Session) SyncToBindedModule(mod ISInner_SendModuleMsg) {
	msg := this.getServerSyncMsg()
	this.rangeBinded(func(moduletype string, moduleid string) bool {
		if moduleid != mod.GetModuleID() {
			mod.SInner_SendModuleMsg(moduleid, msg)
		}
		return true
	})
}

func (this *Session) SendMsg(mod ISInner_SendClientMsg, gatemoduletype string,
	msgid uint16, data []byte) {
	mod.SInner_SendClientMsg(this.GetBind(gatemoduletype),
		this.GetConnectID(), msgid, data)
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
		this.set(SessionKey(k), v)
	}
}
