/**
 * \file ChanConn.go
 * \version
 * \author Jansen
 * \date  2019年01月31日 12:22:43
 * \brief 连接数据管理器
 *
 */

// Package chanconn micserver 中的管道连接，底层默认将其应用于在同一个 App 下的模块间的连接，
// 可以减少消息编解码CPU占用，提高同一APP下的Module交换消息的效率。
package chanconn

import (
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/network/baseio"
	"github.com/liasece/micserver/util/sysutil"
)

// 消息合批时，合并的最大消息数量
const (
	MaxMsgPackSum = 200
)

// ChanConn 的连接状态枚举
const (
	// 未连接
	TCPConnStateNone = 0
	// 已连接
	TCPConnStateLinked = 1
	// 标记不可发送
	TCPConnStateHold = 2
	// 已关闭
	TCPConnStateClosed = 3
)

// chan 连接的错误类型
var (
	ErrSendNilData = errors.New("send nil data")
	ErrCloseed     = errors.New("conn has been closed")
	ErrBufferFull  = errors.New("buffer full")
)

// ChanConn chan 连接
type ChanConn struct {
	*log.Logger

	sendChan chan *msg.MessageBinary
	recvChan chan *msg.MessageBinary

	// 尝试关闭一个连接
	shutdownChan chan struct{}
	state        int32

	// 接收等待通道
	recvmsgchan chan *msg.MessageBinary
	// 消息编解码器
	codec msg.IMsgCodec
}

// Init 初始化一个ChanConn对象
// 	conn: net.Conn对象
// 	sendChanSize: 	发送等待队列中的消息缓冲区大小
// 	sendBufferSize: 发送拼包发送缓冲区大小
// 	recvChanSize: 	接收等待队列中的消息缓冲区大小
// 	recvBufferSize: 接收拼包发送缓冲区大小
// 返回：接收到的 messagebinary 的对象 chan
func (cc *ChanConn) Init(sendChan chan *msg.MessageBinary,
	recvChan chan *msg.MessageBinary) {
	cc.shutdownChan = make(chan struct{})
	cc.sendChan = sendChan
	cc.recvChan = recvChan
	cc.state = TCPConnStateLinked

	// 接收
	cc.recvmsgchan = make(chan *msg.MessageBinary, len(recvChan))
	cc.codec = &msg.DefaultCodec{}
}

// SetMsgCodec chan 中使用的是默认消息编解码器。
// 在 chan 连接中，无法设置消息编解码器，因为其实根本不需要进行消息编解码，
// 为了让 chan 连接实现 IConnect 接口而存在该方法。
func (cc *ChanConn) SetMsgCodec(codec msg.IMsgCodec) {
	cc.Debug("ChanConn can't SetMsgCodec")
}

// GetMsgCodec 该方法将会返回默认消息编解码器。
// 在 chan 连接中，无法设置消息编解码器，因为其实根本不需要进行消息编解码，
// 为了让 chan 连接实现 IConnect 接口而存在该方法。
func (cc *ChanConn) GetMsgCodec() msg.IMsgCodec {
	return cc.codec
}

// SetBanAutoResize 在 chan 连接中，无法设置禁止缓冲区自动扩容。
// 在 chan 连接中，无法设置消息编解码器，因为其实根本不需要进行消息编解码，
// 为了让 chan 连接实现 IConnect 接口而存在该方法。
func (cc *ChanConn) SetBanAutoResize(value bool) {
	cc.Debug("ChanConn can't SetBanAutoResize")
}

// SetLogger 设置连接的 Logger
func (cc *ChanConn) SetLogger(l *log.Logger) {
	cc.Logger = l
}

// StartRecv 开始接收消息
func (cc *ChanConn) StartRecv() {
	go cc.recvThread()
}

// GetRecvMessageChannel 获取消息接收 chan
func (cc *ChanConn) GetRecvMessageChannel() chan *msg.MessageBinary {
	return cc.recvmsgchan
}

// IsAlive 连接是否存活
func (cc *ChanConn) IsAlive() bool {
	if atomic.LoadInt32(&cc.state) == TCPConnStateLinked {
		return true
	}
	return false
}

// RemoteAddr 获取该连接的远程连接地址
func (cc *ChanConn) RemoteAddr() string {
	return "(chan)"
}

