package server

import (
	"github.com/liasece/micserver/connect"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/server/base"
	"github.com/liasece/micserver/servercomm"
)

type serverCmdHandler struct {
	server *Server

	serverHook base.ServerHook
}

func (this *serverCmdHandler) HookServer(serverHook base.ServerHook) {
	this.serverHook = serverHook
}

func (this *serverCmdHandler) onForwardToServer(conn *connect.Server,
	smsg *servercomm.SForwardToServer) {
	if this.serverHook != nil {
		msg := &servercomm.ServerMessage{
			FromServer: conn.ServerInfo,
			MsgID:      smsg.MsgID,
			Data:       smsg.Data,
		}
		this.serverHook.OnServerMessage(msg)
	}
}

func (this *serverCmdHandler) onForwardFromGate(conn *connect.Server,
	smsg *servercomm.SForwardFromGate) {
	if this.serverHook != nil {
		msg := &servercomm.ClientMessage{
			FromServer:   conn.ServerInfo,
			ClientConnID: smsg.ClientConnID,
			Session:      smsg.Session,
			MsgID:        smsg.MsgID,
			Data:         smsg.Data,
		}
		this.serverHook.OnClientMessage(msg)
	}
}

func (this *serverCmdHandler) onForwardToClient(smsg *servercomm.SForwardToClient) {
	err := this.server.DoSendBytesToClient(smsg.FromServerID, smsg.ToGateID,
		smsg.ToClientID, smsg.MsgID, smsg.Data)
	if err != nil {
		this.server.Error("this.doSendBytesToClient Err:%s", err.Error())
	}
}

func (this *serverCmdHandler) onUpdateSession(smsg *servercomm.SUpdateSession) {
	client := this.server.GetClient(smsg.ClientConnID)
	if client != nil {
		client.Session.FromMap(smsg.Session)
		if client.Session.GetUUID() != "" {
			this.server.Info("[gate] 用户登陆成功 %s", smsg.GetJson())
		}
	} else {
		this.server.Warn("serverCmdHandler.OnUpdateSession client == nil[%s]",
			smsg.ClientConnID)
	}
}

func (this *serverCmdHandler) OnRecvSubnetMsg(conn *connect.Server,
	msgbinary *msg.MessageBinary) {
	switch msgbinary.CmdID {
	case servercomm.SForwardToServerID:
		// 服务器间用户空间消息转发
		if this.serverHook != nil {
			layerMsg := &servercomm.SForwardToServer{}
			layerMsg.ReadBinary(msgbinary.ProtoData)
			this.onForwardToServer(conn, layerMsg)
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
