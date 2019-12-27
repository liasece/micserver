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

// 系统中默认的一些 Session 的键
const (
	// 索引绑定的服务器，仅是头部，该Key需要在后面拼接目标索引的Module类型
	SessionKeyBindHead SessionKey = "_s0_bind_"
	// gate 中用于描述链接的唯一ID
	SessionKeyConnectID SessionKey = "_s0_connectid"
	// session 的 UUID 是 session管理器 中的主键
	SessionKeyUUID SessionKey = "_s0_uuid"
)

// 用于提供给 session 向客户端发送消息或者执行某些操作的接口
// 一般情况下，提供 base.Module 即可
type IModuleSessionOptions interface {
	GetModuleID() string
	SInner_SendModuleMsg(gate string, msg msg.MsgStruct)
	SInner_SendClientMsg(gateid string, connectid string, msgid uint16,
		data []byte)
	SInner_CloseSessionConnect(gateid string, connectid string)
}

// 从一个Map结构中实例化一个session
func NewSessionFromMap(session map[string]string) *Session {
	res := &Session{}
	res.FromMap(session)
	return res
}

// 以session的格式从一个Map结构中获取键的值
func getFromMap(session map[string]string, key SessionKey) string {
	if v, ok := session[string(key)]; ok {
		return v
	}
	return ""
}

// 以session的格式从一个Map结构中获取UUID
func GetUUIDFromMap(session map[string]string) string {
	return getFromMap(session, SessionKeyUUID)
}

// 以session的格式从一个Map结构中获取绑定的服务器
func GetBindFromMap(session map[string]string,
	moduleType string) string {
	return getFromMap(session, SessionKeyBindHead+SessionKey(moduleType))
}

// 以session的格式从一个Map结构中获取gate中链接的ID
func GetConnectIDFromMap(session map[string]string) string {
	return getFromMap(session, SessionKeyConnectID)
}

type Session struct {
	m sync.Map
}

func (this *Session) get(key SessionKey) string {
	if vi, ok := this.m.Load(string(key)); ok {
		return vi.(string)
	}
	return ""
}

func (this *Session) set(key SessionKey, value string) {
	this.m.Store(string(key), value)
}

// 遍历所有已绑定的模块
func (this *Session) rangeBinded(
	f func(moduletype string, moduleid string) bool) {
	this.m.Range(func(ki, vi interface{}) bool {
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

// 由于UUID关系到Session的管理，所以不可以直接设置session的UUID，
// 应该通过SessionManager设置，或者使用不推荐的 Set(SessionKeyUUID,uuid) 设置
func (this *Session) setUUID(value string) {
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

// 同步 Session 到 所有已绑定的模块
func (this *Session) SyncToBindedModule(mod IModuleSessionOptions) {
	msg := this.getServerSyncMsg()
	msg.FromModuleID = mod.GetModuleID()
	msg.ToModuleID = "*binded*"
	this.rangeBinded(func(moduletype string, moduleid string) bool {
		if moduleid != mod.GetModuleID() {
			mod.SInner_SendModuleMsg(moduleid, msg)
		}
		return true
	})
}

func (this *Session) SendMsg(mod IModuleSessionOptions, gatemoduletype string,
	msgid uint16, data []byte) {
	mod.SInner_SendClientMsg(this.GetBind(gatemoduletype),
		this.GetConnectID(), msgid, data)
}

func (this *Session) CloseSessionConnect(mod IModuleSessionOptions,
	gatemoduletype string) {
	mod.SInner_CloseSessionConnect(this.GetBind(gatemoduletype),
		this.GetConnectID())
}

func (this *Session) ToMap() map[string]string {
	res := make(map[string]string)
	this.m.Range(func(ki, vi interface{}) bool {
		res[ki.(string)] = vi.(string)
		return true
	})
	return res
}

// 用Map中有的键的值替换当前session中的值
func (this *Session) FromMap(m map[string]string) {
	for k, v := range m {
		this.set(SessionKey(k), v)
	}
}

// 将dir中有且this中没有的键增加到this中，不会修改任何this中已有的键值，
// 只会增加this的键值。
// 这是一种简单的通过另一个Session完善当前Session的方法。
func (this *Session) OnlyAddKeyFromSession(dir *Session) {
	dir.m.Range(func(ki, vi interface{}) bool {
		this.m.LoadOrStore(ki, vi)
		return true
	})
}
