package roc

import (
	"sync"
)

const (
	CACHE_POOL_GROUP_SUM = 8
)

type catchServerInfo struct {
	serverid string
}

type Cache struct {
	catchServer sync.Map
	catchType   sync.Map
}

var _gCache *Cache

// single case
func init() {
	_gCache = &Cache{}
}

func GetCache() *Cache {
	return _gCache
}

func (this *Cache) catchGetTypeMust(objType string) *sync.Map {
	if vi, ok := this.catchType.Load(objType); !ok {
		vi, _ := this.catchType.LoadOrStore(objType, &sync.Map{})
		return vi.(*sync.Map)
	} else {
		return vi.(*sync.Map)
	}
}

func (this *Cache) catchGetServerMust(serverid string) *catchServerInfo {
	if vi, ok := this.catchServer.Load(serverid); !ok {
		vi, _ := this.catchServer.LoadOrStore(serverid, &catchServerInfo{
			serverid: serverid,
		})
		return vi.(*catchServerInfo)
	} else {
		return vi.(*catchServerInfo)
	}
}

func (this *Cache) Set(objType string, objID string, serverid string) {
	m := this.catchGetTypeMust(objType)
	server := this.catchGetServerMust(serverid)
	m.Store(objID, server)
}

func (this *Cache) SetM(objType string, objIDs []string, serverid string) {
	m := this.catchGetTypeMust(objType)
	server := this.catchGetServerMust(serverid)
	for _, v := range objIDs {
		m.Store(v, server)
	}
}

func (this *Cache) Get(objType string, objID string) string {
	m := this.catchGetTypeMust(objType)
	if vi, ok := m.Load(objID); ok && vi != nil {
		return (vi.(*catchServerInfo)).serverid
	}
	return ""
}
