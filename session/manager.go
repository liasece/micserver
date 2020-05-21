/*
Package session 管理当前Module所拥有的Session
*/
package session

import (
	"sync"
)

// Manager 当前Module的Session管理器
type Manager struct {
	sessions sync.Map
}

func (sessionManager *Manager) get(uuid string) *Session {
	if vi, ok := sessionManager.sessions.Load(uuid); ok {
		if vi != nil {
			return vi.(*Session)
		}
	}
	return nil
}

// store 保存一个目标UUID的session，必须存在不为空的uuid，否则不会保存
func (sessionManager *Manager) store(session *Session) {
	uuid := session.GetUUID()
	if uuid == "" {
		return
	}
	sessionManager.sessions.Store(uuid, session)
}

// Store 保存一个Session
func (sessionManager *Manager) Store(session *Session) {
	sessionManager.store(session)
}

func (sessionManager *Manager) loadOrStore(uuid string,
	session *Session) (*Session, bool) {
	vi, isload := sessionManager.sessions.LoadOrStore(uuid, session)
	var v *Session
	if vi != nil {
		v = vi.(*Session)
	}
	return v, isload
}

func (sessionManager *Manager) delete(uuid string) {
	sessionManager.sessions.Delete(uuid)
}

// LoadOrStore 加载或保存一个Session，返回操作成功的Session以及是否时加载行为
func (sessionManager *Manager) LoadOrStore(uuid string,
	session *Session) (*Session, bool) {
	return sessionManager.loadOrStore(uuid, session)
}

// GetSession 根据Session的UUID获取一个Session
func (sessionManager *Manager) GetSession(uuid string) *Session {
	return sessionManager.get(uuid)
}

// MustUpdateFromMap 更新一个 session 到管理器中
// targetSession 可以是不在当前管理器中的，但是其必须拥有UUID
// 如果当前管理器中已存在 targetSession.UUID 指定的session，
// 且两者不是同一个Session，会用 targetSession 完善当前管理器中已存在的Session
func (sessionManager *Manager) MustUpdateFromMap(targetSession *Session,
	data map[string]string) {
	uuid := targetSession.GetUUID()
	if uuid == "" {
		panic("targetSession uuid must exist")
	}
	session, isload := sessionManager.LoadOrStore(uuid, targetSession)
	if isload && session != targetSession {
		// 存在两个主键一样的session，用传入session完善本地已存在的session
		session.OnlyAddKeyFromSession(targetSession)
	}
	session.FromMap(data)
}

// DeleteSession 删除目标uuid的session
func (sessionManager *Manager) DeleteSession(uuid string) {
	sessionManager.delete(uuid)
}

// UpdateSessionUUID 更新 session 绑定的UUID，由于 session manager 使用UUID作为索引
// session 的主键，所以UUID的更改需要同时修改manager中的绑定的
func (sessionManager *Manager) UpdateSessionUUID(uuid string, session *Session) {
	if session == nil {
		return
	}
	// 被修改的Session的原UUID需要从本地管理器中加载出来并判断是否需要修改
	olduuid := session.GetUUID()
	if olduuid != "" {
		// 本地管理器中存在原UUID的session
		// 如果本地和目标修改的session不是同一个，需要复制其内容到新的session中
		oldsession := sessionManager.GetSession(olduuid)
		if oldsession != nil {
			if oldsession != session {
				session.OnlyAddKeyFromSession(oldsession)
			}
			// 主键被修改了，需要删除旧的主键
			if uuid != olduuid {
				sessionManager.DeleteSession(olduuid)
			}
		}
	}
	// 如果目标修改为的UUID已存在于本地，需要将本地已存在的Session合并到目标session，
	// 这里不需要删除本地管理器中主键为uuid的session，由后续的 Store() 实现替换
	localsession := sessionManager.GetSession(uuid)
	if localsession != nil && localsession != session {
		// 替换一个新的session对象
		session.OnlyAddKeyFromSession(localsession)
	}

	// 实际处理session的uuid修改
	session.setUUID(uuid)
	sessionManager.Store(session)
}
