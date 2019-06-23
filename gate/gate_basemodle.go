package gate

import (
	"github.com/liasece/micserver/conf"
	"github.com/liasece/micserver/gate/handle"
	"github.com/liasece/micserver/gate/manager"
	"github.com/liasece/micserver/server/subnet"
)

type GateBase struct {
	subnetManager *subnet.SubnetManager
	modleConf     *conf.ServerConfig

	clientTcpHandler    *handle.ClientTcpHandler
	clientHttpHandler   *handle.ClientHttpHandler
	clientSocketManager *manager.ClientSocketManager
}

func (this *GateBase) Init(config *conf.ServerConfig) {
	this.modleConf = config
	this.subnetManager = &subnet.SubnetManager{}
	// 初始化服务器网络管理器
	this.subnetManager.InitManager()

	this.clientSocketManager = &manager.ClientSocketManager{}
	this.clientSocketManager.Init(this.subnetManager, 3)

	this.clientTcpHandler = &handle.ClientTcpHandler{}
	this.clientTcpHandler.Init(this.clientSocketManager, this.modleConf)

	this.clientHttpHandler = &handle.ClientHttpHandler{}
}

func (this *GateBase) BindGMHttp() {
	// 绑定 GM Http 服务
	gmhttpport := this.modleConf.GetProp("gmhttpport")
	if gmhttpport != "" {
		this.clientHttpHandler.StartAddHttpHandle(":" + gmhttpport)
	}
}

func (this *GateBase) BindOuterTCP() {
	tcpport := this.modleConf.GetPropUint("tcpouterport")
	// 绑定 TCPSocket 服务
	this.clientTcpHandler.StartAddClientTcpSocketHandle("", tcpport)
}
