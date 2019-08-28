package gate

import (
	"github.com/liasece/micserver/conf"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/server/gate/manager"
	"github.com/liasece/micserver/server/subnet"
)

type GateBase struct {
	*log.Logger
	*manager.ClientManager

	subnetManager *subnet.SubnetManager
	modleConf     *conf.TopConfig
}

func (this *GateBase) Init(moduleID string) {
	this.ClientManager = &manager.ClientManager{
		Logger: this.Logger,
	}
	this.ClientManager.Init(moduleID)
}

func (this *GateBase) BindOuterTCP(addr string) {
	// 绑定 TCPSocket 服务
	this.StartAddClientTcpSocketHandle(addr)
}
