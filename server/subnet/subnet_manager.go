package subnet

import (
	"github.com/liasece/micserver/comm"
	"github.com/liasece/micserver/conf"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/server/subnet/serconfs"
	"github.com/liasece/micserver/tcpconn"
	"github.com/liasece/micserver/util"
	"strconv"
	"sync"
)

type ConnectMsgQueueStruct struct {
	task *tcpconn.ServerConn
	msg  *msg.MessageBinary
}

func CheckServerType(servertype uint32) bool {
	if servertype <= 0 || servertype > 10 {
		return false
	}
	return true
}

// websocket连接管理器
type SubnetManager struct {
	*log.Logger

	connInfos serconfs.ConnInfosManager // 所有服务器信息

	moudleConf *conf.ModuleConfig

	connPool tcpconn.ServerConnPool

	// runningMsgQueue  int
	runningMsgChan   []chan *ConnectMsgQueueStruct
	maxRunningMsgNum int32

	connectMutex   sync.Mutex
	serverexitchan map[string]chan bool

	// lastchan         int
	lastwarningtime1 uint32
	lastwarningtime2 uint32

	myServerInfo comm.SServerInfo
}

func (this *SubnetManager) InitManager(moudleConf *conf.ModuleConfig) {
	this.moudleConf = moudleConf
	this.connPool.Logger = this.Logger
	// 初始化连接
	this.BindTCPSubnet(this.moudleConf.Settings)

	if msgthreadnumstr := moudleConf.GetSetting("msgthreadnum"); msgthreadnumstr != "" {
		msgthreadnum, err := strconv.Atoi(msgthreadnumstr)
		if err == nil {
			this.InitMsgQueue(int32(msgthreadnum))
		}
	}
	this.myServerInfo.ServerID = this.moudleConf.ID
}

func (this *SubnetManager) GetLatestVersionConnInfoByType(servertype string) uint64 {
	latestVersion := uint64(0)
	this.connInfos.RangeConnInfo(
		func(value comm.SServerInfo) bool {
			if util.GetServerIDType(value.ServerID) == servertype &&
				value.Version > latestVersion {
				latestVersion = value.Version
			}
			return true
		})
	return latestVersion
}

func (this *SubnetManager) BroadcastByType(
	servertype string, v msg.MsgStruct) {
	this.connPool.BroadcastByType(servertype, v)
}

func (this *SubnetManager) BroadcastAll(v msg.MsgStruct) {
	this.connPool.BroadcastCmd(v)
}

//通知所有服务器列表信息
func (this *SubnetManager) NotifyAllServerInfo(
	tcptask *tcpconn.ServerConn) {
	retmsg := &comm.SNotifyAllInfo{}
	retmsg.Serverinfos = make([]comm.SServerInfo, 0)
	this.connInfos.RangeConnInfo(func(
		value comm.SServerInfo) bool {
		retmsg.Serverinfos = append(retmsg.Serverinfos, value)
		return true
	})
	if len(retmsg.Serverinfos) > 0 {
		this.Debug("[NotifyAllServerInfo] 发送所有服务器列表信息 Msg[%s]",
			retmsg.GetJson())
		tcptask.SendCmd(retmsg)
	}
}

// 广播消息
func (this *SubnetManager) BroadcastCmd(v msg.MsgStruct) {
	this.connPool.BroadcastCmd(v)
}
