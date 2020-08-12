/*
Package subnet micserver中的子网信息，管理了所有模块间的连接
*/
package subnet

import (
	"sync"

	"github.com/liasece/micserver/conf"
	"github.com/liasece/micserver/connect"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/server/subnet/base"
	"github.com/liasece/micserver/server/subnet/serconfs"
	"github.com/liasece/micserver/servercomm"
	"github.com/liasece/micserver/util"
)

// Manager 服务器子网连接管理器
type Manager struct {
	*log.Logger
	// 服务器连接池
	connect.ServerPool
	// 配置信息
	connInfos    serconfs.ConnInfosManager // 所有服务器信息
	moduleConf   *conf.ModuleConfig
	connectMutex sync.Mutex
	// 服务器重连任务相关
	serverexitchan map[string]chan bool
	// 消息处理相关
	runningMsgChan   []chan *ConnectMsgQueueStruct
	maxRunningMsgNum int32
	// 我的服务器信息
	myServerInfo *servercomm.ModuleInfo
	// 子网系统钩子
	subnetHook base.SubnetHook
}

// Init 根据模块配置初始化子网连接管理器
func (manager *Manager) Init(moduleConf *conf.ModuleConfig) {
	manager.myServerInfo = &servercomm.ModuleInfo{}
	manager.moduleConf = moduleConf
	manager.ServerPool.Logger = manager.Logger
	// 初始化消息处理队列
	manager.InitMsgQueue(int32(moduleConf.GetInt64(conf.MsgThreadNum)))
	// 我的服务器信息
	manager.myServerInfo.ModuleID = manager.moduleConf.ID
	manager.connInfos.Logger = manager.Logger
	// 初始化连接
	manager.BindTCPSubnet(manager.moduleConf)
	manager.BindChanSubnet(manager.moduleConf)
}

// HookSubnet 设置子网事件监听者
func (manager *Manager) HookSubnet(subnetHook base.SubnetHook) {
	manager.subnetHook = subnetHook
}

// GetLatestVersionConnInfoByType 指定类型的获取最新版本的服务器版本号
func (manager *Manager) GetLatestVersionConnInfoByType(servertype string) uint64 {
	latestVersion := uint64(0)
	manager.connInfos.Range(func(value *servercomm.ModuleInfo) bool {
		if util.GetModuleIDType(value.ModuleID) == servertype &&
			value.Version > latestVersion {
			latestVersion = value.Version
		}
		return true
	})
	return latestVersion
}

// NotifyAllServerInfo 发送当前已连接的所有服务器信息到目标连接
func (manager *Manager) NotifyAllServerInfo(server *connect.Server) {
	retmsg := &servercomm.SNotifyAllInfo{}
	retmsg.ServerInfos = make([]*servercomm.ModuleInfo, 0)
	manager.connInfos.Range(func(value *servercomm.ModuleInfo) bool {
		retmsg.ServerInfos = append(retmsg.ServerInfos, value)
		return true
	})
	if len(retmsg.ServerInfos) > 0 {
		manager.Debug("[Manager.NotifyAllServerInfo] Send all server list information", log.Reflect("Msg", retmsg))
		server.SendCmd(retmsg)
	}
}
