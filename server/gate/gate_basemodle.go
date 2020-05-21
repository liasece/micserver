/*
Package gate gateway基础模块
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

// Base Gateway基础模块
type Base struct {
	*log.Logger

	subnetManager *subnet.Manager
	modleConf     *conf.TopConfig

	gateHook base.GateHook
	connPool connect.ClientPool
}

// Init 初始化模块
func (gateBase *Base) Init(moduleID string) {
	gateBase.connPool.SetLogger(gateBase.Logger)
	gateBase.connPool.Init()
}

// BindOuterTCP 绑定外部TCP地址
func (gateBase *Base) BindOuterTCP(addr string) {
	// 绑定 TCPSocket 服务
	// 由于部分 NAT 主机没有网卡概念，需要自己配置IP
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		gateBase.Error("[Base.StartAddClientTcpSocketHandle] %s",
			err.Error())
		return
	}
	gateBase.Syslog("[Base.StartAddClientTcpSocketHandle] "+
		"Gateway Client TCP服务启动成功 IPPort[%s]", addr)
	go func() {
		for {
			// 接受连接
			netConn, err := ln.Accept()
			if err != nil {
				// handle error
				gateBase.Error("[Base.StartAddClientTcpSocketHandle] "+
					"Accept() ERR:%q",
					err.Error())
				continue
			}
			gateBase.OnAcceptClientConnect(netConn)
			gateBase.addTCPClient(netConn)
		}
	}()
}

// OnConnectClose 由Client调用，当Client关闭时触发
func (gateBase *Base) OnConnectClose(client *connect.Client) {
	gateBase.remove(client.GetConnectID())

	if gateBase.gateHook != nil {
		client.Syslog("[OnClose] 关闭Client对象")
		gateBase.gateHook.OnCloseClient(client)
	} else {
		client.Debug("[OnNewClient] 关闭Client对象但未处理，缺少GateHook")
	}
}

// OnRecvConnectMessage 由Client调用，当Client收到消息时
func (gateBase *Base) OnRecvConnectMessage(client *connect.Client,
	msgbin *msg.MessageBinary) {
	cmdname := servercomm.MsgIdToString(msgbin.GetMsgID())
	defer msgbin.Free()

	// 检查链接是否已被断开，如果已断开则不处理
	if !client.Check() {
		client.Shutdown()
		client.Debug("[ParseClientJSONMsg] 客户端连接已关闭，丢弃消息 "+
			"MsgID[%d] MsgName[%s] Data[%s]",
			msgbin.GetMsgID(), cmdname, msgbin.String())
		return
	}
	// 接收到有效消息，开始处理
	if gateBase.gateHook != nil {
		client.Syslog("[ParseClientJSONMsg] 收到客户端消息 "+
			"MsgID[%d] Msgname[%s] MsgLen[%d] DataLen[%d]",
			msgbin.GetMsgID(), cmdname, msgbin.GetTotalLength(),
			msgbin.GetProtoLength())
		gateBase.gateHook.OnRecvClientMsg(client, msgbin)
	} else {
		client.Debug("[ParseClientJSONMsg] 收到客户端消息但未处理，缺少GateHook"+
			"MsgID[%d] Msgname[%s] MsgLen[%d] DataLen[%d]",
			msgbin.GetMsgID(), cmdname, msgbin.GetTotalLength(),
			msgbin.GetProtoLength())
	}
}

// OnNewClient 当新建一个Client对象时
func (gateBase *Base) OnNewClient(client *connect.Client) {
	if gateBase.gateHook != nil {
		client.Syslog("[OnNewClient] 创建Client对象")
		gateBase.gateHook.OnNewClient(client)
	} else {
		client.Debug("[OnNewClient] 创建Client对象但未处理，缺少GateHook")
	}
}

// OnAcceptClientConnect 当收到一个客户端net连接时
func (gateBase *Base) OnAcceptClientConnect(conn net.Conn) {
	if gateBase.gateHook != nil {
		gateBase.Syslog("收到Net连接 RemoteAddr[%s]", conn.RemoteAddr().String())
		gateBase.gateHook.OnAcceptClientConnect(conn)
	} else {
		gateBase.Debug("收到Net连接但未处理，缺少GateHook RemoteAddr[%s]",
			conn.RemoteAddr().String())
	}
}

// HookGate 注册gate服务处理钩子
func (gateBase *Base) HookGate(gateHook base.GateHook) {
	gateBase.gateHook = gateHook
}

// GetClient 根据连接ID获取一个Client
func (gateBase *Base) GetClient(connectid string) *connect.Client {
	return gateBase.connPool.Get(connectid)
}

// GetClientCount 获取当前已连接的Client数量
func (gateBase *Base) GetClientCount() uint32 {
	return gateBase.connPool.Len()
}

// remove 删除已连接的Client
func (gateBase *Base) remove(connectid string) {
	gateBase.connPool.Remove(connectid)
}

// Range 遍历所有连接到本模块的客户端
func (gateBase *Base) Range(
	callback func(string, *connect.Client) bool) {
	gateBase.connPool.Range(func(value *connect.Client) bool {
		return callback(value.GetConnectID(), value)
	})
}

// RangeRemove 遍历所有的连接，检查需要移除的连接
func (gateBase *Base) RangeRemove(
	callback func(*connect.Client) bool) {
	removelist := make([]string, 0)
	gateBase.connPool.Range(func(value *connect.Client) bool {
		// 遍历所有的连接
		if callback(value) {
			// 该连接需要被移除
			removelist = append(removelist, value.GetConnectID())
			value.Terminate()
		}
		return true
	})
	for _, v := range removelist {
		gateBase.remove(v)
	}

	gateBase.Syslog("[Base.RangeRemove] "+
		"遍历删除连接数 RemoveSum[%d] NowLinkSum[%d]",
		len(removelist), gateBase.GetClientCount())
}

func (gateBase *Base) addTCPClient(
	netConn net.Conn) (*connect.Client, error) {
	conn, err := gateBase.connPool.NewTCPClient(netConn, gateBase)
	if err != nil {
		return nil, err
	}

	// 当创建一个Client对象时调用
	gateBase.OnNewClient(conn)

	conn.Syslog("[Base.addTCPClient] 新增客户端连接 NowLinkSum[%d]",
		gateBase.GetClientCount())
	// 开始接收数据
	conn.StartRecv()
	return conn, nil
}
