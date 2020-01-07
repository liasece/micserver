/*
gateway基础模块
*/
package gate

import (
	"net"

	"github.com/liasece/micserver/conf"
	"github.com/liasece/micserver/connect"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/server/gate/base"
	"github.com/liasece/micserver/server/subnet"
	"github.com/liasece/micserver/servercomm"
)

// Gateway基础模块
type GateBase struct {
	*log.Logger

	subnetManager *subnet.SubnetManager
	modleConf     *conf.TopConfig

	gateHook base.GateHook
	connPool connect.ClientPool
}

// 初始化模块
func (this *GateBase) Init(moduleID string) {
	this.connPool.SetLogger(this.Logger)
	this.connPool.Init()
}

// 绑定外部TCP地址
func (this *GateBase) BindOuterTCP(addr string) {
	// 绑定 TCPSocket 服务
	// 由于部分 NAT 主机没有网卡概念，需要自己配置IP
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		this.Error("[GateBase.StartAddClientTcpSocketHandle] %s",
			err.Error())
		return
	}
	this.Syslog("[GateBase.StartAddClientTcpSocketHandle] "+
		"Gateway Client TCP服务启动成功 IPPort[%s]", addr)
	go func() {
		for {
			// 接受连接
			netConn, err := ln.Accept()
			if err != nil {
				// handle error
				this.Error("[GateBase.StartAddClientTcpSocketHandle] "+
					"Accept() ERR:%q",
					err.Error())
				continue
			}
			this.OnAcceptClientConnect(netConn)
			this.addTCPClient(netConn)
		}
	}()
}

// 由Client调用，当Client关闭时触发
func (this *GateBase) OnConnectClose(client *connect.Client) {
	this.remove(client.GetConnectID())

	if this.gateHook != nil {
		client.Syslog("[OnClose] 关闭Client对象")
		this.gateHook.OnCloseClient(client)
	} else {
		client.Debug("[OnNewClient] 关闭Client对象但未处理，缺少GateHook")
	}
}

// 由Client调用，当Client收到消息时
func (this *GateBase) OnRecvConnectMessage(client *connect.Client,
	msgbin *msg.MessageBinary) {
	cmdname := servercomm.MsgIdToString(msgbin.GetMsgID())
	defer msgbin.Free()

	// 检查链接是否已被断开，如果已断开则不处理
	if !client.Check() {
		client.Shutdown()
		client.Debug("[ParseClientJsonMsg] 客户端连接已关闭，丢弃消息 "+
			"MsgID[%d] MsgName[%s] Data[%s]",
			msgbin.GetMsgID(), cmdname, msgbin.String())
		return
	}
	// 接收到有效消息，开始处理
	if this.gateHook != nil {
		client.Syslog("[ParseClientJsonMsg] 收到客户端消息 "+
			"MsgID[%d] Msgname[%s] MsgLen[%d] DataLen[%d]",
			msgbin.GetMsgID(), cmdname, msgbin.GetTotalLength(),
			msgbin.GetProtoLength())
		this.gateHook.OnRecvClientMsg(client, msgbin)
	} else {
		client.Debug("[ParseClientJsonMsg] 收到客户端消息但未处理，缺少GateHook"+
			"MsgID[%d] Msgname[%s] MsgLen[%d] DataLen[%d]",
			msgbin.GetMsgID(), cmdname, msgbin.GetTotalLength(),
			msgbin.GetProtoLength())
	}
}

// 当新建一个Client对象时
func (this *GateBase) OnNewClient(client *connect.Client) {
	if this.gateHook != nil {
		client.Syslog("[OnNewClient] 创建Client对象")
		this.gateHook.OnNewClient(client)
	} else {
		client.Debug("[OnNewClient] 创建Client对象但未处理，缺少GateHook")
	}
}

// 当收到一个客户端net连接时
func (this *GateBase) OnAcceptClientConnect(conn net.Conn) {
	if this.gateHook != nil {
		this.Syslog("收到Net连接 RemoteAddr[%s]", conn.RemoteAddr().String())
		this.gateHook.OnAcceptClientConnect(conn)
	} else {
		this.Debug("收到Net连接但未处理，缺少GateHook RemoteAddr[%s]",
			conn.RemoteAddr().String())
	}
}

// 注册gate服务处理钩子
func (this *GateBase) HookGate(gateHook base.GateHook) {
	this.gateHook = gateHook
}

// 根据连接ID获取一个Client
func (this *GateBase) GetClient(connectid string) *connect.Client {
	return this.connPool.Get(connectid)
}

// 获取当前已连接的Client数量
func (this *GateBase) GetClientCount() uint32 {
	return this.connPool.Len()
}

// 删除已连接的Client
func (this *GateBase) remove(connectid string) {
	this.connPool.Remove(connectid)
}

// 遍历所有连接到本模块的客户端
func (this *GateBase) Range(
	callback func(string, *connect.Client) bool) {
	this.connPool.Range(func(value *connect.Client) bool {
		return callback(value.GetConnectID(), value)
	})
}

// 遍历所有的连接，检查需要移除的连接
func (this *GateBase) RangeRemove(
	callback func(*connect.Client) bool) {
	removelist := make([]string, 0)
	this.connPool.Range(func(value *connect.Client) bool {
		// 遍历所有的连接
		if callback(value) {
			// 该连接需要被移除
			removelist = append(removelist, value.GetConnectID())
			value.Terminate()
		}
		return true
	})
	for _, v := range removelist {
		this.remove(v)
	}

	this.Syslog("[GateBase.RangeRemove] "+
		"遍历删除连接数 RemoveSum[%d] NowLinkSum[%d]",
		len(removelist), this.GetClientCount())
}

func (this *GateBase) addTCPClient(
	netConn net.Conn) (*connect.Client, error) {
	conn, err := this.connPool.NewTCPClient(netConn, this)
	if err != nil {
		return nil, err
	}

	// 当创建一个Client对象时调用
	this.OnNewClient(conn)

	conn.Syslog("[GateBase.addTCPClient] 新增客户端连接 NowLinkSum[%d]",
		this.GetClientCount())
	// 开始接收数据
	conn.StartRecv()
	return conn, nil
}
