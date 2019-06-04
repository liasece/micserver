package serconfs

import (
	"base/logger"
	"servercomm"
	"sync"
)

// 服务器配置器接口
// 用来提供给superserver动态生成服务器配置
type IServerConfigBuilder interface {
	NewServerConfig(serverip string, servertype uint32,
		serverNumber uint32, serverID uint32,
		serverPort uint32) *servercomm.SServerInfo
	DeleteServerConfig(*servercomm.SServerInfo)
}

type ServerConfigsManager struct {
	serverconfigs       sync.Map // 所需要的所有服务器信息
	serverConfigSum     uint32
	serverConfigBuilder IServerConfigBuilder
}

func (this *ServerConfigsManager) SetServerConfigBuilder(
	builder IServerConfigBuilder) {
	this.serverConfigBuilder = builder
}

func (this *ServerConfigsManager) GetServerConfigByID(
	serverid uint32) servercomm.SServerInfo {
	if value, found := this.serverconfigs.Load(serverid); found {
		return value.(servercomm.SServerInfo)
	}
	var temp servercomm.SServerInfo
	return temp
}

func (this *ServerConfigsManager) GetServerConfigByInfo(
	tarinfo *servercomm.SLoginCommand) servercomm.SServerInfo {
	var res servercomm.SServerInfo
	if tarinfo.Serverid != 0 {
		// 如果已经指定了ID，直接返回
		info := this.GetServerConfigByID(tarinfo.Serverid)
		if info.Servertype == tarinfo.Servertype &&
			info.Serverid == tarinfo.Serverid {
			// 信息正确
			return info
		}
	}
	if this.serverConfigBuilder != nil {
		res = *this.serverConfigBuilder.NewServerConfig(tarinfo.Serverip,
			tarinfo.Servertype, tarinfo.ServerNumber, tarinfo.Serverid,
			tarinfo.Serverport)
		this.AddServerConfig(res)
	}
	return res
}

func (this *ServerConfigsManager) AddServerConfig(
	newinfo servercomm.SServerInfo) {
	if newinfo.Serverid == 0 {
		logger.Error("[ServerConfigsManager.AddServerConfig] "+
			"尝试添加一个ID为0的服务器 拒绝 Info[%s]", newinfo.GetJson())
		return
	}
	// logger.Debug("[ServerConfigsManager.AddServerConfig] "+
	// 	"添加服务器信息 Info[%s]", newinfo.GetJson())
	if _, finded := this.serverconfigs.Load(newinfo.Serverid); !finded {
		this.serverConfigSum++
	}
	this.serverconfigs.Store(newinfo.Serverid, newinfo)
}

func (this *ServerConfigsManager) RemoveServerConfig(serverid uint32) {
	info := this.GetServerConfigByID(serverid)
	if info.Serverid != 0 && this.serverConfigBuilder != nil {
		this.serverConfigBuilder.DeleteServerConfig(&info)
	}
	this.serverconfigs.Delete(serverid)
}

func (this *ServerConfigsManager) RangeServerConfig(
	callback func(servercomm.SServerInfo) bool) {
	this.serverconfigs.Range(func(tkey interface{},
		tvalue interface{}) bool {
		value := tvalue.(servercomm.SServerInfo)
		return callback(value)
	})
}

func (this *ServerConfigsManager) ExistServerConfig(
	info servercomm.SServerInfo) bool {
	tconfig, finded := this.serverconfigs.Load(info.Serverid)
	config := tconfig.(servercomm.SServerInfo)
	if !finded {
		return false
	}
	if config.Serverid != info.Serverid {
		return false
	}
	return true
}

func (this *ServerConfigsManager) CleanServerConfig() {
	this.serverConfigSum = 0
	this.serverconfigs = sync.Map{}
}

func (this *ServerConfigsManager) CountServerConfig() uint32 {
	return this.serverConfigSum
}
