package subnet

import (
	"github.com/liasece/micserver/conf"
	"github.com/liasece/micserver/connect"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/server/subnet/serconfs"
	"github.com/liasece/micserver/servercomm"
	"github.com/liasece/micserver/util"
	"strconv"
	"sync"
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
	connect.ServerConnPool
	// 回调注册
	SubnetCallback
	// 配置信息
	connInfos    serconfs.ConnInfosManager // 所有服务器信息
	moudleConf   *conf.ModuleConfig
	connectMutex sync.Mutex
	// 服务器重连任务相关
	serverexitchan map[string]chan bool
	// 消息处理相关
	runningMsgChan   []chan *ConnectMsgQueueStruct
	maxRunningMsgNum int32
	// 防止日志频繁
	lastwarningtime1 uint32
	lastwarningtime2 uint32
	// 我的服务器信息
	myServerInfo *servercomm.SServerInfo
}

func (this *SubnetManager) InitManager(moudleConf *conf.ModuleConfig) {
	this.myServerInfo = &servercomm.SServerInfo{}
	this.moudleConf = moudleConf
	this.ServerConnPool.Logger = this.Logger
	// 初始化连接
	this.BindTCPSubnet(this.moudleConf.Settings)
	// 初始化消息处理队列
	if msgthreadnumstr := moudleConf.GetSetting("msgthreadnum"); msgthreadnumstr != "" {
		msgthreadnum, err := strconv.Atoi(msgthreadnumstr)
		if err == nil {
			this.InitMsgQueue(int32(msgthreadnum))
		}
	}
	// 我的服务器信息
	this.myServerInfo.ServerID = this.moudleConf.ID
	this.connInfos.Logger = this.Logger
}

// 指定类型的获取最新版本的服务器版本号
func (this *SubnetManager) GetLatestVersionConnInfoByType(servertype string) uint64 {
	latestVersion := uint64(0)
	this.connInfos.RangeConnInfo(
		func(value *servercomm.SServerInfo) bool {
			if util.GetServerIDType(value.ServerID) == servertype &&
				value.Version > latestVersion {
				latestVersion = value.Version
			}
			return true
		})
	return latestVersion
}

//通知所有服务器列表信息
func (this *SubnetManager) NotifyAllServerInfo(
	tcptask *connect.ServerConn) {
	retmsg := &servercomm.SNotifyAllInfo{}
	retmsg.Serverinfos = make([]*servercomm.SServerInfo, 0)
	this.connInfos.RangeConnInfo(func(
		value *servercomm.SServerInfo) bool {
		retmsg.Serverinfos = append(retmsg.Serverinfos, value)
		return true
	})
	if len(retmsg.Serverinfos) > 0 {
		this.Debug("[NotifyAllServerInfo] 发送所有服务器列表信息 Msg[%s]",
			retmsg.GetJson())
		tcptask.SendCmd(retmsg)
	}
}
