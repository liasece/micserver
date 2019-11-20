package server

import (
	"fmt"
	"github.com/liasece/micserver/conf"
	"github.com/liasece/micserver/connect"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/msg"
	serverbase "github.com/liasece/micserver/server/base"
	"github.com/liasece/micserver/server/gate"
	gatebase "github.com/liasece/micserver/server/gate/base"
	"github.com/liasece/micserver/server/subnet"
	"github.com/liasece/micserver/servercomm"
	"github.com/liasece/micserver/session"
	"github.com/liasece/micserver/util"
)

type Server struct {
	*log.Logger
	// event libs
	ROCServer

	serverCmdHandler   serverCmdHandler
	clientEventHandler clientEventHandler
	subnetManager      *subnet.SubnetManager
	gateBase           *gate.GateBase
	sessionManager     session.SessionManager

	// server info
	moduleid string
	isStop   bool
	stopChan chan bool
}

func (this *Server) Init(moduleid string) {
	this.moduleid = moduleid
	this.stopChan = make(chan bool)
	this.ROCServer.Init(this)
}

func (this *Server) InitSubnet(conf *conf.ModuleConfig) {
	// 初始化服务器网络管理器
	if this.subnetManager == nil {
		this.subnetManager = &subnet.SubnetManager{}
	}
	this.serverCmdHandler.server = this
	this.subnetManager.Logger = this.Logger.Clone()
	this.subnetManager.Init(conf)
	this.subnetManager.HookSubnet(&this.serverCmdHandler)
}

func (this *Server) HookServer(serverHook serverbase.ServerHook) {
	this.serverCmdHandler.HookServer(serverHook)
}

func (this *Server) HookGate(gateHook gatebase.GateHook) {
	this.clientEventHandler.HookGate(gateHook)
}

func (this *Server) BindSubnet(subnetAddrMap map[string]string) {
	for k, addr := range subnetAddrMap {
		if k != this.moduleid {
			this.subnetManager.TryConnectServer(k, addr)
		}
	}
}

func (this *Server) InitGate(gateaddr string) {
	this.gateBase = &gate.GateBase{
		Logger: this.Logger,
	}
	this.clientEventHandler.server = this
	this.gateBase.Init(this.moduleid)
	this.gateBase.BindOuterTCP(gateaddr)

	// 事件监听
	this.gateBase.HookGate(&this.clientEventHandler)
}

func (this *Server) SetLogger(source *log.Logger) {
	if source == nil {
		this.Logger = nil
		return
	}
	this.Logger = source.Clone()
}

// 获取一个客户端连接
func (this *Server) GetClient(tmpid string) *connect.Client {
	if this.gateBase != nil {
		return this.gateBase.GetClient(tmpid)
	}
	return nil
}

// 发送一个服务器消息到另一个服务器
func (this *Server) SendModuleMsg(
	to string, msgstr msg.MsgStruct) {
	conn := this.subnetManager.GetServer(to)
	if conn != nil {
		conn.SendCmd(this.getModuleMsgPack(msgstr, conn))
	}
}

// 发送一个服务器消息到另一个服务器,仅框架内使用
func (this *Server) SInner_SendModuleMsg(
	to string, msgstr msg.MsgStruct) {
	conn := this.subnetManager.GetServer(to)
	if conn != nil {
		conn.SendCmd(msgstr)
	} else {
		this.Error("Server.SInner_SendServerMsg conn == nil[%s]", to)
	}
}

// 发送一个服务器消息到另一个服务器,仅框架内使用
func (this *Server) SInner_SendClientMsg(
	gateid string, connectid string, msgid uint16, data []byte) {
	this.SendBytesToClient(gateid, connectid, msgid, data)
}

// 转发一个客户端消息到另一个服务器
func (this *Server) ForwardClientMsgToModule(fromconn *connect.Client,
	to string, msgid uint16, data []byte) {
	conn := this.subnetManager.GetServer(to)
	if conn != nil {
		conn.SendCmd(this.getFarwardFromGateMsgPack(msgid, data, fromconn, conn))
	} else {
		this.Error("Server.ForwardClientMsgToServer conn == nil [%s]",
			to)
	}
}

// 广播一个消息到连接到本服务器的所有服务器
func (this *Server) BroadcastModuleCmd(msgstr msg.MsgStruct) {
	this.subnetManager.BroadcastCmd(this.getModuleMsgPack(msgstr, nil))
}

// 获取一个均衡的负载服务器
func (this *Server) GetBalanceModuleID(moduletype string) string {
	server := this.subnetManager.GetRandomServer(moduletype)
	if server != nil {
		return server.GetTempID()
	}
	return ""
}

// 删除本地维护的 session
func (this *Server) DeleteSession(uuid string) {
	this.sessionManager.DeleteSession(uuid)
}

// 获取本地维护的 session
func (this *Server) GetSession(uuid string) *session.Session {
	return this.sessionManager.GetSession(uuid)
}

// 发送一个消息到客户端
func (this *Server) SendBytesToClient(gateid string,
	to string, msgid uint16, data []byte) error {
	sec := false
	if this.moduleid == gateid {
		if this.DoSendBytesToClient(
			this.moduleid, gateid, to, msgid, data) == nil {
			sec = true
		}
	} else {
		conn := this.subnetManager.GetServer(gateid)
		if conn != nil {
			forward := &servercomm.SForwardToClient{}
			forward.FromModuleID = this.moduleid
			forward.MsgID = msgid
			forward.ToClientID = to
			forward.ToGateID = gateid
			forward.Data = make([]byte, len(data))
			copy(forward.Data, data)
			conn.SendCmd(forward)
			sec = true
		} else {
			this.Error("目标服务器连接不存在 GateID[%s]", gateid)
		}
	}
	if !sec {
		return fmt.Errorf("目标客户端连接不存在")
	}
	return nil
}

// 发送一个消息到连接到本服务器的客户端
func (this *Server) DoSendBytesToClient(fromserver string, gateid string,
	to string, msgid uint16, data []byte) error {
	sec := false
	if this.gateBase != nil {
		conn := this.gateBase.GetClient(to)
		if conn != nil {
			if fromserver != gateid {
				conn.Session.SetBind(util.GetModuleIDType(fromserver),
					fromserver)
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

// 获取一个服务器消息的服务器间转发协议
func (this *Server) getModuleMsgPack(msgstr msg.MsgStruct,
	tarconn *connect.Server) msg.MsgStruct {
	res := &servercomm.SForwardToModule{}
	res.FromModuleID = this.moduleid
	if tarconn != nil {
		res.ToModuleID = tarconn.ModuleInfo.ModuleID
	}
	res.MsgID = msgstr.GetMsgId()
	size := msgstr.GetSize()
	res.Data = make([]byte, size)
	msgstr.WriteBinary(res.Data)
	return res
}

// 获取一个客户端消息到其他服务器间的转发协议
func (this *Server) getFarwardFromGateMsgPack(msgid uint16, data []byte,
	fromconn *connect.Client, tarconn *connect.Server) msg.MsgStruct {
	res := &servercomm.SForwardFromGate{}
	res.FromModuleID = this.moduleid
	if tarconn != nil {
		res.ToModuleID = tarconn.ModuleInfo.ModuleID
	}
	if fromconn != nil {
		res.Session = make(map[string]string)
		res.ClientConnID = fromconn.GetConnectID()
		res.Session = fromconn.ToMap()
	}
	res.MsgID = msgid
	size := len(data)
	res.Data = make([]byte, size)
	copy(res.Data, data)
	return res
}

func (this *Server) Stop() {
	this.isStop = true
}
