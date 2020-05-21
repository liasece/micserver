/**
 * \file TCPConn.go
 * \version
 * \author Jansen
 * \date  2019年01月31日 12:22:43
 * \brief 连接数据管理器
 *
 */

// Package tcpconn micserver 中的 TCP 连接管理，默认情况下，跨 App 的 Module 使用的便是 TCP 进行消息
// 通信，常用的客户端连接的协议也是 TCP 协议。
package tcpconn

import (
	"fmt"
	"io"
	"net"
	"sync/atomic"
	"time"

	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/network/baseio"
	"github.com/liasece/micserver/util/buffer"
	"github.com/liasece/micserver/util/sysutil"
)

// 消息合批时，合并的最大消息数量
const (
	MaxMsgPackSum = 200
)

// TCPConn 的连接状态枚举
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

// TCPConn TCP 连接
type TCPConn struct {
	*log.Logger
	conn net.Conn
	work baseio.Worker

	// 尝试关闭一个连接
	shutdownChan chan struct{}
	state        int32

	// 发送等待通道
	sendmsgchan chan *msg.MessageBinary
	// 发送缓冲区
	sendBuffer *buffer.IOBuffer
	// 当前等待发送出去的数据总大小
	waitingSendBufferLength int64
	// 等待发送数据的总大小
	maxWaitingSendBufferLength int

	// 已合批发送消息缓冲数组
	sendJoinedMessageBinaryBuffer []*msg.MessageBinary

	// 接收等待通道
	recvmsgchan chan *msg.MessageBinary
	// 接收缓冲区
	recvBuffer *buffer.IOBuffer
	// 消息编解码器
	codec msg.IMsgCodec
}

// Init 初始化一个TCPConn对象
// 	conn: net.Conn对象
// 	sendChanSize: 	发送等待队列中的消息缓冲区大小
// 	sendBufferSize: 发送拼包发送缓冲区大小
// 	recvChanSize: 	接收等待队列中的消息缓冲区大小
// 	recvBufferSize: 接收拼包发送缓冲区大小
// 返回：接收到的 messagebinary 的对象 chan
func (tcpConn *TCPConn) Init(conn net.Conn,
	sendChanSize int, sendBufferSize int,
	recvChanSize int, recvBufferSize int) {
	tcpConn.shutdownChan = make(chan struct{})
	tcpConn.conn = conn
	tcpConn.work.Init(conn)
	tcpConn.state = TCPConnStateLinked

	// 发送
	tcpConn.sendmsgchan = make(chan *msg.MessageBinary, sendChanSize)
	tcpConn.maxWaitingSendBufferLength = msg.MessageMaxSize * sendChanSize
	tcpConn.sendBuffer = buffer.NewIOBuffer(nil, sendBufferSize)
	tcpConn.sendBuffer.Logger = tcpConn.Logger
	tcpConn.sendJoinedMessageBinaryBuffer = make([]*msg.MessageBinary,
		MaxMsgPackSum)
	go tcpConn.sendThread()

	// 接收
	tcpConn.recvmsgchan = make(chan *msg.MessageBinary, recvChanSize)
	tcpConn.recvBuffer = buffer.NewIOBuffer(tcpConn, recvBufferSize)
	tcpConn.recvBuffer.Logger = tcpConn.Logger
	tcpConn.codec = &msg.DefaultCodec{}
}

// SetBanAutoResize 设置禁止缓冲区自动扩容
func (tcpConn *TCPConn) SetBanAutoResize(value bool) {
	tcpConn.sendBuffer.SetBanAutoResize(value)
	tcpConn.recvBuffer.SetBanAutoResize(value)
}

// SetMsgCodec 设置消息编解码器
func (tcpConn *TCPConn) SetMsgCodec(codec msg.IMsgCodec) {
	tcpConn.codec = codec
}