// HookProtocal chan 中使用的是默认网络协议（chan通信）。
// 在 chan 连接中，无法设置消息编解码器，因为其实根本不需要进行消息编解码，
// 为了让 chan 连接实现 IConnect 接口而存在该方法。
func (cc *ChanConn) HookProtocal(p baseio.Protocal) {
	cc.Error("ChanConn not support HookProtocal")
}

// Shutdown 尝试关闭此连接
func (cc *ChanConn) Shutdown() error {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if stackInfo, err := sysutil.GetPanicInfo(recover()); err != nil {
			cc.Warn("[ChanConn.shutdownThread] Panic", log.ErrorField(err), log.String("Stack", stackInfo))
		}
	}()
	if cc.state != TCPConnStateLinked {
		return fmt.Errorf("isn't linked")
	}
	sec := atomic.CompareAndSwapInt32(&cc.state,
		TCPConnStateLinked, TCPConnStateHold)
	if sec {
		close(cc.shutdownChan)
	}
	return nil
}

// chan 中使用的是默认网络协议（chan通信），该消息返回空。
// 在 chan 连接中，无法设置消息编解码器，因为其实根本不需要进行消息编解码，
// Read 为了让 chan 连接实现 IConnect 接口而存在该方法。
func (cc *ChanConn) Read(toData []byte) (int, error) {
	return 0, nil
}

// chan 中使用的是默认网络协议（chan通信），该消息返回空。
// 在 chan 连接中，无法设置消息编解码器，因为其实根本不需要进行消息编解码，
// Write 为了让 chan 连接实现 IConnect 接口而存在该方法。
func (cc *ChanConn) Write(data []byte) (int, error) {
	return 0, nil
}

// SendBytes 发送消息 ID 及 Bytes 构成的消息
func (cc *ChanConn) SendBytes(
	cmdid uint16, protodata []byte) error {
	if cc.state >= TCPConnStateHold {
		cc.Warn("[ChanConn.SendBytes] Connection disabled, cancel transmission")
		return ErrCloseed
	}
	msgbinary := msg.DefaultEncodeBytes(cmdid, protodata)

	return cc.SendMessageBinary(msgbinary)
}

// SendMessageBinary 发送 MsgBinary 消息
func (cc *ChanConn) SendMessageBinary(
	msgbinary *msg.MessageBinary) error {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if stackInfo, err := sysutil.GetPanicInfo(recover()); err != nil {
			cc.Warn("[ChanConn.SendMessageBinary] Panic", log.ErrorField(err), log.String("Stack", stackInfo))
		}
	}()
	// 检查连接是否已死亡
	if cc.state >= TCPConnStateHold {
		cc.Warn("[ChanConn.SendMessageBinary] Connection disabled, cancel transmission")
		return ErrCloseed
	}
	// 如果发送数据为空
	if msgbinary == nil {
		cc.Debug("[ChanConn.SendMessageBinary] Send message is empty, cancel sending")
		return ErrSendNilData
	}

	// 检查发送channel是否已经关闭
	select {
	case <-cc.shutdownChan:
		cc.Warn("[ChanConn.SendMessageBinary] Send Channel is off, cancel sending")
		return ErrCloseed
	default:
	}

	// 确认发送channel是否已经关闭
	select {
	case <-cc.shutdownChan:
		cc.Warn("[ChanConn.SendMessageBinary] Send Channel is off, cancel sending")
		return ErrCloseed
	case cc.sendChan <- msgbinary:
		// 遍历已经发送的消息
		// 调用发送回调函数
		msgbinary.OnSendFinish()
		// default:
		// 	cc.Warn("[ChanConn.SendMessageBinary] Send channel buffer full, blocking timeout")
		// 	return ErrBufferFull
	}
	return nil
}

//关闭socket 应该在消息尝试发送完之后执行
func (cc *ChanConn) closeSocket() error {
	cc.state = TCPConnStateClosed
	close(cc.sendChan)
	return nil
}

func (cc *ChanConn) recvThread() {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if stackInfo, err := sysutil.GetPanicInfo(recover()); err != nil {
			cc.Error("[ChanConn.recvThread] Panic", log.ErrorField(err), log.String("Stack", stackInfo))
		}
		close(cc.recvmsgchan)
	}()

	for true {
		if cc.state != TCPConnStateLinked {
			return
		}
		select {
		case msgbinary, ok := <-cc.recvChan:
			if !ok {
				// 连接关闭
				return
			}
			cc.recvmsgchan <- msgbinary
			// cc.Debug("ChanConn Receive messages", log.Uint16("MsgID", msgbinary.GetMsgID()))
		}
	}
}
