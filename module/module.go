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

func (this *BaseModule) BroadcastServerCmd(msgstr msg.MsgStruct) {
	this.subnetManager.BroadcastCmd(this.getServerMsgPack(msgstr, nil))
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
}

func (this *BaseModule) IsStopped() bool {
	return this.hasStopped
}
