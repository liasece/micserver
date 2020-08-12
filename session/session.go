/*
Package session 客户端在连接到服务器网络后，除了Gate能取到客户端的实际连接Client外，
其他模块只能通过客户端的Session操作客户端。
*/
package session

import (
	"fmt"
	"strings"
	"sync"

	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/servercomm"
	"github.com/liasece/micserver/util/conv"
)

// TKey Session中字段的键的类型
type TKey string

// 系统中默认的一些 Session 的键
const (
	// 索引绑定的服务器，仅是头部，该Key需要在后面拼接目标索引的Module类型
	SessionKeyBindHead TKey = "_s0_bind_"
	// gate 中用于描述链接的唯一ID
	SessionKeyConnectID TKey = "_s0_connectid"
	// session 的 UUID 是 session管理器 中的主键
	SessionKeyUUID TKey = "_s0_uuid"
)

// IModuleSessionOptions 用于提供给 session 向客户端发送消息或者执行某些操作的接口
// 一般情况下，提供 base.Module 即可
type IModuleSessionOptions interface {
	GetModuleID() string
	SInnerSendModuleMsg(gate string, msg msg.IMsgStruct)
	SInnerSendClientMsg(gateID string, connectID string, msgID uint16, data []byte)
	SInnerCloseSessionConnect(gateID string, connectID string)
}

// NewSessionFromMap 从一个Map结构中实例化一个session
func NewSessionFromMap(session map[string]string) *Session {
	res := &Session{}
	res.FromMap(session)
	return res
}

// getFromMap 以session的格式从一个Map结构中获取键的值
func getFromMap(session map[string]string, key TKey) string {
	if v, ok := session[string(key)]; ok {
		return v
	}
	return ""
}

// GetUUIDFromMap 以session的格式从一个Map结构中获取UUID
func GetUUIDFromMap(session map[string]string) string {
	return getFromMap(session, SessionKeyUUID)
}

// GetBindFromMap 以session的格式从一个Map结构中获取绑定的服务器
func GetBindFromMap(session map[string]string, moduleType string) string {
	return getFromMap(session, SessionKeyBindHead+TKey(moduleType))
}

// GetConnectIDFromMap 以session的格式从一个Map结构中获取gate中链接的ID
func GetConnectIDFromMap(session map[string]string) string {
	return getFromMap(session, SessionKeyConnectID)
}

// Session 客户端连接会话
type Session struct {
	m sync.Map
}

func (s *Session) get(key TKey) string {
	if vi, ok := s.m.Load(string(key)); ok {
		return vi.(string)
	}
	return ""
}

func (s *Session) set(key TKey, value string) {
	s.m.Store(string(key), value)
}

