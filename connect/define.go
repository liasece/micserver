package connect

import (
	"github.com/liasece/micserver/msg"
)

/*
connect 中使用的一些常量
*/

// TServerSCType server server-client type
type TServerSCType uint32

// TServerSCType
const (
	ServerSCTypeNone   TServerSCType = 1
	ServerSCTypeTask   TServerSCType = 2
	ServerSCTypeClient TServerSCType = 3
)

// ServerSendChanSize 服务器连接发送消息缓冲要考虑到服务器处理消息的能力
const ServerSendChanSize = 100000

// ServerSendBufferSize 发送缓冲大小，用于将多个小消息拼接发送的缓冲大小
const ServerSendBufferSize = msg.MessageMaxSize * 10

// ServerRecvChanSize 服务器连接发送消息缓冲要考虑到服务器处理消息的能力
const ServerRecvChanSize = 100000

// ServerRecvBufferSize 发送缓冲大小，用于将多个小消息拼接发送的缓冲大小
const ServerRecvBufferSize = msg.MessageMaxSize * 10

// ClientConnSendChanSize 客户端连接发送消息缓冲不宜过大， 10*64KiB*100000连接=64GiB
const ClientConnSendChanSize = 256

// ClientConnSendBufferSize 发送缓冲大小，用于将多个小消息拼接发送的缓冲大小
const ClientConnSendBufferSize = 16 * 1024

// ClientConnRecvChanSize 客户端连接发送消息缓冲不宜过大， 10*64KiB*100000连接=64GiB
const ClientConnRecvChanSize = 256

// ClientConnRecvBufferSize 发送缓冲大小，用于将多个小消息拼接发送的缓冲大小
const ClientConnRecvBufferSize = 256 * 1024

// mClientPoolGroupSum 客户端连接池分组数量
const mClientPoolGroupSum = 10
