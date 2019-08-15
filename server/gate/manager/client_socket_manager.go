package manager

import (
	"github.com/liasece/micserver/connect"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/server/gate/handle"
	"github.com/liasece/micserver/util"
	"io"
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
	conn, err := this.connPool.NewClientConn(netConn)
	conn.Logger = this.Logger
	if err != nil {
		return nil, err
	}
	curtime := time.Now().Unix()
	conn.SetTerminateTime(curtime + 20) // 20秒以后还没有验证通过就断开连接

	conn.Debug("[ClientSocketManager.AddClientTcpSocket] "+
		"新增连接数 当前连接数量 NowSum[%d]",
		this.GetClientTcpSocketCount())
	return conn, nil
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
		callback(value.Tempid, value)
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
			removelist = append(removelist, value.Tempid)
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

func (this *ClientSocketManager) onNewConnect(netConn net.Conn) {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := util.GetPanicInfo(recover()); err != nil {
			this.Error("[onNewConnect] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
		}
	}()
	this.Debug("[onNewConnect] Receive one netConn connect json")
	conn, err := this.AddClientTcpSocket(netConn)
	if err != nil || conn == nil {
		this.Error("[onNewConnect] "+
			"创建 ClientTcpSocket 对象失败，断开连接 Err[%s]", err.Error())
		return
	}
	netbuffer := util.NewIOBuffer(netConn, 64*1024)
	msgReader := msg.NewMessageBinaryReader(netbuffer)

	// 所有连接都需要经过加密
	// conn.Encryption = base.EncryptionTypeXORSimple

	for {
		if !conn.Check() {
			// 强制移除客户端连接
			this.RemoveTaskByTmpID(conn.Tempid)
			return
		}
		// 设置阻塞读取过期时间
		err := netConn.SetReadDeadline(
			time.Now().Add(time.Duration(time.Millisecond * 250)))
		if err != nil {
			conn.Error("[onNewConnect] SetReadDeadline Err[%s]",
				err.Error())
		}
		// buffer从连接中读取socket数据
		_, err = netbuffer.ReadFromReader()

		// 异常
		if err != nil {
			if err == io.EOF {
				conn.Debug("[onNewConnect] "+
					"Scoket数据读写异常,断开连接了,"+
					"scoket返回 Err[%s]", err.Error())
				this.RemoveTaskByTmpID(conn.Tempid)
				return
			} else {
				continue
			}
		}

		err = msgReader.RangeMsgBinary(func(msgbinary *msg.MessageBinary) {
			if conn.Encryption != msg.EncryptionTypeNone &&
				msgbinary.CmdMask != conn.Encryption {
				conn.Error("加密方式错误，加密方式应为 %d 此消息为 %d "+
					"MsgID[%d]", conn.Encryption,
					msgbinary.CmdMask, msgbinary.CmdID)
			} else {
				// 解析消息
				this.OnRecvSocketPackage(conn, msgbinary)
			}
		})
		if err != nil {
			conn.Error("[onNewConnect] 解析消息错误，断开连接 "+
				"Err[%s]", err.Error())
			// 强制移除客户端连接
			this.RemoveTaskByTmpID(conn.Tempid)
			return
		}
	}
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
			go this.onNewConnect(netConn)
		}
	}()
}
