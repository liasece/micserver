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
type ClientConnManager struct {
	*log.Logger
	handle.ClientTcpHandler

	connPool connect.ClientConnPool
}

func (this *ClientConnManager) Init(moduleID string) {
	this.connPool.Init(int32(util.GetStringHash(moduleID)))
}

func (this *ClientConnManager) addTCPClient(
	netConn net.Conn) (*connect.ClientConn, error) {
	conn, err := this.connPool.NewClientConn(netConn, this.OnConnectRecv,
		this.onConnectClose)
	if err != nil {
		return nil, err
	}
	conn.SetLogger(this.Logger)

	this.OnNewConn(conn)
	curtime := time.Now().Unix()
	conn.SetTerminateTime(curtime + 20) // 20秒以后还没有验证通过就断开连接

	conn.Debug("[ClientConnManager.addTCPClient] "+
		"新增连接数 当前连接数量 NowSum[%d]",
		this.GetClientTcpSocketCount())
	return conn, nil
}

func (this *ClientConnManager) onConnectClose(conn *connect.ClientConn) {
	this.RemoveTaskByTmpID(conn.GetConnectID())
}

func (this *ClientConnManager) GetClientConn(
	webtaskid string) *connect.ClientConn {
	return this.connPool.Get(webtaskid)
}

func (this *ClientConnManager) GetClientTcpSocketCount() uint32 {
	return this.connPool.Len()
}

func (this *ClientConnManager) remove(connectid string) {
	value := this.GetClientConn(connectid)
	if value == nil {
		return
	}
	this.connPool.Remove(connectid)
}

func (this *ClientConnManager) RemoveTaskByTmpID(
	connectid string) {
	this.remove(connectid)
}

// 遍历所有的连接
func (this *ClientConnManager) Range(
	callback func(string, *connect.ClientConn) bool) {
	this.connPool.Range(func(value *connect.ClientConn) bool {
		return callback(value.GetConnectID(), value)
	})
}

// 遍历所有的连接，检查需要移除的连接
func (this *ClientConnManager) RangeRemove(
	callback func(*connect.ClientConn) bool) {
	removelist := make([]string, 0)
	this.connPool.Range(func(value *connect.ClientConn) bool {
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

	this.Debug("[ClientConnManager.ExecRemove] "+
		"条件删除连接数 RemoveSum[%d] 当前连接数量 LinkSum[%d]",
		len(removelist), this.GetClientTcpSocketCount())
}

func (this *ClientConnManager) StartAddClientTcpSocketHandle(addr string) {
	// 由于部分 NAT 主机没有网卡概念，需要自己配置IP
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		this.Error("[ClientConnManager.StartAddClientTcpSocketHandle] %s",
			err.Error())
		return
	}
	this.Debug("[ClientConnManager.StartAddClientTcpSocketHandle] "+
		"Gateway Client TCP服务启动成功 IPPort[%s]", addr)
	go func() {
		for {
			// 接受连接
			netConn, err := ln.Accept()
			if err != nil {
				// handle error
				this.Error("[ClientConnManager.StartAddClientTcpSocketHandle] "+
					"Accept() ERR:%q",
					err.Error())
				continue
			}
			this.addTCPClient(netConn)
		}
	}()
}