// GetMsgCodec 获取消息编解码器
func (tcpConn *TCPConn) GetMsgCodec() msg.IMsgCodec {
	return tcpConn.codec
}

// SetLogger 设置 Logger
func (tcpConn *TCPConn) SetLogger(l *log.Logger) {
	tcpConn.Logger = l
}

// StartRecv 开始接收消息
func (tcpConn *TCPConn) StartRecv() {
	go tcpConn.recvThread()
}

// GetRecvMessageChannel 获取接收消息 chan
func (tcpConn *TCPConn) GetRecvMessageChannel() chan *msg.MessageBinary {
	return tcpConn.recvmsgchan
}

// IsAlive 判断连接是否存活
func (tcpConn *TCPConn) IsAlive() bool {
	if atomic.LoadInt32(&tcpConn.state) == TCPConnStateLinked {
		return true
	}
	return false
}

// RemoteAddr 获取连接的远程地址
func (tcpConn *TCPConn) RemoteAddr() string {
	return tcpConn.conn.RemoteAddr().String()
}

// HookProtocal 设置网络层协议
func (tcpConn *TCPConn) HookProtocal(p baseio.Protocal) {
	tcpConn.work.HookProtocal(p)
}

// Shutdown 尝试关闭此连接
func (tcpConn *TCPConn) Shutdown() error {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if stackInfo, err := sysutil.GetPanicInfo(recover()); err != nil {
			tcpConn.Warn("[TCPConn.shutdownThread] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
		}
	}()
	if tcpConn.state != TCPConnStateLinked {
		return fmt.Errorf("isn't linked")
	}
	sec := atomic.CompareAndSwapInt32(&tcpConn.state,
		TCPConnStateLinked, TCPConnStateHold)
	if sec {
		close(tcpConn.shutdownChan)
	}
	return nil
}

// Read 读数据
func (tcpConn *TCPConn) Read(toData []byte) (int, error) {
	return tcpConn.work.Read(toData)
}

// Write 写数据
func (tcpConn *TCPConn) Write(data []byte) (int, error) {
	return tcpConn.work.Write(data)
}

// SendBytes 发送由消息 ID 及 Bytes 构成的消息
func (tcpConn *TCPConn) SendBytes(
	cmdid uint16, protodata []byte) error {
	if tcpConn.state >= TCPConnStateHold {
		tcpConn.Warn("[TCPConn.SendBytes] 连接已失效，取消发送")
		return ErrCloseed
	}
	msgbinary := tcpConn.codec.EncodeBytes(cmdid, protodata)

	return tcpConn.SendMessageBinary(msgbinary)
}

// SendMessageBinary 发送 MsgBinary 消息
func (tcpConn *TCPConn) SendMessageBinary(
	msgbinary *msg.MessageBinary) error {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if stackInfo, err := sysutil.GetPanicInfo(recover()); err != nil {
			tcpConn.Warn("[TCPConn.SendMessageBinary] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
		}
	}()
	// 检查连接是否已死亡
	if tcpConn.state >= TCPConnStateHold {
		tcpConn.Warn("[TCPConn.SendMessageBinary] 连接已失效，取消发送")
		return ErrCloseed
	}
	// 如果发送数据为空
	if msgbinary == nil {
		tcpConn.Debug("[TCPConn.SendMessageBinary] 发送消息为空，取消发送")
		return ErrSendNilData
	}

	// 检查发送channel是否已经关闭
	select {
	case <-tcpConn.shutdownChan:
		tcpConn.Warn("[TCPConn.SendMessageBinary] 发送Channel已关闭，取消发送")
		return ErrCloseed
	default:
	}

	// 检查等待缓冲区数据是否已满
	// if tcpConn.waitingSendBufferLength > int64(tcpConn.maxWaitingSendBufferLength) {
	// 	tcpConn.Error("[TCPConn.SendMessageBinary] 等待发送缓冲区满")
	// 	return ErrBufferFull
	// }

	// 确认发送channel是否已经关闭
	select {
	case <-tcpConn.shutdownChan:
		tcpConn.Warn("[TCPConn.SendMessageBinary] 发送Channel已关闭，取消发送")
		return ErrCloseed
	case tcpConn.sendmsgchan <- msgbinary:
		atomic.AddInt64(&tcpConn.waitingSendBufferLength,
			int64(msgbinary.GetTotalLength()))
		// default:
		// 	tcpConn.Warn("[TCPConn.SendMessageBinary] 发送Channel缓冲区满，阻塞超时")
		// 	return ErrBufferFull
	}
	return nil
}

