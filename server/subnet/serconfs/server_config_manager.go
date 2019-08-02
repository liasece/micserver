package serconfs

import (
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/servercomm"
	"sync"
)

type ConnInfosManager struct {
	*log.Logger
	ConnInfos   sync.Map // 所需要的所有服务器信息
	ConnInfoSum uint32
}

func (this *ConnInfosManager) GetConnInfoByID(
	serverid string) *servercomm.SServerInfo {
	if value, found := this.ConnInfos.Load(serverid); found {
		return value.(*servercomm.SServerInfo)
	}
	return &servercomm.SServerInfo{}
}

func (this *ConnInfosManager) GetConnInfoByInfo(
	tarinfo *servercomm.SLoginCommand) *servercomm.SServerInfo {
	var res *servercomm.SServerInfo
	if tarinfo.ServerID != "" {
		// 如果已经指定了ID，直接返回
		info := this.GetConnInfoByID(tarinfo.ServerID)
		if info.ServerID == tarinfo.ServerID {
			// 信息正确
			return info
		}
	}
	return res
}

func (this *ConnInfosManager) AddConnInfo(
	newinfo *servercomm.SServerInfo) {
	if newinfo.ServerID == "" {
		log.Error("[ConnInfosManager.AddConnInfo] "+
			"尝试添加一个ID为空的服务器 拒绝 Info[%s]", newinfo.GetJson())
		return
	}
	log.Debug("[ConnInfosManager.AddConnInfo] "+
		"添加服务器信息 Info[%s]", newinfo.GetJson())
	if _, finded := this.ConnInfos.Load(newinfo.ServerID); !finded {
		this.ConnInfoSum++
	}
	this.ConnInfos.Store(newinfo.ServerID, newinfo)
}

func (this *ConnInfosManager) RemoveConnInfo(serverid string) {
	this.ConnInfos.Delete(serverid)
}

func (this *ConnInfosManager) RangeConnInfo(
	callback func(*servercomm.SServerInfo) bool) {
	this.ConnInfos.Range(func(tkey interface{},
		tvalue interface{}) bool {
		value := tvalue.(*servercomm.SServerInfo)
		return callback(value)
	})
}

func (this *ConnInfosManager) ExistConnInfo(
	info *servercomm.SServerInfo) bool {
	tconfig, finded := this.ConnInfos.Load(info.ServerID)
	config := tconfig.(*servercomm.SServerInfo)
	if !finded {
		return false
	}
	if config.ServerID != info.ServerID {
		return false
	}
	return true
}

func (this *ConnInfosManager) CleanConnInfo() {
	this.ConnInfoSum = 0
	this.ConnInfos = sync.Map{}
}

func (this *ConnInfosManager) CountConnInfo() uint32 {
	return this.ConnInfoSum
}
