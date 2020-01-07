/*
连接到本模块的服务器配置信息管理器
*/
package serconfs

import (
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/servercomm"
	"sync"
)

// 连接到本模块的服务器配置信息管理器
type ConnInfosManager struct {
	*log.Logger
	ConnInfos   sync.Map // 所需要的所有服务器信息
	ConnInfoSum uint32
}

// 获取目标连接的配置信息，这不是由本地配置决定的，而是由目标方更新过来的
func (this *ConnInfosManager) Get(moduleid string) *servercomm.ModuleInfo {
	if value, found := this.ConnInfos.Load(moduleid); found {
		return value.(*servercomm.ModuleInfo)
	}
	return &servercomm.ModuleInfo{}
}

// 增加一个连接的配置信息
func (this *ConnInfosManager) Add(newinfo *servercomm.ModuleInfo) {
	if newinfo.ModuleID == "" {
		log.Error("[ConnInfosManager.AddConnInfo] "+
			"尝试添加一个ID为空的服务器 拒绝 Info[%s]", newinfo.GetJson())
		return
	}
	log.Debug("[ConnInfosManager.AddConnInfo] "+
		"添加服务器信息 Info[%s]", newinfo.GetJson())
	if _, finded := this.ConnInfos.Load(newinfo.ModuleID); !finded {
		this.ConnInfoSum++
	}
	this.ConnInfos.Store(newinfo.ModuleID, newinfo)
}

// 删除一个连接的配置信息
func (this *ConnInfosManager) Delete(moduleid string) {
	this.ConnInfos.Delete(moduleid)
}

// 遍历所有连接的配置信息
func (this *ConnInfosManager) Range(
	callback func(*servercomm.ModuleInfo) bool) {
	this.ConnInfos.Range(func(tkey interface{},
		tvalue interface{}) bool {
		value := tvalue.(*servercomm.ModuleInfo)
		return callback(value)
	})
}

// 判断目标信息是否存在
func (this *ConnInfosManager) Exist(
	info *servercomm.ModuleInfo) bool {
	tconfig, finded := this.ConnInfos.Load(info.ModuleID)
	config := tconfig.(*servercomm.ModuleInfo)
	if !finded {
		return false
	}
	if config.ModuleID != info.ModuleID {
		return false
	}
	return true
}

// 清空当前配置信息
func (this *ConnInfosManager) Clean() {
	this.ConnInfoSum = 0
	this.ConnInfos = sync.Map{}
}

// 当前连接配置信息的数量
func (this *ConnInfosManager) Len() uint32 {
	return this.ConnInfoSum
}
