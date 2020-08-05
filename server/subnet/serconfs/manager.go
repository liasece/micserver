/*
Package serconfs 连接到本模块的服务器配置信息管理器
*/
package serconfs

import (
	"sync"

	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/servercomm"
)

// ConnInfosManager 连接到本模块的服务器配置信息管理器
type ConnInfosManager struct {
	*log.Logger
	ConnInfos   sync.Map // 所需要的所有服务器信息
	ConnInfoSum uint32
}

// Get 获取目标连接的配置信息，这不是由本地配置决定的，而是由目标方更新过来的
func (connInfosManager *ConnInfosManager) Get(moduleid string) *servercomm.ModuleInfo {
	if value, found := connInfosManager.ConnInfos.Load(moduleid); found {
		return value.(*servercomm.ModuleInfo)
	}
	return &servercomm.ModuleInfo{}
}

// Add 增加一个连接的配置信息
func (connInfosManager *ConnInfosManager) Add(newinfo *servercomm.ModuleInfo) {
	if newinfo.ModuleID == "" {
		log.Error("[ConnInfosManager.AddConnInfo] Try to add a server with an empty ID, denied", log.Reflect("Info", newinfo))
		return
	}
	log.Debug("[ConnInfosManager.AddConnInfo] Adding server information", log.Reflect("Info", newinfo))
	if _, finded := connInfosManager.ConnInfos.Load(newinfo.ModuleID); !finded {
		connInfosManager.ConnInfoSum++
	}
	connInfosManager.ConnInfos.Store(newinfo.ModuleID, newinfo)
}

// Delete 删除一个连接的配置信息
func (connInfosManager *ConnInfosManager) Delete(moduleid string) {
	connInfosManager.ConnInfos.Delete(moduleid)
}

// Range 遍历所有连接的配置信息
func (connInfosManager *ConnInfosManager) Range(callback func(*servercomm.ModuleInfo) bool) {
	connInfosManager.ConnInfos.Range(func(tkey interface{},
		tvalue interface{}) bool {
		value := tvalue.(*servercomm.ModuleInfo)
		return callback(value)
	})
}

// Exist 判断目标信息是否存在
func (connInfosManager *ConnInfosManager) Exist(info *servercomm.ModuleInfo) bool {
	tconfig, finded := connInfosManager.ConnInfos.Load(info.ModuleID)
	config := tconfig.(*servercomm.ModuleInfo)
	if !finded {
		return false
	}
	if config.ModuleID != info.ModuleID {
		return false
	}
	return true
}

// Clean 清空当前配置信息
func (connInfosManager *ConnInfosManager) Clean() {
	connInfosManager.ConnInfoSum = 0
	connInfosManager.ConnInfos = sync.Map{}
}

// Len 当前连接配置信息的数量
func (connInfosManager *ConnInfosManager) Len() uint32 {
	return connInfosManager.ConnInfoSum
}
