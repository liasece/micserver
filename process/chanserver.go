// Package process 处于同一个进程中的 module ，可以使用 chan 进行通信，不必利用 TCP 中转，
// 通过 process ， Module 可以知道自己的连接目标是否是本进程的，如果是本进程的即可
// 使用 process 包获取到对方的消息通信 chan ，通过 chan 连接对端。
package process

import (
	"sync"

	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/servercomm"
)

// ChanServerHandshake chan 连接握手协议， Module 之间如果想要使用 chan 通信，则先需要通过一个握手 chan
// 进行握手，交换双方的消息接收发送 chan 。
type ChanServerHandshake struct {
	ModuleInfo    *servercomm.ModuleInfo
	ClientMsgChan chan *msg.MessageBinary
	ServerMsgChan chan *msg.MessageBinary
	Seq           int
}

var (
	_gServerChan sync.Map
)

// AddServerChan 增加一个模块的握手 chan
func AddServerChan(id string, serverChan chan *ChanServerHandshake) {
	_gServerChan.Store(id, serverChan)
}

// DeleteServerChan 删除一个模块的握手 chan
func DeleteServerChan(id string) {
	_gServerChan.Delete(id)
}

// GetServerChan 获取一个模块的握手 chan
func GetServerChan(id string) chan *ChanServerHandshake {
	if vi, ok := _gServerChan.Load(id); ok {
		if v, ok := vi.(chan *ChanServerHandshake); ok {
			return v
		}
	}
	return nil
}
