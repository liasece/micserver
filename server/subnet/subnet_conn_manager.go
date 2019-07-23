package subnet

import (
	"github.com/liasece/micserver/tcpconn"
	"net"
)

// 添加一个Client
func (this *SubnetManager) AddTCPConn(sctype tcpconn.TServerSCType,
	conn net.Conn, serverid string) *tcpconn.ServerConn {
	tcptask := this.connPool.NewServerConn(sctype, conn, serverid)

	this.Debug("[SubnetManager.AddTCPClient] "+
		"AddTCPClient ServerID[%s] 当前连接数量 LinkSum[%d]",
		serverid, this.connPool.Len())

	return tcptask
}

// 根据tempid 获取一个client
func (this *SubnetManager) GetTCPConn(
	tempid string) *tcpconn.ServerConn {
	return this.connPool.Get(tempid)
}

func (this *SubnetManager) RangeConn(
	callback func(*tcpconn.ServerConn) bool) {
	this.connPool.Range(callback)
}

// 根据tempid移除一个client
func (this *SubnetManager) RemoveTCPConn(tempid string) {
	this.connPool.Remove(tempid)
}

// 修改链接的 tempip
func (this *SubnetManager) ChangeTCPConnTempid(
	tcptask *tcpconn.ServerConn, newtempid string) error {
	return this.connPool.ChangeTempid(tcptask, newtempid)
}

func (this *SubnetManager) RandomGetTCPConn(
	servertype string) *tcpconn.ServerConn {
	return this.connPool.GetRandom(servertype)
}
