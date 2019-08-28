package manager

import (
	"github.com/liasece/micserver/connect"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/server/gate/handle"
	"github.com/liasece/micserver/util"
	"net"
	"time"
)

// websocket连接管理器
type ClientManager struct {
	*log.Logger
	handle.ClientTcpHandler

	connPool connect.ClientPool
}

func (this *ClientManager) Init(moduleID string) {
	this.connPool.Init(int32(util.GetStringHash(moduleID)))
}

func (this *ClientManager) addTCPClient(
	netConn net.Conn) (*connect.Client, error) {
	conn, err := this.connPool.NewClient(netConn, this.OnConnectRecv,
		this.onConnectClose)
	if err != nil {
		return nil, err
	}
	conn.SetLogger(this.Logger)

	this.OnNewClient(conn)
	curtime := time.Now().Unix()
	conn.SetTerminateTime(curtime + 20) // 20秒以后还没有验证通过就断开连接

	conn.Debug("[ClientManager.addTCPClient] "+
		"新增连接数 当前连接数量 NowSum[%d]",
		this.GetClientTcpSocketCount())
	return conn, nil
}

func (this *ClientManager) onConnectClose(conn *connect.Client) {
	this.RemoveTaskByTmpID(conn.GetConnectID())
}

func (this *ClientManager) GetClient(
	webtaskid string) *connect.Client {
	return this.connPool.Get(webtaskid)
}

func (this *ClientManager) GetClientTcpSocketCount() uint32 {
	return this.connPool.Len()
}

func (this *ClientManager) remove(connectid string) {
	value := this.GetClient(connectid)
	if value == nil {
		return
	}
	this.connPool.Remove(connectid)
}

func (this *ClientManager) RemoveTaskByTmpID(
	connectid string) {
	this.remove(connectid)
}

// 遍历所有的连接
func (this *ClientManager) Range(
	callback func(string, *connect.Client) bool) {
	this.connPool.Range(func(value *connect.Client) bool {
		return callback(value.GetConnectID(), value)
	})
}

// 遍历所有的连接，检查需要移除的连接
func (this *ClientManager) RangeRemove(
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

	this.Debug("[ClientManager.ExecRemove] "+
		"条件删除连接数 RemoveSum[%d] 当前连接数量 LinkSum[%d]",
		len(removelist), this.GetClientTcpSocketCount())
}

func (this *ClientManager) StartAddClientTcpSocketHandle(addr string) {
	// 由于部分 NAT 主机没有网卡概念，需要自己配置IP
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		this.Error("[ClientManager.StartAddClientTcpSocketHandle] %s",
			err.Error())
		return
	}
	this.Debug("[ClientManager.StartAddClientTcpSocketHandle] "+
		"Gateway Client TCP服务启动成功 IPPort[%s]", addr)
	go func() {
		for {
			// 接受连接
			netConn, err := ln.Accept()
			if err != nil {
				// handle error
				this.Error("[ClientManager.StartAddClientTcpSocketHandle] "+
					"Accept() ERR:%q",
					err.Error())
				continue
			}
			this.OnAcceptConnect(netConn)
			this.addTCPClient(netConn)
		}
	}()
}
