package server

import (
	"github.com/liasece/micserver/connect"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/server/base"
	"github.com/liasece/micserver/servercomm"
	"github.com/liasece/micserver/session"
)

type serverCmdHandler struct {
	server *Server

	serverHook base.ServerHook
}

func (this *serverCmdHandler) HookServer(serverHook base.ServerHook) {
	this.serverHook = serverHook
}

func (this *serverCmdHandler) onForwardToModule(conn *connect.Server,
	smsg *servercomm.SForwardToModule) {
	if this.serverHook != nil {
		msg := &servercomm.ModuleMessage{
			FromModule: conn.ModuleInfo,
			MsgID:      smsg.MsgID,
			Data:       smsg.Data,
		}
		this.serverHook.OnModuleMessage(msg)
	}
}

func (this *serverCmdHandler) onForwardFromGate(conn *connect.Server,
	smsg *servercomm.SForwardFromGate) {
	if this.serverHook != nil {
		msg := &servercomm.ClientMessage{
			FromModule:   conn.ModuleInfo,
			ClientConnID: smsg.ClientConnID,
			Session:      smsg.Session,
			MsgID:        smsg.MsgID,
			Data:         smsg.Data,
		}
		this.serverHook.OnClientMessage(msg)
	}
}

func (this *serverCmdHandler) onForwardToClient(smsg *servercomm.SForwardToClient) {
	err := this.server.DoSendBytesToClient(smsg.FromModuleID, smsg.ToGateID,
		smsg.ToClientID, smsg.MsgID, smsg.Data)
	if err != nil {
		this.server.Error("this.doSendBytesToClient Err:%s", err.Error())
	}
}

func (this *serverCmdHandler) onUpdateSession(smsg *servercomm.SUpdateSession) {
	var connectedSession *session.Session
	if this.server.gateBase != nil {
		client := this.server.GetClient(smsg.ClientConnID)
		if client != nil {
			client.Session.FromMap(smsg.Session)
			connectedSession = client.Session
			// if client.Session.GetUUID() != "" {
			// 	this.server.Info("[gate] 用户登陆成功 %s", smsg.GetJson())
			// }
		}
	}

	// 尝试更新本地 session
	if smsg.SessionUUID != "" {
		s := connectedSession
		if s == nil {
			s = this.server.sessionManager.GetSession(smsg.SessionUUID)
			if s == nil {
				s = &session.Session{}
				s.SetUUID(smsg.SessionUUID)
			}
		}
		this.server.sessionManager.MustUpdateFromMap(s, smsg.Session)
		this.server.Debug("Session Manager Update: %+v", smsg.Session)
	}
}

func (this *serverCmdHandler) OnRecvSubnetMsg(conn *connect.Server,
	msgbinary *msg.MessageBinary) {
	switch msgbinary.CmdID {
	case servercomm.SForwardToModuleID:
		// 服务器间用户空间消息转发
		if this.serverHook != nil {
			layerMsg := &servercomm.SForwardToModule{}
			layerMsg.ReadBinary(msgbinary.ProtoData)
			this.onForwardToModule(conn, layerMsg)
		}
	case servercomm.SForwardFromGateID:
		var layerMsg *servercomm.SForwardFromGate
		if obj := msgbinary.GetObj(); obj != nil {
			if m, ok := obj.(*servercomm.SForwardFromGate); ok {
				layerMsg = m
			}
		}
		if layerMsg == nil {
			layerMsg = &servercomm.SForwardFromGate{}
			layerMsg.ReadBinary(msgbinary.ProtoData)
		}
		// Gateway 转发过来的客户端消息
		this.onForwardFromGate(conn, layerMsg)
	case servercomm.SForwardToClientID:
		// 其他服务器转发过来的，要发送到客户端的消息
		layerMsg := &servercomm.SForwardToClient{}
		layerMsg.ReadBinary(msgbinary.ProtoData)
		this.onForwardToClient(layerMsg)
	case servercomm.SUpdateSessionID:
		// 客户端会话更新
		layerMsg := &servercomm.SUpdateSession{}
		layerMsg.ReadBinary(msgbinary.ProtoData)
		this.onUpdateSession(layerMsg)
	case servercomm.SStartMyNotifyCommandID:
	case servercomm.SROCBindID:
		layerMsg := &servercomm.SROCBind{}
		layerMsg.ReadBinary(msgbinary.ProtoData)
		this.server.ROCServer.onMsgROCBind(layerMsg)
	case servercomm.SROCRequestID:
		layerMsg := &servercomm.SROCRequest{}
		layerMsg.ReadBinary(msgbinary.ProtoData)
		this.server.ROCServer.onMsgROCRequest(layerMsg)
	default:
		msgid := msgbinary.CmdID
		msgname := servercomm.MsgIdToString(msgid)
		this.server.Error("[SubnetManager.OnRecvTCPMsg] 未知消息 %d:%s",
			msgid, msgname)
	}
}
