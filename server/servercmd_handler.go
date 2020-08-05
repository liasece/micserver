package server

import (
	"github.com/liasece/micserver/connect"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/server/base"
	"github.com/liasece/micserver/servercomm"
	"github.com/liasece/micserver/session"
)

// 服务消息处理
type serverCmdHandler struct {
	server *Server

	serverHook base.ServerHook
}

// HookServer 设置服务器消息事件的监听者
func (handler *serverCmdHandler) HookServer(serverHook base.ServerHook) {
	handler.serverHook = serverHook
}

// onForwardToModule 当需要将一个消息转发到其他服务器中时调用
func (handler *serverCmdHandler) onForwardToModule(conn *connect.Server, smsg *servercomm.SForwardToModule) {
	if handler.serverHook != nil {
		msg := &servercomm.ModuleMessage{
			FromModule: conn.ModuleInfo,
			MsgID:      smsg.MsgID,
			Data:       smsg.Data,
		}
		handler.serverHook.OnModuleMessage(msg)
	}
}

// onForwardFromGate 当收到一个从网关转发过来的消息时调用
func (handler *serverCmdHandler) onForwardFromGate(conn *connect.Server, smsg *servercomm.SForwardFromGate) {
	if handler.serverHook != nil {
		msg := &servercomm.ClientMessage{
			FromModule:   conn.ModuleInfo,
			ClientConnID: smsg.ClientConnID,
			MsgID:        smsg.MsgID,
			Data:         smsg.Data,
		}
		uuid := session.GetUUIDFromMap(smsg.Session)
		var se *session.Session
		if uuid != "" {
			se = handler.server.GetSession(uuid)
		}
		if se == nil {
			se = session.NewSessionFromMap(smsg.Session)
		}
		handler.serverHook.OnClientMessage(se, msg)
	}
}

// onForwardToClient 当收到转发一个消息到客户端时调用
func (handler *serverCmdHandler) onForwardToClient(smsg *servercomm.SForwardToClient) {
	err := handler.server.DoSendBytesToClient(smsg.FromModuleID, smsg.ToGateID,
		smsg.ToClientID, smsg.MsgID, smsg.Data)
	if err != nil {
		if err == ErrTargetClientDontExist {
			handler.server.Debug("[serverCmdHandler.onForwardToClient] ErrTargetClientDontExist", log.ErrorField(err),
				log.String("FromModuleID", smsg.FromModuleID), log.String("ToGateID", smsg.ToGateID), log.String("ToClientID", smsg.ToClientID),
				log.Uint16("MsgID", smsg.MsgID), log.ByteString("Data", smsg.Data))
		} else {
			handler.server.Error("[serverCmdHandler.onForwardToClient] error", log.ErrorField(err),
				log.String("FromModuleID", smsg.FromModuleID), log.String("ToGateID", smsg.ToGateID), log.String("ToClientID", smsg.ToClientID),
				log.Uint16("MsgID", smsg.MsgID), log.ByteString("Data", smsg.Data))
		}
	}
}

// onUpdateSession 当收到Session更新消息时调用
func (handler *serverCmdHandler) onUpdateSession(smsg *servercomm.SUpdateSession) {
	var connectedSession *session.Session
	if handler.server.gateBase != nil {
		client := handler.server.GetClient(smsg.ClientConnID)
		if client != nil {
			client.Session.FromMap(smsg.Session)
			connectedSession = client.Session
			// if client.Session.GetUUID() != "" {
			// 	handler.server.Info("[gate] 用户登陆成功 %s", smsg.GetJSON())
			// }
		}
	}

	// 尝试更新本地 session
	if smsg.SessionUUID != "" {
		// 先从连接中的session复制
		s := connectedSession
		localsession := handler.server.sessionManager.GetSession(smsg.SessionUUID)
		if localsession != nil {
			if connectedSession != nil && connectedSession != localsession {
				// 不是同一个session对象，需要将本地session复制为最新链接的session
				connectedSession.OnlyAddKeyFromSession(localsession)
				handler.server.sessionManager.Store(connectedSession)
				localsession = connectedSession
			}
			s = localsession
		}
		if s == nil {
			s = &session.Session{}
			handler.server.sessionManager.UpdateSessionUUID(smsg.SessionUUID, s)
		}
		handler.server.sessionManager.MustUpdateFromMap(s, smsg.Session)
		handler.server.Syslog("[serverCmdHandler.onUpdateSession] Session Manager Update: %+v From:%s To:%s", log.Reflect("Session", smsg.Session),
			log.String("FromModuleID", smsg.FromModuleID), log.String("ToModuleID", smsg.ToModuleID))
	}
}

