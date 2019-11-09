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

func CheckServerType(servertype uint32) bool {
	if servertype <= 0 || servertype > 10 {
		return false
	}
	return true
}

// websocket连接管理器
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
	myServerInfo *servercomm.ServerInfo
	// 子网系统钩子
	subnetHook base.SubnetHook
}

func (this *SubnetManager) Init(moudleConf *conf.ModuleConfig) {
	this.myServerInfo = &servercomm.ServerInfo{}
	this.moudleConf = moudleConf
	this.ServerPool.Logger = this.Logger
	// 初始化消息处理队列
	this.InitMsgQueue(int32(moudleConf.GetInt64(conf.MsgThreadNum)))
	// 我的服务器信息
	this.myServerInfo.ServerID = this.moudleConf.ID
	this.connInfos.Logger = this.Logger
	// 初始化连接
	this.BindTCPSubnet(this.moudleConf)
	this.BindChanSubnet(this.moudleConf)
}

func (this *SubnetManager) HookSubnet(subnetHook base.SubnetHook) {
	this.subnetHook = subnetHook
}

// 指定类型的获取最新版本的服务器版本号
func (this *SubnetManager) GetLatestVersionConnInfoByType(servertype string) uint64 {
	latestVersion := uint64(0)
	this.connInfos.RangeConnInfo(
		func(value *servercomm.ServerInfo) bool {
			if util.GetModuleIDType(value.ServerID) == servertype &&
				value.Version > latestVersion {
				latestVersion = value.Version
			}
			return true
		})
	return latestVersion
}

//通知所有服务器列表信息
func (this *SubnetManager) NotifyAllServerInfo(
	tcptask *connect.Server) {
	retmsg := &servercomm.SNotifyAllInfo{}
	retmsg.ServerInfos = make([]*servercomm.ServerInfo, 0)
	this.connInfos.RangeConnInfo(func(
		value *servercomm.ServerInfo) bool {
		retmsg.ServerInfos = append(retmsg.ServerInfos, value)
		return true
	})
	if len(retmsg.ServerInfos) > 0 {
		this.Debug("[NotifyAllServerInfo] 发送所有服务器列表信息 Msg[%s]",
			retmsg.GetJson())
		tcptask.SendCmd(retmsg)
	}
}