// sendThread 消息发送线程 必须单线程执行
func (tcpConn *TCPConn) sendThread() {
	for {
		if tcpConn.asyncSendCmd() {
			// 正常退出
			break
		}
	}
	// 用于通知发送线程，发送channel已关闭
	tcpConn.Debug("[TCPConn.sendThread] 断开连接")
	// close(tcpConn.stopChan)
	err := tcpConn.closeSocket()
	if err != nil {
		tcpConn.Error("[TCPConn.sendThread] closeSocket Err[%s]",
			err.Error())
	}
}

//关闭socket 应该在消息尝试发送完之后执行
func (tcpConn *TCPConn) closeSocket() error {
	tcpConn.state = TCPConnStateClosed
	return tcpConn.conn.Close()
}

// asyncSendCmd 异步方式发送消息 必须单线程执行
func (tcpConn *TCPConn) asyncSendCmd() (normalreturn bool) {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if stackInfo, err := sysutil.GetPanicInfo(recover()); err != nil {
			tcpConn.Error("[TCPConn.asyncSendCmd] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
			normalreturn = false
		}
	}()

	isrunning := true
	for isrunning {
		select {
		case msg, ok := <-tcpConn.sendmsgchan:
			if msg == nil || !ok {
				tcpConn.Warn("[TCPConn.asyncSendCmd] " +
					"Channle已关闭，发送行为终止")
				break
			}
			tcpConn.sendMsgList(msg)
		case <-tcpConn.shutdownChan:
			// 线程被主动关闭
			isrunning = false
			break
		}
	}

	// 线程准备退出，执行收尾工作，尝试将未发送的消息发送出去
	waitting := true
	for waitting {
		select {
		case msg, ok := <-tcpConn.sendmsgchan:
			// 从发送chan中获取一条消息
			if msg == nil || !ok {
				tcpConn.Warn("[TCPConn.asyncSendCmd] " +
					"Channle已关闭，发送行为终止")
				break
			}
			tcpConn.sendMsgList(msg)
		default:
			waitting = false
			break
		}
	}

	return true
}

// 发送拼接消息 必须单线程执行
// sendMsgList 	tmsg 首消息，如果没有需要加入的第一个消息，直接给Nil即可
func (tcpConn *TCPConn) sendMsgList(tmsg *msg.MessageBinary) {
	// 开始拼包
	msglist := tcpConn.joinMsgByFunc(
		func(nowpkgsum int, nowpkglen int) *msg.MessageBinary {
			if nowpkgsum == 0 && tmsg != nil {
				// 如果这是第一个包，且包含首包
				return tmsg
			}
			if nowpkgsum >= MaxMsgPackSum {
				// 如果当前拼包消息数量已大到最大
				return nil
			}
			// 单次最大发送长度
			if tcpConn.sendBuffer.TotalSize() < msg.MessageMaxSize ||
				nowpkglen > tcpConn.sendBuffer.TotalSize()-msg.MessageMaxSize {
				// 超过最大限制长度，停止拼包
				return nil
			}
			// 遍历消息发送通道
			select {
			case msg, ok := <-tcpConn.sendmsgchan:
				// 取到了数据
				if msg == nil || !ok {
					// 通道中的数据不合法
					tcpConn.Warn("[TCPConn.sendMsgList] " +
						"Channle已关闭，发送行为终止")
					return nil
				}
				// 返回取到的消息
				return msg
			default:
				// 通道中没有数据了，停止拼包
				return nil
			}
		})
	// 拼包总消息长度
	nowpkgsum := len(msglist)
	if nowpkgsum == 0 {
		// 当前没有需要发送的消息
		return
	}

	bs, err := tcpConn.sendBuffer.SeekAll()
	if err != nil {
		tcpConn.Error("[TCPConn.sendMsgList] "+
			"tcpConn.sendBuffer.SeekAll() Err[%s]",
			err.Error())
	} else {
		secn, err := tcpConn.work.Write(bs)
		if err != nil {
			tcpConn.Debug("[TCPConn.sendMsgList] "+
				"缓冲区发送消息异常 Err[%s]",
				err.Error())
		} else {
			tcpConn.sendBuffer.MoveStart(secn)
			// 发送
			// 发送缓冲区长度减少
			atomic.AddInt64(&tcpConn.waitingSendBufferLength, int64(-secn))
		}
	}

	// 遍历已经发送的消息
	for _, msg := range msglist {
		// 调用发送回调函数
		msg.OnSendFinish()
		msg.Free()
	}
}

