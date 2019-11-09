package connect

import (
	"math/rand"
	"net"

	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/servercomm"
)

type TServerSCType uint32

const (
	ServerSCTypeNone   TServerSCType = 1
	ServerSCTypeTask   TServerSCType = 2
	ServerSCTypeClient TServerSCType = 3
)

// 服务器连接发送消息缓冲要考虑到服务器处理消息的能力
const ServerSendChanSize = 100000

// 发送缓冲大小，用于将多个小消息拼接发送的缓冲大小
const ServerSendBufferSize = msg.MessageMaxSize * 10

// 服务器连接发送消息缓冲要考虑到服务器处理消息的能力
const ServerRecvChanSize = 100000

// 发送缓冲大小，用于将多个小消息拼接发送的缓冲大小
const ServerRecvBufferSize = msg.MessageMaxSize * 10

type Server struct {
	BaseConnect

	// 建立连接优先级
	ConnectPriority int64
	// 该连接对方服务器信息
	ServerInfo *servercomm.ServerInfo
	// 用于区分该连接是服务器 client task 连接
	serverSCType TServerSCType
}

// 初始化一个新的服务器连接
// sctype: 连接的 客户端/服务器 类型
// netconn: 连接的net.Conn对象
func (this *Server) InitTCP(sctype TServerSCType, netconn net.Conn,
	onRecv func(*Server, *msg.MessageBinary),
	onClose func(*Server)) {
	this.BaseConnect.Init()
	this.ServerInfo = &servercomm.ServerInfo{}
	this.SetSC(sctype)
	this.ConnectPriority = rand.Int63()
	this.IConnection = NewTCP(netconn, this.Logger,
		ServerSendChanSize, ServerSendBufferSize,
		ServerRecvChanSize, ServerRecvBufferSize)
	// 禁止连接自动扩容缓冲区
	this.IConnection.SetBanAutoResize(true)
	this.IConnection.StartRecv()
	go this.recvMsgThread(this.IConnection.GetRecvMessageChannel(),
		onRecv, onClose)
}

// 初始化一个新的服务器连接
// sctype: 连接的 客户端/服务器 类型
// sendChan: 发送消息管道
// recvChan: 接收消息管道
func (this *Server) InitChan(sctype TServerSCType,
	sendChan chan *msg.MessageBinary, recvChan chan *msg.MessageBinary,
	onRecv func(*Server, *msg.MessageBinary),
	onClose func(*Server)) {
	this.BaseConnect.Init()
	this.ServerInfo = &servercomm.ServerInfo{}
	this.SetSC(sctype)
	this.ConnectPriority = rand.Int63()
	this.IConnection = NewChan(sendChan, recvChan, this.Logger)
	this.IConnection.StartRecv()
	go this.recvMsgThread(this.IConnection.GetRecvMessageChannel(),
		onRecv, onClose)
}

func (this *Server) recvMsgThread(c chan *msg.MessageBinary,
	onRecv func(*Server, *msg.MessageBinary),
	onClose func(*Server)) {
	defer func() {
		if onClose != nil {
			onClose(this)
		}
	}()

	for {
		select {
		case m, ok := <-c:
			if !ok || m == nil {
				return
			}
			if onRecv != nil {
				onRecv(this, m)
			}
		}
	}
}

func (this *Server) SetSC(sctype TServerSCType) {
	this.serverSCType = sctype
}

func (this *Server) GetSCType() TServerSCType {
	return this.serverSCType
}
