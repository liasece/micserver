package connect

import (
	"math/rand"
	"net"

	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/servercomm"
)

// Server 服务器连接，在一个模块的 SubnetManager 中，连接至该模块的任何模块都在该模块中
// 存在一个 Server 连接。
type Server struct {
	BaseConnect

	// 建立连接优先级
	ConnectPriority int64
	// 该连接对方服务器信息
	ModuleInfo *servercomm.ModuleInfo
	// 用于区分该连接是服务器 client task 连接
	serverSCType TServerSCType
}

// InitTCP 初始化一个新的服务器连接
// sctype: 连接的 客户端/服务器 类型
// netconn: 连接的net.Conn对象
func (s *Server) InitTCP(sctype TServerSCType, netconn net.Conn,
	onRecv func(*Server, *msg.MessageBinary),
	onClose func(*Server)) {
	s.BaseConnect.Init()
	s.ModuleInfo = &servercomm.ModuleInfo{}
	s.SetSC(sctype)
	s.ConnectPriority = rand.Int63()
	s.IConnection = NewTCP(netconn, s.Logger,
		ServerSendChanSize, ServerSendBufferSize,
		ServerRecvChanSize, ServerRecvBufferSize)
	// 禁止连接自动扩容缓冲区
	s.IConnection.SetBanAutoResize(true)
	s.IConnection.StartRecv()
	go s.recvMsgThread(s.IConnection.GetRecvMessageChannel(),
		onRecv, onClose)
}

// InitChan 初始化一个新的服务器连接
// sctype: 连接的 客户端/服务器 类型
// sendChan: 发送消息管道
// recvChan: 接收消息管道
func (s *Server) InitChan(sctype TServerSCType,
	sendChan chan *msg.MessageBinary, recvChan chan *msg.MessageBinary,
	onRecv func(*Server, *msg.MessageBinary),
	onClose func(*Server)) {
	s.BaseConnect.Init()
	s.ModuleInfo = &servercomm.ModuleInfo{}
	s.SetSC(sctype)
	s.ConnectPriority = rand.Int63()
	s.IConnection = NewChan(sendChan, recvChan, s.Logger)
	s.IConnection.StartRecv()
	go s.recvMsgThread(s.IConnection.GetRecvMessageChannel(),
		onRecv, onClose)
}

func (s *Server) recvMsgThread(c chan *msg.MessageBinary,
	onRecv func(*Server, *msg.MessageBinary),
	onClose func(*Server)) {
	defer func() {
		if onClose != nil {
			onClose(s)
		}
	}()

	for {
		select {
		case m, ok := <-c:
			if !ok || m == nil {
				return
			}
			if onRecv != nil {
				onRecv(s, m)
			}
		}
	}
}

// SetSC 设置该服务器连接是连接方还是受连接方
func (s *Server) SetSC(sctype TServerSCType) {
	s.serverSCType = sctype
}

// GetSCType 获取该服务器连接是连接方还是受连接方
func (s *Server) GetSCType() TServerSCType {
	return s.serverSCType
}
