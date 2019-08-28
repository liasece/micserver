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
type ClientSocketManager struct {
	*log.Logger
	handle.ClientTcpHandler

	connPool connect.ClientConnPool
}

func (this *ClientSocketManager) Init(moduleID string) {
	this.connPool.Init(int32(util.GetStringHash(moduleID)))
}

func (this *ClientSocketManager) AddClientTcpSocket(
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

	conn.Debug("[ClientSocketManager.AddClientTcpSocket] "+
		"新增连接数 当前连接数量 NowSum[%d]",
		this.GetClientTcpSocketCount())
	return conn, nil
}

func (this *ClientSocketManager) onConnectClose(conn *connect.ClientConn) {
	this.RemoveTaskByTmpID(conn.GetConnectID())
}

func (this *ClientSocketManager) GetTaskByTmpID(
	webtaskid string) *connect.ClientConn {
	return this.connPool.Get(webtaskid)
}

func (this *ClientSocketManager) GetClientTcpSocketCount() uint32 {
	return this.connPool.Len()
}

func (this *ClientSocketManager) remove(tempid string) {
	value := this.GetTaskByTmpID(tempid)
	if value == nil {
		return
	}
	this.connPool.Remove(tempid)
}

func (this *ClientSocketManager) RemoveTaskByTmpID(
	tempid string) {
	this.remove(tempid)
}

// 遍历所有的连接
func (this *ClientSocketManager) ExecAllUsers(
	callback func(string, *connect.ClientConn)) {
	this.connPool.Range(func(value *connect.ClientConn) {
		callback(value.GetConnectID(), value)
	})
}

// 遍历所有的连接，检查需要移除的连接
func (this *ClientSocketManager) ExecRemove(
	callback func(*connect.ClientConn) bool) {
	removelist := make([]string, 0)
	this.connPool.Range(func(value *connect.ClientConn) {
		// 遍历所有的连接
		if callback(value) {
			// 该连接需要被移除
			removelist = append(removelist, value.GetConnectID())
			value.Terminate()
		}
	})
	for _, v := range removelist {
		this.remove(v)
	}

	this.Debug("[ClientSocketManager.ExecRemove] "+
		"条件删除连接数 RemoveSum[%d] 当前连接数量 LinkSum[%d]",
		len(removelist), this.GetClientTcpSocketCount())
}

func (this *ClientSocketManager) StartAddClientTcpSocketHandle(addr string) {
	// 由于部分 NAT 主机没有网卡概念，需要自己配置IP
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		this.Error("[ClientSocketManager.StartAddClientTcpSocketHandle] %s",
			err.Error())
		return
	}
	this.Debug("[ClientSocketManager.StartAddClientTcpSocketHandle] "+
		"Gateway Client TCP服务启动成功 IPPort[%s]", addr)
	go func() {
		for {
			// 接受连接
			netConn, err := ln.Accept()
			if err != nil {
				// handle error
				this.Error("[ClientSocketManager.StartAddClientTcpSocketHandle] "+
					"Accept() ERR:%q",
					err.Error())
				continue
			}
			this.AddClientTcpSocket(netConn)
		}
	}()
}
