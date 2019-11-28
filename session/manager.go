package session

import (
	"sync"
)

type SessionManager struct {
	sessions sync.Map
}

func (this *SessionManager) get(uuid string) *Session {
	if vi, ok := this.sessions.Load(uuid); ok {
		if vi == nil {
			return nil
		}
		return vi.(*Session)
	}
	return nil
}

func (this *SessionManager) store(uuid string, session *Session) {
	this.sessions.Store(uuid, session)
}

func (this *SessionManager) Store(uuid string, session *Session) {
	this.store(uuid, session)
}

func (this *SessionManager) loadOrStore(uuid string,
	session *Session) (*Session, bool) {
	vi, isload := this.sessions.LoadOrStore(uuid, session)
	var v *Session
	if vi != nil {
		v = vi.(*Session)
	}
	return v, isload
}

func (this *SessionManager) delete(uuid string) {
	this.sessions.Delete(uuid)
}

func (this *SessionManager) LoadOrStore(uuid string,
	session *Session) (*Session, bool) {
	return this.loadOrStore(uuid, session)
}

func (this *SessionManager) GetSession(uuid string) *Session {
	return this.get(uuid)
}

// 更新一个 session 到管理器中
func (this *SessionManager) MustUpdateFromMap(targetSession *Session,
	data map[string]string) {
	uuid := targetSession.GetUUID()
	if uuid == "" {
		panic("targetSession uuid must exist")
		return
	}
	session, _ := this.LoadOrStore(uuid, targetSession)
	session.FromMap(data)
}

func (this *SessionManager) DeleteSession(uuid string) {
	this.delete(uuid)
}

func (this *SessionManager) UpdateSessionUUID(uuid string, session *Session) {
	if session == nil {
		return
	}
	olduuid := session.GetUUID()
	var oldsession *Session
	if olduuid != "" {
		oldsession = this.GetSession(olduuid)
	}
	if oldsession != nil {
		session.OnlyAddKeyFromSession(oldsession)
	}
	if uuid != olduuid {
		this.DeleteSession(olduuid)
	}
	session.setUUID(uuid)
	this.Store(uuid, session)
}
