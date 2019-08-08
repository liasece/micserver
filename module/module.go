package module

import (
	"github.com/liasece/micserver/conf"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/server/gate"
	"github.com/liasece/micserver/server/subnet"
	"github.com/liasece/micserver/servercomm"
	"github.com/liasece/micserver/tcpconn"
	"github.com/liasece/micserver/util"
	"time"
)

type IModule interface {
	GetModuleID() string
	InitModule(conf.ModuleConfig)
	InitSubnet(map[string]string)
	AfterInitModule()
	TopRunner()
	KillModule()
	IsStopped() bool
	GetConfiger() *conf.ModuleConfig
}

type BaseModule struct {
	*log.Logger
	util.TimerManager

	ModuleID string
	Configer *conf.ModuleConfig

	subnetManager   *subnet.SubnetManager
	gateBase        *gate.GateBase
	hasKilledModule bool
	hasStopped      bool
}

func (this *BaseModule) InitModule(configer conf.ModuleConfig) {
	this.Configer = &configer
	// 初始化logger
	if this.Configer.HasSetting("logpath") {
		this.Logger = log.NewLogger(this.Configer.GetModuleSettingMap())
		this.SetLogName(this.ModuleID)
	} else {
		this.Logger = log.GetDefaultLogger()
	}
	// 申请内存
	if this.subnetManager == nil {
		this.subnetManager = &subnet.SubnetManager{}
	}
	this.subnetManager.Logger = this.Logger
	// 初始化服务器网络管理器
	this.subnetManager.InitManager(this.Configer)

	this.Debug("[BaseModule.InitModule] module initting...")
	// gateway初始化
	if gateaddr := this.Configer.GetModuleSetting("gatetcpaddr"); gateaddr != "" {
		this.gateBase = &gate.GateBase{
			Logger: this.Logger,
		}
		this.gateBase.Init(this.GetModuleID())
		this.gateBase.BindOuterTCP(gateaddr)
	}
}

func (this *BaseModule) AfterInitModule() {
	this.Debug("[BaseModule.AfterInitModule] 模块 [%s] 初始化完成",
		this.GetModuleID())
}

func (this *BaseModule) GetConfiger() *conf.ModuleConfig {
	return this.Configer
}

func (this *BaseModule) InitSubnet(subnetAddrMap map[string]string) {
	for k, addr := range subnetAddrMap {
		if k != this.GetModuleID() {
			this.subnetManager.TryConnectServer(k, addr)
		}
	}
}

func (this *BaseModule) SendServerMsgByTmpID(
	serverTmpID string, msgstr msg.MsgStruct) {
	conn := this.subnetManager.GetServerConn(serverTmpID)
	if conn != nil {
		conn.SendCmd(this.getServerMsgPack(msgstr, conn))
	}
}

func (this *BaseModule) SendGateMsgByTmpID(fromconn *tcpconn.ClientConn,
	serverTmpID string, msgname string, data []byte) {
	conn := this.subnetManager.GetServerConn(serverTmpID)
	if conn != nil {
		conn.SendCmd(this.getGateServerMsgPack(msgname, data, fromconn, conn))
	}
}

func (this *BaseModule) BroadcastServerCmd(msgstr msg.MsgStruct) {
	this.subnetManager.BroadcastCmd(this.getServerMsgPack(msgstr, nil))
}

func (this *BaseModule) GetBalanceServerID(servertype string) string {
	return this.subnetManager.GetRandomServerConn(servertype).Tempid
}

func (this *BaseModule) getServerMsgPack(msgstr msg.MsgStruct,
	tarconn *tcpconn.ServerConn) msg.MsgStruct {
	res := &servercomm.SForwardToServer{}
	res.FromServerID = this.ModuleID
	if tarconn != nil {
		res.ToServerID = tarconn.Serverinfo.ServerID
	}
	res.MsgName = msgstr.GetMsgName()
	size := msgstr.GetSize()
	res.Data = make([]byte, size)
	msgstr.WriteBinary(res.Data)
	return res
}

func (this *BaseModule) getGateServerMsgPack(msgname string, data []byte,
	fromconn *tcpconn.ClientConn, tarconn *tcpconn.ServerConn) msg.MsgStruct {
	res := &servercomm.SForwardFromGate{}
	res.FromServerID = this.ModuleID
	if tarconn != nil {
		res.ToServerID = tarconn.Serverinfo.ServerID
	}
	if fromconn != nil {
		res.ClientConnID = fromconn.Tempid
	}
	res.MsgName = msgname
	size := len(data)
	res.Data = make([]byte, size)
	copy(res.Data, data)
	return res
}

func (this *BaseModule) GetGate() *gate.GateBase {
	return this.gateBase
}

func (this *BaseModule) GetSubnetManager() *subnet.SubnetManager {
	return this.subnetManager
}

func (this *BaseModule) GetModuleID() string {
	return this.ModuleID
}

func (this *BaseModule) TopRunner() {
	this.RegTimer(time.Minute, 0, false, func(t time.Duration) {
		this.Debug("Timer 1 Minute...")
	})
}

func (this *BaseModule) KillModule() {
	this.Debug("KillModule...")
	this.hasKilledModule = true
	this.KillRegister()

	// 退出完成
	this.hasStopped = true
	this.Logger.CloseLogger()
}

func (this *BaseModule) IsStopped() bool {
	return this.hasStopped
}
