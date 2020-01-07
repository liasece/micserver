/*
micserver中的子网信息，管理了所有模块间的连接
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

// 服务器子网连接管理器
type SubnetManager struct {
	*log.Logger
	// 服务器连接池
	connect.ServerPool
	// 配置信息
	connInfos    serconfs.ConnInfosManager // 所有服务器信息
	moudleConf   *conf.ModuleConfig
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

// 根据模块配置初始化子网连接管理器
func (this *SubnetManager) Init(moudleConf *conf.ModuleConfig) {
	this.myServerInfo = &servercomm.ModuleInfo{}
	this.moudleConf = moudleConf
	this.ServerPool.Logger = this.Logger
	// 初始化消息处理队列
	this.InitMsgQueue(int32(moudleConf.GetInt64(conf.MsgThreadNum)))
	// 我的服务器信息
	this.myServerInfo.ModuleID = this.moudleConf.ID
	this.connInfos.Logger = this.Logger
	// 初始化连接
	this.BindTCPSubnet(this.moudleConf)
	this.BindChanSubnet(this.moudleConf)
}

// 设置子网事件监听者
func (this *SubnetManager) HookSubnet(subnetHook base.SubnetHook) {
	this.subnetHook = subnetHook
}

// 指定类型的获取最新版本的服务器版本号
func (this *SubnetManager) GetLatestVersionConnInfoByType(servertype string) uint64 {
	latestVersion := uint64(0)
	this.connInfos.Range(func(value *servercomm.ModuleInfo) bool {
		if util.GetModuleIDType(value.ModuleID) == servertype &&
			value.Version > latestVersion {
			latestVersion = value.Version
		}
		return true
	})
	return latestVersion
}

// 发送当前已连接的所有服务器信息到目标连接
func (this *SubnetManager) NotifyAllServerInfo(server *connect.Server) {
	retmsg := &servercomm.SNotifyAllInfo{}
	retmsg.ServerInfos = make([]*servercomm.ModuleInfo, 0)
	this.connInfos.Range(func(value *servercomm.ModuleInfo) bool {
		retmsg.ServerInfos = append(retmsg.ServerInfos, value)
		return true
	})
	if len(retmsg.ServerInfos) > 0 {
		this.Debug("[NotifyAllServerInfo] 发送所有服务器列表信息 Msg[%s]",
			retmsg.GetJson())
		server.SendCmd(retmsg)
	}
}
