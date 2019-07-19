package gate

import (
	"github.com/liasece/micserver/conf"
	"github.com/liasece/micserver/server/gate/handle"
	"github.com/liasece/micserver/server/gate/manager"
	"github.com/liasece/micserver/server/subnet"
)

type GateBase struct {
	subnetManager *subnet.SubnetManager
	modleConf     *conf.TopConfig

	clientTcpHandler    *handle.ClientTcpHandler
	clientHttpHandler   *handle.ClientHttpHandler
	clientSocketManager *manager.ClientSocketManager
}

func (this *GateBase) BindOuterTCP(addr string) {
	// 绑定 TCPSocket 服务
	this.clientTcpHandler.StartAddClientTcpSocketHandle(addr)
}
