package serconfs

import (
	"github.com/liasece/micserver/comm"
	"github.com/liasece/micserver/log"
	"sync"
)

// 服务器配置器接口
// 用来提供给superserver动态生成服务器配置
type ITopConfigBuilder interface {
	NewTopConfig(serverip string, servertype uint32,
		serverNumber uint32, serverID uint32,
		serverPort uint32) *comm.SServerInfo
	DeleteTopConfig(*comm.SServerInfo)
}

type TopConfigsManager struct {
	TopConfigs       sync.Map // 所需要的所有服务器信息
	TopConfigSum     uint32
	TopConfigBuilder ITopConfigBuilder
}

func (this *TopConfigsManager) SetTopConfigBuilder(
	builder ITopConfigBuilder) {
	this.TopConfigBuilder = builder
}

func (this *TopConfigsManager) GetTopConfigByID(
	serverid uint32) comm.SServerInfo {
	if value, found := this.TopConfigs.Load(serverid); found {
		return value.(comm.SServerInfo)
	}
	var temp comm.SServerInfo
	return temp
}

func (this *TopConfigsManager) GetTopConfigByInfo(
	tarinfo *comm.SLoginCommand) comm.SServerInfo {
	var res comm.SServerInfo
	if tarinfo.Serverid != 0 {
		// 如果已经指定了ID，直接返回
		info := this.GetTopConfigByID(tarinfo.Serverid)
		if info.Servertype == tarinfo.Servertype &&
			info.Serverid == tarinfo.Serverid {
			// 信息正确
			return info
		}
	}
	if this.TopConfigBuilder != nil {
		res = *this.TopConfigBuilder.NewTopConfig(tarinfo.Serverip,
			tarinfo.Servertype, tarinfo.ServerNumber, tarinfo.Serverid,
			tarinfo.Serverport)
		this.AddTopConfig(res)
	}
	return res
}

func (this *TopConfigsManager) AddTopConfig(
	newinfo comm.SServerInfo) {
	if newinfo.Serverid == 0 {
		log.Error("[TopConfigsManager.AddTopConfig] "+
			"尝试添加一个ID为0的服务器 拒绝 Info[%s]", newinfo.GetJson())
		return
	}
	// log.Debug("[TopConfigsManager.AddTopConfig] "+
	// 	"添加服务器信息 Info[%s]", newinfo.GetJson())
	if _, finded := this.TopConfigs.Load(newinfo.Serverid); !finded {
		this.TopConfigSum++
	}
	this.TopConfigs.Store(newinfo.Serverid, newinfo)
}

func (this *TopConfigsManager) RemoveTopConfig(serverid uint32) {
	info := this.GetTopConfigByID(serverid)
	if info.Serverid != 0 && this.TopConfigBuilder != nil {
		this.TopConfigBuilder.DeleteTopConfig(&info)
	}
	this.TopConfigs.Delete(serverid)
}

func (this *TopConfigsManager) RangeTopConfig(
	callback func(comm.SServerInfo) bool) {
	this.TopConfigs.Range(func(tkey interface{},
		tvalue interface{}) bool {
		value := tvalue.(comm.SServerInfo)
		return callback(value)
	})
}

func (this *TopConfigsManager) ExistTopConfig(
	info comm.SServerInfo) bool {
	tconfig, finded := this.TopConfigs.Load(info.Serverid)
	config := tconfig.(comm.SServerInfo)
	if !finded {
		return false
	}
	if config.Serverid != info.Serverid {
		return false
	}
	return true
}

func (this *TopConfigsManager) CleanTopConfig() {
	this.TopConfigSum = 0
	this.TopConfigs = sync.Map{}
}

func (this *TopConfigsManager) CountTopConfig() uint32 {
	return this.TopConfigSum
}
