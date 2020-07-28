// Package roc 每个 micserver 进程会持有一份 ROC 缓存单例，维护了所有已知的ROC对象所处的位置
package roc

import (
	"math/rand"
	"sync"
)

// ROC缓存分组数量
const (
	CachePoolGroupSum = 8
)

type catchServerInfo struct {
	moduleid string
}

type serverInfoMap map[string]*catchServerInfo
type objIDToServerMap map[string]*catchServerInfo

// Cache ROC缓存管理器
type Cache struct {
	catchServer serverInfoMap
	catchType   map[ObjType]objIDToServerMap
	mutex       sync.Mutex
}

var _gCache Cache

// GetCache 获取ROC缓存
func GetCache() *Cache {
	return &_gCache
}

func (c *Cache) catchGetTypeMust(objType ObjType) objIDToServerMap {
	if c.catchType == nil {
		c.catchType = make(map[ObjType]objIDToServerMap)
	}
	if v, ok := c.catchType[objType]; !ok {
		v = make(objIDToServerMap)
		c.catchType[objType] = v
		return v
	}
	return c.catchType[objType]
}

func (c *Cache) catchGetServerMust(moduleid string) *catchServerInfo {
	if c.catchServer == nil {
		c.catchServer = make(serverInfoMap)
	}
	if v, ok := c.catchServer[moduleid]; !ok {
		v = &catchServerInfo{
			moduleid: moduleid,
		}
		c.catchServer[moduleid] = v
		return v
	}
	return nil
}

// Set 添加目标对象
func (c *Cache) Set(objType ObjType, objID string, moduleid string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	m := c.catchGetTypeMust(objType)
	server := c.catchGetServerMust(moduleid)
	m[objID] = server
}

// SetM 同时添加多个
func (c *Cache) SetM(objType ObjType, objIDs []string, moduleid string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	m := c.catchGetTypeMust(objType)
	server := c.catchGetServerMust(moduleid)
	for _, v := range objIDs {
		m[v] = server
	}
}

// Del 删除目标对象
func (c *Cache) Del(objType ObjType, objID string, moduleid string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	m := c.catchGetTypeMust(objType)
	if info, ok := m[objID]; ok && info.moduleid == moduleid {
		delete(m, objID)
	}
}

// DelM 同时删除多个
func (c *Cache) DelM(objType ObjType, objIDs []string, moduleid string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	m := c.catchGetTypeMust(objType)
	for _, v := range objIDs {
		if info, ok := m[v]; ok && info.moduleid == moduleid {
			delete(m, v)
		}
	}
}

// Get 获取缓存的目标对象在哪个模块上
func (c *Cache) Get(objType ObjType, objID string) string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	m := c.catchGetTypeMust(objType)
	if v, ok := m[objID]; ok && v != nil {
		return v.moduleid
	}
	return ""
}

// RangeByType 遍历指定类型的ROC对象
func (c *Cache) RangeByType(objType ObjType,
	f func(id string, location string) bool,
	limitModuleIDs map[string]bool) {
	// 防止 f 中调用其他加锁函数导致死锁，需要备份map
	back := make(objIDToServerMap)

	c.mutex.Lock()
	m := c.catchGetTypeMust(objType)
	for id, v := range m {
		if limitModuleIDs == nil || limitModuleIDs[v.moduleid] == true {
			back[id] = v
		}
	}
	c.mutex.Unlock()

	for id, v := range back {
		if !f(id, v.moduleid) {
			break
		}
	}
	return
}

// RandomObjIDByType 随机获取一个目标类型的缓存对象ID
func (c *Cache) RandomObjIDByType(objType ObjType,
	limitModuleIDs map[string]bool) string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	m := c.catchGetTypeMust(objType)
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