// onReqCloseConnect 当收到一个关闭客户端连接的请求时调用
func (handler *serverCmdHandler) onReqCloseConnect(smsg *servercomm.SReqCloseConnect) {
	handler.server.Syslog("[serverCmdHandler.onReqCloseConnect] Request close client connect", log.String("FromModule", smsg.FromModuleID),
		log.String("ToModuleID", smsg.ToModuleID), log.String("ClientID", smsg.ClientConnID))
	handler.server.ReqCloseConnect(smsg.ToModuleID, smsg.ClientConnID)
}

// OnServerJoinSubnet 当一个服务器成功加入网络时调用
func (handler *serverCmdHandler) OnServerJoinSubnet(server *connect.Server) {
	handler.server.onServerJoinSubnet(server)
}

// OnRecvSubnetMsg 当收到一个其他服务发过来的消息时调用
func (handler *serverCmdHandler) OnRecvSubnetMsg(conn *connect.Server, msgbinary *msg.MessageBinary) {
	switch msgbinary.GetMsgID() {
	case servercomm.SForwardToModuleID:
		// 服务器间用户空间消息转发
		if handler.serverHook != nil {
			layerMsg := &servercomm.SForwardToModule{}
			layerMsg.ReadBinary(msgbinary.ProtoData)
			handler.onForwardToModule(conn, layerMsg)
		}
	case servercomm.SForwardFromGateID:
		// Gateway 转发过来的客户端消息
		layerMsg := &servercomm.SForwardFromGate{}
		layerMsg.ReadBinary(msgbinary.ProtoData)
		handler.onForwardFromGate(conn, layerMsg)
	case servercomm.SForwardToClientID:
		// 其他服务器转发过来的，要发送到客户端的消息
		var layerMsg *servercomm.SForwardToClient
		if obj := msgbinary.GetObj(); obj != nil {
			if m, ok := obj.(*servercomm.SForwardToClient); ok {
				layerMsg = m
			}
		}
		if layerMsg == nil {
			layerMsg = &servercomm.SForwardToClient{}
			layerMsg.ReadBinary(msgbinary.ProtoData)
		}
		handler.onForwardToClient(layerMsg)
	case servercomm.SUpdateSessionID:
		// 客户端会话更新
		layerMsg := &servercomm.SUpdateSession{}
		layerMsg.ReadBinary(msgbinary.ProtoData)
		handler.onUpdateSession(layerMsg)
	case servercomm.SReqCloseConnectID:
		// 关闭客户端连接
		layerMsg := &servercomm.SReqCloseConnect{}
		layerMsg.ReadBinary(msgbinary.ProtoData)
		handler.onReqCloseConnect(layerMsg)
	case servercomm.SStartMyNotifyCommandID:
	case servercomm.SROCBindID:
		// ROC 对象绑定
		layerMsg := &servercomm.SROCBind{}
		layerMsg.ReadBinary(msgbinary.ProtoData)
		handler.server.ROCServer.onMsgROCBind(layerMsg)
	case servercomm.SROCRequestID:
		// ROC 调用请求
		layerMsg := &servercomm.SROCRequest{}
		layerMsg.ReadBinary(msgbinary.ProtoData)
		handler.server.ROCServer.onMsgROCRequest(layerMsg)
	case servercomm.SROCResponseID:
		// ROC 调用返回
		layerMsg := &servercomm.SROCResponse{}
		layerMsg.ReadBinary(msgbinary.ProtoData)
		handler.server.ROCServer.onMsgROCResponse(layerMsg)
	default:
		msgid := msgbinary.GetMsgID()
		msgname := servercomm.MsgIdToString(msgid)
		handler.server.Error("[SubnetManager.OnRecvTCPMsg] Unknow message", log.Uint16("MsgID", msgid), log.String("MsgName", msgname))
	}
}
