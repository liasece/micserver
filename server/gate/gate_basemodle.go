package gate

import (
	"github.com/liasece/micserver/conf"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/server/gate/manager"
	"github.com/liasece/micserver/server/subnet"
)

type GateBase struct {
	*log.Logger

	subnetManager *subnet.SubnetManager
	modleConf     *conf.TopConfig

	clientSocketManager *manager.ClientSocketManager
}

func (this *GateBase) Init(moduleID string) {
	this.clientSocketManager = &manager.ClientSocketManager{
		Logger: this.Logger,
	}
	this.clientSocketManager.Init(moduleID)
}

func (this *GateBase) BindOuterTCP(addr string) {
	// 绑定 TCPSocket 服务
	this.clientSocketManager.StartAddClientTcpSocketHandle(addr)
}