// rangeBind 遍历所有已绑定的模块
func (s *Session) rangeBind(
	f func(moduletype string, moduleID string) bool) {
	s.m.Range(func(ki, vi interface{}) bool {
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

// GetBindList 获取该 Session 绑定的所有模块
// 返回值 键为模块类型，值为模块ID
func (s *Session) GetBindList() map[string]string {
	res := make(map[string]string)
	s.rangeBind(func(moduletype string, moduleID string) bool {
		res[moduletype] = moduleID
		return true
	})
	return res
}

// GetUUID get this session uuid
func (s *Session) GetUUID() string {
	return s.get(SessionKeyUUID)
}

// 由于UUID关系到Session的管理，所以不可以直接设置session的UUID，
// setUUID 应该通过SessionManager设置，或者使用不推荐的 Set(SessionKeyUUID,uuid) 设置
func (s *Session) setUUID(value string) {
	s.set(SessionKeyUUID, value)
}

// Get 获取指定键的值
func (s *Session) Get(key TKey) string {
	return s.get(key)
}

// GetBool 获取指定键的 bool 值
func (s *Session) GetBool(key TKey) bool {
	return conv.MustInterfaceToBool(s.get(key))
}

// GetInt64 获取指定键的 int64 值
func (s *Session) GetInt64(key TKey) int64 {
	return conv.MustInterfaceToInt64(s.get(key))
}

// Set 设置指定键的值
func (s *Session) Set(key TKey, value string) {
	s.set(key, value)
}

// SetBool 设置指定键的 bool 值
func (s *Session) SetBool(key TKey, value bool) {
	if value {
		s.set(key, "true")
	} else {
		s.set(key, "false")
	}
}

// SetInt64 设置指定键的 int64 值
func (s *Session) SetInt64(key TKey, value int64) {
	s.set(key, fmt.Sprint(value))
}

// GetConnectID 获取Session的客户端连接ID
func (s *Session) GetConnectID() string {
	return s.get(SessionKeyConnectID)
}

// SetConnectID 设置Session的客户端连接ID
func (s *Session) SetConnectID(value string) {
	s.set(SessionKeyConnectID, value)
}

// GetBind 获取当前绑定的指定类型模块的ID
func (s *Session) GetBind(moduleType string) string {
	return s.get(SessionKeyBindHead + TKey(moduleType))
}

// SetBind 设置当前绑定的指定类型模块的ID
func (s *Session) SetBind(moduleType string, value string) {
	s.set(SessionKeyBindHead+TKey(moduleType), value)
}

// HasBind 判断当前是否已经绑定指定类型的模块
func (s *Session) HasBind(moduleType string) bool {
	return s.HasKey(SessionKeyBindHead + TKey(moduleType))
}

// HasKey 判断当前是否存在指定键的值
func (s *Session) HasKey(key TKey) bool {
	_, ok := s.m.Load(string(key))
	return ok
}

// IsVerify 判断当前Session是否已经经过验证，如果一个客户端连接经过了验证，则一定会存在一个
// 用户UUID绑定到此Session上。
func (s *Session) IsVerify() bool {
	if !s.HasKey(SessionKeyUUID) {
		return false
	}
	if s.GetUUID() == "" {
		return false
	}
	return true
}

// getServerSyncMsg 获取用于在服务器间同步的消息
func (s *Session) getServerSyncMsg() *servercomm.SUpdateSession {
	smsg := &servercomm.SUpdateSession{
		Session:      s.ToMap(),
		ClientConnID: s.GetConnectID(),
		SessionUUID:  s.GetUUID(),
	}
	return smsg
}

// SyncToBindModule 同步 Session 到 所有已绑定的模块
func (s *Session) SyncToBindModule(mod IModuleSessionOptions) {
	msg := s.getServerSyncMsg()
	msg.FromModuleID = mod.GetModuleID()
	msg.ToModuleID = "*bind*"
	s.rangeBind(func(moduletype string, moduleID string) bool {
		if moduleID != mod.GetModuleID() {
			mod.SInnerSendModuleMsg(moduleID, msg)
		}
		return true
	})
}

// SendMsg 向该Session指定的客户端发送一个消息
func (s *Session) SendMsg(mod IModuleSessionOptions, gatemoduletype string, msgID uint16, data []byte) {
	mod.SInnerSendClientMsg(s.GetBind(gatemoduletype), s.GetConnectID(), msgID, data)
}

// CloseSessionConnect 请求管理该Session的网关关闭该Session的连接
func (s *Session) CloseSessionConnect(mod IModuleSessionOptions, gatemoduletype string) {
	mod.SInnerCloseSessionConnect(s.GetBind(gatemoduletype), s.GetConnectID())
}

// ToMap 将当前Session的键值到处成为 map[string]string 的类型
func (s *Session) ToMap() map[string]string {
	res := make(map[string]string)
	s.m.Range(func(ki, vi interface{}) bool {
		res[ki.(string)] = vi.(string)
		return true
	})
	return res
}

// FromMap 用Map中有的键的值替换当前session中的值
func (s *Session) FromMap(m map[string]string) {
	for k, v := range m {
		s.set(TKey(k), v)
	}
}

// OnlyAddKeyFromSession 将dir中有且this中没有的键增加到this中，不会修改任何this中已有的键值，
// 只会增加this的键值。
// 这是一种简单的通过另一个Session完善当前Session的方法。
func (s *Session) OnlyAddKeyFromSession(dir *Session) {
	dir.m.Range(func(ki, vi interface{}) bool {
		s.m.LoadOrStore(ki, vi)
		return true
	})
}
