package gate

import (
	"github.com/liasece/micserver/conf"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/server/gate/manager"
	"github.com/liasece/micserver/server/subnet"
)

type GateBase struct {
	*log.Logger
	*manager.ClientConnManager

	subnetManager *subnet.SubnetManager
	modleConf     *conf.TopConfig
}

func (this *GateBase) Init(moduleID string) {
	this.ClientConnManager = &manager.ClientConnManager{
		Logger: this.Logger,
	}
	this.ClientConnManager.Init(moduleID)
}

func (this *GateBase) BindOuterTCP(addr string) {
	// 绑定 TCPSocket 服务
	this.StartAddClientTcpSocketHandle(addr)
}
