package server

import (
	"github.com/liasece/micserver/servercomm"
)

type serverCmdHandler struct {
	server *Server

	fonForwardToServer func(msg *servercomm.SForwardToServer)
	fonForwardFromGate func(msg *servercomm.SForwardFromGate)
	fonForwardToClient func(msg *servercomm.SForwardToClient)
}

func (this *serverCmdHandler) RegOnForwardToServer(
	cb func(msg *servercomm.SForwardToServer)) {
	this.fonForwardToServer = cb
}

func (this *serverCmdHandler) onForwardToServer(smsg *servercomm.SForwardToServer) {
	if this.fonForwardToServer != nil {
		this.fonForwardToServer(smsg)
	}
}

func (this *serverCmdHandler) RegOnForwardFromGate(
	cb func(msg *servercomm.SForwardFromGate)) {
	this.fonForwardFromGate = cb
}

func (this *serverCmdHandler) onForwardFromGate(smsg *servercomm.SForwardFromGate) {
	if this.fonForwardFromGate != nil {
		this.fonForwardFromGate(smsg)
	}
}

func (this *serverCmdHandler) RegOnForwardToClient(
	cb func(msg *servercomm.SForwardToClient)) {
	this.fonForwardToClient = cb
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
