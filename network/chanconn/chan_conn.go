/**
 * \file ChanConn.go
 * \version
 * \author Jansen
 * \date  2019年01月31日 12:22:43
 * \brief 连接数据管理器
 *
 */

/*
micserver 中的管道连接，底层默认将其应用于在同一个 App 下的模块间的连接，
可以减少消息编解码CPU占用，提高同一APP下的Module交换消息的效率。
*/
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
	TCPCONNSTATE_NONE = 0
	// 已连接
	TCPCONNSTATE_LINKED = 1
	// 标记不可发送
	TCPCONNSTATE_HOLD = 2
	// 已关闭
	TCPCONNSTATE_CLOSED = 3
)

// chan 连接的错误类型
var (
	ErrSendNilData = errors.New("send nil data")
	ErrCloseed     = errors.New("conn has been closed")
	ErrBufferFull  = errors.New("buffer full")
)

// chan 连接
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

// 初始化一个ChanConn对象
// 	conn: net.Conn对象
// 	sendChanSize: 	发送等待队列中的消息缓冲区大小
// 	sendBufferSize: 发送拼包发送缓冲区大小
// 	recvChanSize: 	接收等待队列中的消息缓冲区大小
// 	recvBufferSize: 接收拼包发送缓冲区大小
// 返回：接收到的 messagebinary 的对象 chan
func (this *ChanConn) Init(sendChan chan *msg.MessageBinary,
	recvChan chan *msg.MessageBinary) {
	this.shutdownChan = make(chan struct{})
	this.sendChan = sendChan
	this.recvChan = recvChan
	this.state = TCPCONNSTATE_LINKED

	// 接收
	this.recvmsgchan = make(chan *msg.MessageBinary, len(recvChan))
	this.codec = &msg.DefaultCodec{}
}

// chan 中使用的是默认消息编解码器。
// 在 chan 连接中，无法设置消息编解码器，因为其实根本不需要进行消息编解码，
// 为了让 chan 连接实现 IConnect 接口而存在该方法。
func (this *ChanConn) SetMsgCodec(codec msg.IMsgCodec) {
	this.Debug("ChanConn can't SetMsgCodec")
}

// 该方法将会返回默认消息编解码器。
// 在 chan 连接中，无法设置消息编解码器，因为其实根本不需要进行消息编解码，
// 为了让 chan 连接实现 IConnect 接口而存在该方法。
func (this *ChanConn) GetMsgCodec() msg.IMsgCodec {
	return this.codec
}

// 在 chan 连接中，无法设置禁止缓冲区自动扩容。
// 在 chan 连接中，无法设置消息编解码器，因为其实根本不需要进行消息编解码，
// 为了让 chan 连接实现 IConnect 接口而存在该方法。
func (this *ChanConn) SetBanAutoResize(value bool) {
	this.Debug("ChanConn can't SetBanAutoResize")
}

// 设置连接的 Logger
func (this *ChanConn) SetLogger(l *log.Logger) {
	this.Logger = l
}

// 开始接收消息
func (this *ChanConn) StartRecv() {
	go this.recvThread()
}

// 获取消息接收 chan
func (this *ChanConn) GetRecvMessageChannel() chan *msg.MessageBinary {
	return this.recvmsgchan
}

// 连接是否存活
func (this *ChanConn) IsAlive() bool {
	if atomic.LoadInt32(&this.state) == TCPCONNSTATE_LINKED {
		return true
	}
	return false
}

// 获取该连接的远程连接地址
func (this *ChanConn) RemoteAddr() string {
	return "(chan)"
}

// chan 中使用的是默认网络协议（chan通信）。
// 在 chan 连接中，无法设置消息编解码器，因为其实根本不需要进行消息编解码，
// 为了让 chan 连接实现 IConnect 接口而存在该方法。
func (this *ChanConn) HookProtocal(p baseio.Protocal) {
	this.Error("ChanConn not support HookProtocal")
}

