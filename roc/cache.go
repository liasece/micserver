// Package roc 每个 micserver 进程会持有一份 ROC 缓存单例，维护了所有已知的ROC对象所处的位置
package roc

import (
	"errors"
	"math/rand"
	"sync"
)

// ROC缓存分组数量
const (
	CachePoolGroupSum = 8
)

type catchServerInfo struct {
	moduleID string
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
	} else {
		return v
	}
}

func (c *Cache) catchGetServerMust(moduleID string) *catchServerInfo {
	if c.catchServer == nil {
		c.catchServer = make(serverInfoMap)
	}
	if v, ok := c.catchServer[moduleID]; !ok {
		v = &catchServerInfo{
			moduleID: moduleID,
		}
		c.catchServer[moduleID] = v
		return v
	} else {
		return v
	}
}

// Set 添加目标对象
func (c *Cache) Set(objType ObjType, objID string, moduleID string) error {
	if objType == "" || objID == "" || moduleID == "" {
		return errors.New("objType or objID or moduleID can't be empty")
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()

	m := c.catchGetTypeMust(objType)
	server := c.catchGetServerMust(moduleID)
	m[objID] = server
	return nil
}

// SetM 同时添加多个
func (c *Cache) SetM(objType ObjType, objIDs []string, moduleID string) error {
	if objType == "" || moduleID == "" {
		return errors.New("objType or objID or moduleID can't be empty")
	}
	for _, v := range objIDs {
		if v == "" {
			return errors.New("objType or objID or moduleID can't be empty")
		}
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	m := c.catchGetTypeMust(objType)
	server := c.catchGetServerMust(moduleID)
	for _, v := range objIDs {
		m[v] = server
	}
	return nil
}

// Del 删除目标对象
func (c *Cache) Del(objType ObjType, objID string, moduleID string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	m := c.catchGetTypeMust(objType)
	if info, ok := m[objID]; ok && info.moduleID == moduleID {
		delete(m, objID)
		return true
	}
	return false
}

// DelM 同时删除多个
// 返回成功删除的数量
func (c *Cache) DelM(objType ObjType, objIDs []string, moduleID string) int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	res := 0
	m := c.catchGetTypeMust(objType)
	for _, v := range objIDs {
		if info, ok := m[v]; ok && info.moduleID == moduleID {
			delete(m, v)
			res++
		}
	}
	return res
}

// Get 获取缓存的目标对象在哪个模块上
func (c *Cache) Get(objType ObjType, objID string) string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	m := c.catchGetTypeMust(objType)
	if v, ok := m[objID]; ok && v != nil {
		return v.moduleID
	}
	return ""
}

// RangeByType 遍历指定类型的ROC对象
func (c *Cache) RangeByType(objType ObjType, f func(id string, location string) bool, limitModuleIDs map[string]bool) {
	// 防止 f 中调用其他加锁函数导致死锁，需要备份map
	back := make(objIDToServerMap)

	c.mutex.Lock()
	m := c.catchGetTypeMust(objType)
	for id, v := range m {
		if limitModuleIDs == nil || limitModuleIDs[v.moduleID] == true {
			back[id] = v
		}
	}
	c.mutex.Unlock()

	for id, v := range back {
		if !f(id, v.moduleID) {
			break
		}
	}
	return
}

// RandomObjIDByType 随机获取一个目标类型的缓存对象ID
func (c *Cache) RandomObjIDByType(objType ObjType, limitModuleIDs map[string]bool) string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	m := c.catchGetTypeMust(objType)
	tempList := make([]string, 0)
	for id, v := range m {
		if limitModuleIDs == nil || limitModuleIDs[v.moduleID] == true {
			tempList = append(tempList, id)
		}
	}
	lenm := len(tempList)
	if lenm > 0 {
		// 只有目标类型的对象超过一个，才能从中随机一个
		n := rand.Intn(lenm)
		return tempList[n]
	}
	return ""
}
