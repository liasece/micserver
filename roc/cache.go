package roc

import (
	"math/rand"
	"sync"
)

const (
	CACHE_POOL_GROUP_SUM = 8
)

type catchServerInfo struct {
	moduleid string
}

type serverInfoMap map[string]*catchServerInfo
type objIDToServerMap map[string]*catchServerInfo

type Cache struct {
	catchServer serverInfoMap
	catchType   map[ROCObjType]objIDToServerMap
	mutex       sync.Mutex
}

var _gCache Cache

func GetCache() *Cache {
	return &_gCache
}

func (this *Cache) catchGetTypeMust(objType ROCObjType) objIDToServerMap {
	if this.catchType == nil {
		this.catchType = make(map[ROCObjType]objIDToServerMap)
	}
	if v, ok := this.catchType[objType]; !ok {
		v = make(objIDToServerMap)
		this.catchType[objType] = v
		return v
	} else {
		return v
	}
}

func (this *Cache) catchGetServerMust(moduleid string) *catchServerInfo {
	if this.catchServer == nil {
		this.catchServer = make(serverInfoMap)
	}
	if v, ok := this.catchServer[moduleid]; !ok {
		v = &catchServerInfo{
			moduleid: moduleid,
		}
		this.catchServer[moduleid] = v
		return v
	} else {
		return v
	}
}

// 添加目标对象
func (this *Cache) Set(objType ROCObjType, objID string, moduleid string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	m := this.catchGetTypeMust(objType)
	server := this.catchGetServerMust(moduleid)
	m[objID] = server
}

// 同时添加多个
func (this *Cache) SetM(objType ROCObjType, objIDs []string, moduleid string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	m := this.catchGetTypeMust(objType)
	server := this.catchGetServerMust(moduleid)
	for _, v := range objIDs {
		m[v] = server
	}
}

// 删除目标对象
func (this *Cache) Del(objType ROCObjType, objID string, moduleid string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	m := this.catchGetTypeMust(objType)
	if info, ok := m[objID]; ok && info.moduleid == moduleid {
		delete(m, objID)
	}
}

// 同时删除多个
func (this *Cache) DelM(objType ROCObjType, objIDs []string, moduleid string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	m := this.catchGetTypeMust(objType)
	for _, v := range objIDs {
		if info, ok := m[v]; ok && info.moduleid == moduleid {
			delete(m, v)
		}
	}
}

// 获取缓存的目标对象在哪个模块上
func (this *Cache) Get(objType ROCObjType, objID string) string {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	m := this.catchGetTypeMust(objType)
	if v, ok := m[objID]; ok && v != nil {
		return v.moduleid
	}
	return ""
}

// 遍历指定类型的ROC对象
func (this *Cache) RangeByType(objType ROCObjType,
	f func(id string, location string) bool) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	m := this.catchGetTypeMust(objType)
	for id, v := range m {
		if !f(id, v.moduleid) {
			break
		}
	}
	return
}

// 随机获取一个目标类型的缓存对象ID
func (this *Cache) RandomObjIDByType(objType ROCObjType) string {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	m := this.catchGetTypeMust(objType)
	lenm := len(m)
	if lenm > 0 {
		// 只有目标类型的对象超过一个，才能从中随机一个
		n := rand.Intn(lenm)
		for id, _ := range m {
			if n <= 0 {
				return id
			}
			n--
		}
	}
	return ""
}
