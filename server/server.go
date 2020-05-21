/*
Package server micserver中管理与其他服务器连接的管理器
*/
package server

import (
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

// Server 一个Module就是一个Server
type Server struct {
	*log.Logger
	// event libs
	ROCServer

	serverCmdHandler   serverCmdHandler
	clientEventHandler clientEventHandler
	subnetManager      *subnet.Manager
	gateBase           *gate.Base
	sessionManager     session.Manager

	// server info
	moduleid     string
	moduleConfig *conf.ModuleConfig
	isStop       bool
	stopChan     chan bool
}

// Init 初始化本服务
func (s *Server) Init(moduleid string) {
	s.moduleid = moduleid
	s.stopChan = make(chan bool)
	s.ROCServer.Init(s)
}

// InitSubnet 初始化本服务的子网管理器
func (s *Server) InitSubnet(conf *conf.ModuleConfig) {
	s.moduleConfig = conf
	// 初始化服务器网络管理器
	if s.subnetManager == nil {
		s.subnetManager = &subnet.Manager{}
	}
	s.serverCmdHandler.server = s
	s.subnetManager.Logger = s.Logger.Clone()
	s.subnetManager.Init(conf)
	s.subnetManager.HookSubnet(&s.serverCmdHandler)
}

// HookServer 设置本服务的服务事件监听者
func (s *Server) HookServer(serverHook serverbase.ServerHook) {
	s.serverCmdHandler.HookServer(serverHook)
}

// HookGate 设置本服务的网关事件监听者，如果本服务没有启用网关，将不会收到任何事件
func (s *Server) HookGate(gateHook gatebase.GateHook) {
	s.clientEventHandler.HookGate(gateHook)
}

// BindSubnet 尝试连接本服务子网中的其他服务器
func (s *Server) BindSubnet(subnetAddrMap map[string]string) {
	for k, addr := range subnetAddrMap {
		if k != s.moduleid {
			s.subnetManager.TryConnectServer(k, addr)
		}
	}
}

// InitGate 初始化本服务的网关部分
func (s *Server) InitGate(gateaddr string) {
	s.gateBase = &gate.Base{
		Logger: s.Logger,
	}
	s.clientEventHandler.server = s
	s.gateBase.Init(s.moduleid)
	s.gateBase.BindOuterTCP(gateaddr)

	// 事件监听
	s.gateBase.HookGate(&s.clientEventHandler)
}

// SetLogger 设置本服务的Logger
func (s *Server) SetLogger(source *log.Logger) {
	if source == nil {
		s.Logger = nil
		return
	}
	s.Logger = source.Clone()
}

// GetClient 获取一个客户端连接
func (s *Server) GetClient(tmpid string) *connect.Client {
	if s.gateBase != nil {
		return s.gateBase.GetClient(tmpid)
	}
	return nil
}

// RangeClient 获取一个客户端连接
func (s *Server) RangeClient(
	f func(tmpid string, client *connect.Client) bool) {
	if s.gateBase != nil {
		s.gateBase.Range(f)
	}
}

// onServerJoinSubnet 当一个服务器成功加入网络时调用
func (s *Server) onServerJoinSubnet(server *connect.Server) {
	s.Debug("服务器 ModuleID[%s] 加入子网成功",
		server.ModuleInfo.ModuleID)
	s.ROCServer.onServerJoinSubnet(server)
}

// SendModuleMsg 发送一个服务器消息到另一个服务器
func (s *Server) SendModuleMsg(
	to string, msgstr msg.IMsgStruct) {
	conn := s.subnetManager.GetServer(to)
	if conn != nil {
		conn.SendCmd(s.getModuleMsgPack(msgstr, conn))
	}
}

// SInnerCloseSessionConnect 断开一个客户端连接,仅框架内使用
func (s *Server) SInnerCloseSessionConnect(gateid string, connectid string) {
	s.ReqCloseConnect(gateid, connectid)
}

// ReqCloseConnect 请求关闭远程瞪的目标客户端连接
func (s *Server) ReqCloseConnect(gateid string, connectid string) {
	if s.moduleid == gateid {
		s.doCloseConnect(connectid)
	} else {
		// 向gate请求
		conn := s.subnetManager.GetServer(gateid)
		if conn != nil {
			msg := &servercomm.SReqCloseConnect{
				FromModuleID: s.moduleid,
				ToModuleID:   gateid,
				ClientConnID: connectid,
			}
			conn.SendCmd(msg)
		} else {
			s.Error("Server.ReqCloseConnect "+
				"target module does not exist GateID[%s]",
				gateid)
		}
	}
}

// doCloseConnect 关闭本地的目标客户端连接
func (s *Server) doCloseConnect(connectid string) {
	if s.gateBase == nil {
		s.Error("Server.doCloseConnect s module isn't gate")
		return
	}
	client := s.gateBase.GetClient(connectid)
	if client == nil {
		s.Error("Server.doCloseConnect client does not exist ConnectID[%s]",
			connectid)
		return
	}
	client.Terminate()
}

// SInnerSendModuleMsg 发送一个服务器消息到另一个服务器,仅框架内使用
func (s *Server) SInnerSendModuleMsg(
	to string, msgstr msg.IMsgStruct) {
	conn := s.subnetManager.GetServer(to)
	if conn != nil {
		conn.SendCmd(msgstr)
	} else {
		s.Error("Server.SInnerSendServerMsg conn == nil[%s]", to)
	}
}

// SInnerSendClientMsg 发送一个服务器消息到另一个服务器,仅框架内使用
func (s *Server) SInnerSendClientMsg(
	gateid string, connectid string, msgid uint16, data []byte) {
	s.SendBytesToClient(gateid, connectid, msgid, data)
}

// ForwardClientMsgToModule 转发一个客户端消息到另一个服务器
func (s *Server) ForwardClientMsgToModule(fromconn *connect.Client,
	to string, msgid uint16, data []byte) {
	conn := s.subnetManager.GetServer(to)
	if conn != nil {
		conn.SendCmd(s.getFarwardFromGateMsgPack(msgid, data, fromconn, conn))
	} else {
		s.Error("Server.ForwardClientMsgToServer conn == nil [%s]",
			to)
	}
}

// BroadcastModuleCmd 广播一个消息到连接到本服务器的所有服务器
func (s *Server) BroadcastModuleCmd(msgstr msg.IMsgStruct) {
	s.subnetManager.BroadcastCmd(s.getModuleMsgPack(msgstr, nil))
}

// GetBalanceModuleID 获取一个均衡的负载服务器
func (s *Server) GetBalanceModuleID(moduletype string) string {
	server := s.subnetManager.GetRandomServer(moduletype)
	if server != nil {
		return server.GetTempID()
	}
	return ""
}

// DeleteSession 删除本地维护的 session
func (s *Server) DeleteSession(uuid string) {
	s.sessionManager.DeleteSession(uuid)
}

// GetSession 获取本地维护的 session
func (s *Server) GetSession(uuid string) *session.Session {
	return s.sessionManager.GetSession(uuid)
}

// MustUpdateSessionFromMap 更新本地的Session，如果没有的话注册它
func (s *Server) MustUpdateSessionFromMap(uuid string, data map[string]string) {
	se := s.server.sessionManager.GetSession(uuid)
	if se == nil {
		se = &session.Session{}
		s.server.sessionManager.UpdateSessionUUID(uuid, se)
	}
	s.server.sessionManager.MustUpdateFromMap(se, data)
	s.server.Syslog("[MustUpdateSessionFromMap] Session Manager Update: %+v",
		data)
}

// UpdateSessionUUID 更新目标Session的UUID
func (s *Server) UpdateSessionUUID(uuid string, session *session.Session) {
	s.server.sessionManager.UpdateSessionUUID(uuid, session)
}

// SendBytesToClient 发送一个消息到客户端
func (s *Server) SendBytesToClient(gateid string,
	to string, msgid uint16, data []byte) error {
	sec := false
	if s.moduleid == gateid {
		if s.DoSendBytesToClient(
			s.moduleid, gateid, to, msgid, data) == nil {
			sec = true
		}
	} else {
		conn := s.subnetManager.GetServer(gateid)
		if conn != nil {
			forward := &servercomm.SForwardToClient{}
			forward.FromModuleID = s.moduleid
			forward.MsgID = msgid
			forward.ToClientID = to
			forward.ToGateID = gateid
			forward.Data = make([]byte, len(data))
			copy(forward.Data, data)
			conn.SendCmd(forward)
			sec = true
		} else {
			s.Error("目标服务器连接不存在 GateID[%s]", gateid)
		}
	}
	if !sec {
		return ErrTargetClientDontExist
	}
	return nil
}

// DoSendBytesToClient 发送一个消息到连接到本服务器的客户端
func (s *Server) DoSendBytesToClient(fromserver string, gateid string,
	to string, msgid uint16, data []byte) error {
	sec := false
	if s.gateBase != nil {
		conn := s.gateBase.GetClient(to)
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
		return ErrTargetClientDontExist
	}
	return nil
}

// getModuleMsgPack 获取一个服务器消息的服务器间转发协议
func (s *Server) getModuleMsgPack(msgstr msg.IMsgStruct,
	tarconn *connect.Server) msg.IMsgStruct {
	res := &servercomm.SForwardToModule{}
	res.FromModuleID = s.moduleid
	if tarconn != nil {
		res.ToModuleID = tarconn.ModuleInfo.ModuleID
	}
	res.MsgID = msgstr.GetMsgId()
	size := msgstr.GetSize()
	res.Data = make([]byte, size)
	msgstr.WriteBinary(res.Data)
	return res
}

// getFarwardFromGateMsgPack 获取一个客户端消息到其他服务器间的转发协议
func (s *Server) getFarwardFromGateMsgPack(msgid uint16, data []byte,
	fromconn *connect.Client, tarconn *connect.Server) msg.IMsgStruct {
	res := &servercomm.SForwardFromGate{}
	res.FromModuleID = s.moduleid
	if tarconn != nil {
		res.ToModuleID = tarconn.ModuleInfo.ModuleID
	}
	if fromconn != nil {
		res.ClientConnID = fromconn.GetConnectID()
		res.Session = fromconn.ToMap()
	}
	res.MsgID = msgid
	size := len(data)
	res.Data = make([]byte, size)
	copy(res.Data, data)
	return res
}

// Stop stop server
func (s *Server) Stop() {
	s.isStop = true
}