// 尝试关闭此连接
func (this *ChanConn) Shutdown() error {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := sysutil.GetPanicInfo(recover()); err != nil {
			this.Warn("[ChanConn.shutdownThread] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
		}
	}()
	if this.state != TCPCONNSTATE_LINKED {
		return fmt.Errorf("isn't linked")
	}
	sec := atomic.CompareAndSwapInt32(&this.state,
		TCPCONNSTATE_LINKED, TCPCONNSTATE_HOLD)
	if sec {
		close(this.shutdownChan)
	}
	return nil
}

// chan 中使用的是默认网络协议（chan通信），该消息返回空。
// 在 chan 连接中，无法设置消息编解码器，因为其实根本不需要进行消息编解码，
// 为了让 chan 连接实现 IConnect 接口而存在该方法。
func (this *ChanConn) Read(toData []byte) (int, error) {
	return 0, nil
}

// chan 中使用的是默认网络协议（chan通信），该消息返回空。
// 在 chan 连接中，无法设置消息编解码器，因为其实根本不需要进行消息编解码，
// 为了让 chan 连接实现 IConnect 接口而存在该方法。
func (this *ChanConn) Write(data []byte) (int, error) {
	return 0, nil
}

// 发送消息 ID 及 Bytes 构成的消息
func (this *ChanConn) SendBytes(
	cmdid uint16, protodata []byte) error {
	if this.state >= TCPCONNSTATE_HOLD {
		this.Warn("[ChanConn.SendBytes] 连接已失效，取消发送")
		return ErrCloseed
	}
	msgbinary := msg.DefaultEncodeBytes(cmdid, protodata)

	return this.SendMessageBinary(msgbinary)
}

// 发送 MsgBinary 消息
func (this *ChanConn) SendMessageBinary(
	msgbinary *msg.MessageBinary) error {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := sysutil.GetPanicInfo(recover()); err != nil {
			this.Warn("[ChanConn.SendMessageBinary] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
		}
	}()
	// 检查连接是否已死亡
	if this.state >= TCPCONNSTATE_HOLD {
		this.Warn("[ChanConn.SendMessageBinary] 连接已失效，取消发送")
		return ErrCloseed
	}
	// 如果发送数据为空
	if msgbinary == nil {
		this.Debug("[ChanConn.SendMessageBinary] 发送消息为空，取消发送")
		return ErrSendNilData
	}

	// 检查发送channel是否已经关闭
	select {
	case <-this.shutdownChan:
		this.Warn("[ChanConn.SendMessageBinary] 发送Channel已关闭，取消发送")
		return ErrCloseed
	default:
	}

	// 确认发送channel是否已经关闭
	select {
	case <-this.shutdownChan:
		this.Warn("[ChanConn.SendMessageBinary] 发送Channel已关闭，取消发送")
		return ErrCloseed
	case this.sendChan <- msgbinary:
		// 遍历已经发送的消息
		// 调用发送回调函数
		msgbinary.OnSendFinish()
		// default:
		// 	this.Warn("[ChanConn.SendMessageBinary] 发送Channel缓冲区满，阻塞超时")
		// 	return ErrBufferFull
	}
	return nil
}

//关闭socket 应该在消息尝试发送完之后执行
func (this *ChanConn) closeSocket() error {
	this.state = TCPCONNSTATE_CLOSED
	close(this.sendChan)
	return nil
}

func (this *ChanConn) recvThread() {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := sysutil.GetPanicInfo(recover()); err != nil {
			this.Error("[ChanConn.recvThread] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
		}
		close(this.recvmsgchan)
	}()

	for true {
		if this.state != TCPCONNSTATE_LINKED {
			return
		}
		select {
		case msgbinary, ok := <-this.recvChan:
			if !ok {
				// 连接关闭
				return
			}
			this.recvmsgchan <- msgbinary
			// this.Debug("ChanConn 接收消息 %d", msgbinary.CmdID)
		}
	}
}
