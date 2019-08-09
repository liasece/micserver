package module

import (
	"fmt"
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
		this.Logger = log.GetDefaultLogger().Clone()
		this.Logger.SetLogName(this.ModuleID)
	}
	// 申请内存
	if this.subnetManager == nil {
		this.subnetManager = &subnet.SubnetManager{}
	}
	this.subnetManager.Logger = this.Logger
	// 初始化服务器网络管理器
	this.subnetManager.InitManager(this.Configer)
	this.subnetManager.RegHandleToClientMsg(this.handleToClientMsg)

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

// 获取一个客户端连接
func (this *BaseModule) GetClientConn(tmpid string) *tcpconn.ClientConn {
	if this.gateBase != nil {
		return this.gateBase.GetTaskByTmpID(tmpid)
	}
	return nil
}

// 初始化服务器集群网络
func (this *BaseModule) InitSubnet(subnetAddrMap map[string]string) {
	for k, addr := range subnetAddrMap {
		if k != this.GetModuleID() {
			this.subnetManager.TryConnectServer(k, addr)
		}
	}
}

// 发送一个服务器消息到另一个服务器
func (this *BaseModule) SendServerCmdToServer(
	to string, msgstr msg.MsgStruct) {
	conn := this.subnetManager.GetServerConn(to)
	if conn != nil {
		conn.SendCmd(this.getServerMsgPack(msgstr, conn))
	}
}

// 转发一个客户端消息到另一个服务器
func (this *BaseModule) ForwardClientMsgToServer(fromconn *tcpconn.ClientConn,
	to string, msgname string, data []byte) {
	conn := this.subnetManager.GetServerConn(to)
	if conn != nil {
		conn.SendCmd(this.getGateServerMsgPack(msgname, data, fromconn, conn))
	}
}

// 发送一个消息到客户端
func (this *BaseModule) SendMsgToClient(gateid string,
	to string, msgstr msg.MsgStruct) error {
	sec := false
	if this.ModuleID == gateid {
		if this.doSendMsgToClient(this.ModuleID, gateid, to, msgstr) == nil {
			sec = true
		}
	} else {
		conn := this.subnetManager.GetServerConn(gateid)
		if conn != nil {
			forward := &servercomm.SForwardToClient{}
			forward.FromServerID = this.ModuleID
			forward.MsgName = msgstr.GetMsgName()
			forward.MsgID = msgstr.GetMsgId()
			forward.ToClientID = to
			forward.ToGateID = gateid
			forward.Data = make([]byte, msgstr.GetSize())
			msgstr.WriteBinary(forward.Data)
			conn.SendCmd(forward)
			sec = true
		}
	}
	if !sec {
		return fmt.Errorf("目标客户端连接不存在")
	}
	return nil
}

// 发送一个消息到连接到本服务器的客户端
func (this *BaseModule) doSendMsgToClient(fromserver string, gateid string,
	to string, msgstr msg.MsgStruct) error {
	sec := false
	gate := this.GetGate()
	if gate != nil {
		conn := gate.GetTaskByTmpID(to)
		if conn != nil {
			if fromserver != gateid {
				conn.Session[util.GetServerIDType(fromserver)] = fromserver
			}
			conn.SendCmd(msgstr)
			sec = true
		}
	}
	if !sec {
		return fmt.Errorf("目标客户端连接不存在")
	}
	return nil
}

// 发送一个消息到客户端
func (this *BaseModule) SendBytesToClient(gateid string,
	to string, msgid uint16, msgname string, data []byte) error {
	sec := false
	if this.ModuleID == gateid {
		if this.doSendBytesToClient(
			this.ModuleID, gateid, to, msgid, msgname, data) == nil {
			sec = true
		}
	} else {
		conn := this.subnetManager.GetServerConn(gateid)
		if conn != nil {
			forward := &servercomm.SForwardToClient{}
			forward.FromServerID = this.ModuleID
			forward.MsgName = msgname
			forward.MsgID = msgid
			forward.ToClientID = to
			forward.ToGateID = gateid
			forward.Data = make([]byte, len(data))
			copy(forward.Data, data)
			conn.SendCmd(forward)
			sec = true
		}
	}
	if !sec {
		return fmt.Errorf("目标客户端连接不存在")
	}
	return nil
}

// 发送一个消息到连接到本服务器的客户端
func (this *BaseModule) doSendBytesToClient(fromserver string, gateid string,
	to string, msgid uint16, msgname string, data []byte) error {
	sec := false
	gate := this.GetGate()
	if gate != nil {
		conn := gate.GetTaskByTmpID(to)
		if conn != nil {
			if fromserver != gateid {
				conn.Session[util.GetServerIDType(fromserver)] = fromserver
			}
			conn.SendBytes(msgid, data)
			sec = true
		}
	}
	if !sec {
		return fmt.Errorf("目标客户端连接不存在")
	}
	return nil
}

// 广播一个消息到连接到本服务器的所有服务器
func (this *BaseModule) BroadcastServerCmd(msgstr msg.MsgStruct) {
	this.subnetManager.BroadcastCmd(this.getServerMsgPack(msgstr, nil))
}

// 获取一个均衡的负载服务器
func (this *BaseModule) GetBalanceServerID(servertype string) string {
	server := this.subnetManager.GetRandomServerConn(servertype)
	if server != nil {
		return server.Tempid
	}
	return ""
}

// 获取一个服务器消息的服务器间转发协议
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

// 获取一个客户端消息到其他服务器间的转发协议
func (this *BaseModule) getGateServerMsgPack(msgname string, data []byte,
	fromconn *tcpconn.ClientConn, tarconn *tcpconn.ServerConn) msg.MsgStruct {
	res := &servercomm.SForwardFromGate{}
	res.FromServerID = this.ModuleID
	if tarconn != nil {
		res.ToServerID = tarconn.Serverinfo.ServerID
	}
	if fromconn != nil {
		res.ClientConnID = fromconn.Tempid
		res.Session = fromconn.Session
	}
	res.MsgName = msgname
	size := len(data)
	res.Data = make([]byte, size)
	copy(res.Data, data)
	return res
}

// 获取本服务器的gate管理器
func (this *BaseModule) GetGate() *gate.GateBase {
	return this.gateBase
}

// 获取本服务器的集群子网管理器
func (this *BaseModule) GetSubnetManager() *subnet.SubnetManager {
	return this.subnetManager
}

func (this *BaseModule) GetModuleID() string {
	return this.ModuleID
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

func (this *BaseModule) TopRunner() {
	this.RegTimer(time.Minute, 0, false, func(t time.Duration) {
		this.Debug("Timer 1 Minute...")
	})
}

// 处理经过本服务器发送到客户端的消息
func (this *BaseModule) handleToClientMsg(smsg *servercomm.SForwardToClient) {
	this.doSendBytesToClient(smsg.FromServerID, smsg.ToGateID, smsg.ToClientID,
		smsg.MsgID, smsg.MsgName, smsg.Data)
}