// 从指定接口中拼接消息 必须单线程执行
// 	回调参数：
// 		当前消息数量
// 		当前消息总大小
//
//  返回：
//  	拼接的	消息列表
//  			二进制列表
//  	总长度
// joinMsgByFunc  	最大延迟
func (tcpConn *TCPConn) joinMsgByFunc(getMsg func(int, int) *msg.MessageBinary) []*msg.MessageBinary {
	// 初始化变量
	var (
		nowpkgsum = int(0)
		nowpkglen = int(0)
	)
	for {
		msg := getMsg(nowpkgsum, nowpkglen)
		if msg == nil {
			break
		}
		// 拼接一个消息
		sendlen := int(msg.GetTotalLength())
		sendata := msg.GetBuffer()[:sendlen]
		err := tcpConn.sendBuffer.Write(sendata)
		if err != nil {
			tcpConn.Error("[TCPConn.joinMsgByFunc] "+
				"tcpConn.sendBuffer.Write(sendata) Err:%s", err.Error())
		}
		tcpConn.sendJoinedMessageBinaryBuffer[nowpkgsum] = msg
		nowpkgsum++
		nowpkglen += sendlen
	}

	return tcpConn.sendJoinedMessageBinaryBuffer[:nowpkgsum]
}

// recvThread 消息接收线程
func (tcpConn *TCPConn) recvThread() {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if stackInfo, err := sysutil.GetPanicInfo(recover()); err != nil {
			tcpConn.Error("[TCPConn.recvThread] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
		}
		close(tcpConn.recvmsgchan)
	}()
	if tcpConn.recvBuffer.TotalSize() == 0 {
		return
	}
	for true {
		if tcpConn.state != TCPConnStateLinked {
			return
		}
		// 设置阻塞过期时间
		derr := tcpConn.conn.SetReadDeadline(time.Now().Add(time.Duration(time.Millisecond * 250)))
		if derr != nil {
			// 设置阻塞过期时间失败
			return
		}
		// 从连接中读取数据
		_, err := tcpConn.recvBuffer.ReadFromReader()
		if err != nil {
			if err == io.EOF {
				// 连接关闭
				return
			}
			// 其他错误
			continue
		}
		// 循环读取当前缓冲区中的所有消息
		err = tcpConn.codec.RangeMsgBinary(tcpConn.recvBuffer,
			func(msgbinary *msg.MessageBinary) {
				// 解析消息
				tcpConn.recvmsgchan <- msgbinary
			})
		if err != nil {
			tcpConn.Error("[TCPConn.recvThread] "+
				"RangeMsgBinary读消息失败，断开连接 Err[%s]", err.Error())
			return
		}
	}
}
