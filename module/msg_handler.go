package module

import (
	"github.com/liasece/micserver/servercomm"
)

type msgHandler struct {
	mod *BaseModule

	fonForwardToServer func(msg *servercomm.SForwardToServer)
	fonForwardFromGate func(msg *servercomm.SForwardFromGate)
	fonForwardToClient func(msg *servercomm.SForwardToClient)
}

func (this *msgHandler) RegOnForwardToServer(
	cb func(msg *servercomm.SForwardToServer)) {
	this.fonForwardToServer = cb
}

func (this *msgHandler) onForwardToServer(smsg *servercomm.SForwardToServer) {
	if this.fonForwardToServer != nil {
		this.fonForwardToServer(smsg)
	}
}

func (this *msgHandler) RegOnForwardFromGate(
	cb func(msg *servercomm.SForwardFromGate)) {
	this.fonForwardFromGate = cb
}

func (this *msgHandler) onForwardFromGate(smsg *servercomm.SForwardFromGate) {
	if this.fonForwardFromGate != nil {
		this.fonForwardFromGate(smsg)
	}
}

func (this *msgHandler) RegOnForwardToClient(
	cb func(msg *servercomm.SForwardToClient)) {
	this.fonForwardToClient = cb
}

func (this *msgHandler) onForwardToClient(smsg *servercomm.SForwardToClient) {
	err := this.mod.doSendBytesToClient(smsg.FromServerID, smsg.ToGateID,
		smsg.ToClientID, smsg.MsgID, smsg.Data)
	if err != nil {
		this.mod.Error("this.doSendBytesToClient Err:%s", err.Error())
	}
}

func (this *msgHandler) onUpdateSession(smsg *servercomm.SUpdateSession) {
	client := this.mod.GetClient(smsg.ClientConnID)
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
