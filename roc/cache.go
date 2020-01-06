/*
每个 micserver 进程会持有一份 ROC 缓存单例，维护了所有已知的ROC对象所处的位置
*/
package roc

import (
	"math/rand"
	"sync"
)

// ROC缓存分组数量
const (
	CACHE_POOL_GROUP_SUM = 8
)

type catchServerInfo struct {
	moduleid string
}

type serverInfoMap map[string]*catchServerInfo
type objIDToServerMap map[string]*catchServerInfo

// ROC缓存管理器
type Cache struct {
	catchServer serverInfoMap
	catchType   map[ROCObjType]objIDToServerMap
	mutex       sync.Mutex
}

var _gCache Cache

// 获取ROC缓存
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
	f func(id string, location string) bool,
	limitModuleIDs map[string]bool) {
	// 防止 f 中调用其他加锁函数导致死锁，需要备份map
	back := make(objIDToServerMap)

	this.mutex.Lock()
	m := this.catchGetTypeMust(objType)
	for id, v := range m {
		if limitModuleIDs == nil || limitModuleIDs[v.moduleid] == true {
			back[id] = v
		}
	}
	this.mutex.Unlock()

	for id, v := range back {
		if !f(id, v.moduleid) {
			break
		}
	}
	return
}

// 随机获取一个目标类型的缓存对象ID
func (this *Cache) RandomObjIDByType(objType ROCObjType,
	limitModuleIDs map[string]bool) string {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	m := this.catchGetTypeMust(objType)
	tmplist := make([]string, 0)
	for id, v := range m {
		if limitModuleIDs == nil || limitModuleIDs[v.moduleid] == true {
			tmplist = append(tmplist, id)
		}
	}
	lenm := len(tmplist)
	if lenm > 0 {
		// 只有目标类型的对象超过一个，才能从中随机一个
		n := rand.Intn(lenm)
		return tmplist[n]
	}
	return ""
}
