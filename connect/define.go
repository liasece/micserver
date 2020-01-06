package connect

import (
	"github.com/liasece/micserver/msg"
)

/*
connect 中使用的一些常量
*/

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

// 客户端连接发送消息缓冲不宜过大， 10*64KiB*100000连接=64GiB
const ClientConnSendChanSize = 256

// 发送缓冲大小，用于将多个小消息拼接发送的缓冲大小
const ClientConnSendBufferSize = 16 * 1024

// 客户端连接发送消息缓冲不宜过大， 10*64KiB*100000连接=64GiB
const ClientConnRecvChanSize = 256

// 发送缓冲大小，用于将多个小消息拼接发送的缓冲大小
const ClientConnRecvBufferSize = 256 * 1024

// 客户端连接池分组数量
const mClientPoolGroupSum = 10
