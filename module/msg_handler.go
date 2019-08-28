package module

import (
	"github.com/liasece/micserver/servercomm"
)

type msgHandler struct {
	mod *BaseModule

	regForwardToServer func(msg *servercomm.SForwardToServer)
	regForwardFromGate func(msg *servercomm.SForwardFromGate)
	regForwardToClient func(msg *servercomm.SForwardToClient)
}

func (this *msgHandler) RegForwardToServer(
	cb func(msg *servercomm.SForwardToServer)) {
	this.regForwardToServer = cb
}

func (this *msgHandler) onForwardToServer(smsg *servercomm.SForwardToServer) {
	if this.regForwardToServer != nil {
		this.regForwardToServer(smsg)
	}
}

func (this *msgHandler) RegForwardFromGate(
	cb func(msg *servercomm.SForwardFromGate)) {
	this.regForwardFromGate = cb
}

func (this *msgHandler) onForwardFromGate(smsg *servercomm.SForwardFromGate) {
	if this.regForwardFromGate != nil {
		this.regForwardFromGate(smsg)
	}
}

func (this *msgHandler) RegForwardToClient(
	cb func(msg *servercomm.SForwardToClient)) {
	this.regForwardToClient = cb
}

func (this *msgHandler) onForwardToClient(smsg *servercomm.SForwardToClient) {
	err := this.mod.doSendBytesToClient(smsg.FromServerID, smsg.ToGateID,
		smsg.ToClientID, smsg.MsgID, smsg.Data)
	if err != nil {
		this.mod.Error("this.doSendBytesToClient Err:%s", err.Error())
	}
}

func (this *msgHandler) onUpdateSession(smsg *servercomm.SUpdateSession) {
	client := this.mod.GetClientConn(smsg.ClientConnID)
	if client != nil {
		for k, v := range smsg.Session {
			client.Session[k] = v
		}
		if client.Session.GetUUID() != "" {
			this.mod.Info("[gate] 用户登陆成功 %s", smsg.GetJson())
		}
	} else {
		this.mod.Warn("msgHandler.OnUpdateSession client == nil[%s]",
			smsg.ClientConnID)
	}
}
