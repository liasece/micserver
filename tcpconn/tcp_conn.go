/**
 * \file TCPConn.go
 * \version
 * \author liaojiansheng
 * \date  2019年01月31日 12:22:43
 * \brief 连接数据管理器
 *
 */

package tcpconn

import (
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/util"
	"io"
	"io/ioutil"
	"net"
	"sync/atomic"
	"time"
)

// 消息合批时，合并的最大消息数量
const (
	MaxMsgPackSum = 5000
)

// TCPConn 的连接状态枚举
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

type TCPConn struct {
	*log.Logger
	net.Conn

	// 尝试关闭一个连接
	shutdownChan chan struct{}
	state        int32

	// 发送等待通道
	sendmsgchan chan *msg.MessageBinary
	// 发送缓冲区
	sendBuffer *util.IOBuffer
	// 当前等待发送出去的数据总大小
	waitingSendBufferLength int64
	// 等待发送数据的总大小
	maxWaitingSendBufferLength int

	// 接收等待通道
	recvmsgchan chan *msg.MessageBinary
	// 接收缓冲区
	recvBuffer *util.IOBuffer
}

// 初始化一个TCPConn对象
// 	conn: net.Conn对象
// 	sendChanSize: 	发送等待队列中的消息缓冲区大小
// 	sendBufferSize: 发送拼包发送缓冲区大小
// 	recvChanSize: 	接收等待队列中的消息缓冲区大小
// 	recvBufferSize: 接收拼包发送缓冲区大小
// 返回：接收到的 messagebinary 的对象 chan
func (this *TCPConn) Init(conn net.Conn,
	sendChanSize int, sendBufferSize int,
	recvChanSize int, recvBufferSize int) chan *msg.MessageBinary {
	this.shutdownChan = make(chan struct{})
	this.Conn = conn
	this.state = TCPCONNSTATE_LINKED

	// 发送
	this.sendmsgchan = make(chan *msg.MessageBinary, sendChanSize)
	this.maxWaitingSendBufferLength = msg.MessageMaxSize * sendChanSize
	this.sendBuffer = util.NewIOBuffer(nil, sendBufferSize)
	go this.sendThread()

	// 接收
	this.recvmsgchan = make(chan *msg.MessageBinary, recvChanSize)
	this.recvBuffer = util.NewIOBuffer(conn, recvBufferSize)
	go this.recvThread()

	return this.recvmsgchan
}

func (this *TCPConn) IsAlive() bool {
	if this.state == TCPCONNSTATE_LINKED {
		return true
	}
	return false
}

// 尝试关闭此连接
func (this *TCPConn) Shutdown() {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := util.GetPanicInfo(recover()); err != nil {
			this.Warn("[TCPConn.shutdownThread] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
		}
	}()
	if this.state != TCPCONNSTATE_LINKED {
		return
	}
	sec := atomic.CompareAndSwapInt32(&this.state,
		TCPCONNSTATE_LINKED, TCPCONNSTATE_HOLD)
	if sec {
		close(this.shutdownChan)
	}
}

func (this *TCPConn) ReadAll() ([]byte, error) {
	return ioutil.ReadAll(this.Conn)
}

// 异步发送一条消息，不带发送完成回调
func (this *TCPConn) SendCmd(v msg.MsgStruct) error {
	if this.state >= TCPCONNSTATE_HOLD {
		this.Warn("[TCPConn.SendCmd] 连接已被关闭，取消发送 Msg[%s]",
			v.GetMsgName())
		return ErrCloseed
	}

	msg := msg.MakeMessageByJson(v)

	return this.SendMessageBinary(msg)
}

// 异步发送一条消息，带发送完成回调
func (this *TCPConn) SendCmdWithCallback(v msg.MsgStruct,
	callback func(interface{}), cbarg interface{}) error {
	if this.state >= TCPCONNSTATE_HOLD {
		this.Warn("[TCPConn.SendCmdWithCallback] 连接已失效，取消发送")
		return ErrCloseed
	}
	msg := msg.MakeMessageByJson(v)
	msg.OnSendDone = callback
	msg.OnSendDoneArg = cbarg

	return this.SendMessageBinary(msg)
}

// 发送 Bytes
func (this *TCPConn) SendBytes(
	cmdid uint16, protodata []byte) error {
	if this.state >= TCPCONNSTATE_HOLD {
		this.Warn("[TCPConn.SendBytes] 连接已失效，取消发送")
		return ErrCloseed
	}
	msgbinary := msg.MakeMessageByBytes(cmdid, protodata)

	return this.SendMessageBinary(msgbinary)
}

// 发送 MsgBinary
func (this *TCPConn) SendMessageBinary(
	msgbinary *msg.MessageBinary) error {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := util.GetPanicInfo(recover()); err != nil {
			this.Warn("[TCPConn.SendMessageBinary] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
		}
	}()
	// 检查连接是否已死亡
	if this.state >= TCPCONNSTATE_HOLD {
		this.Warn("[TCPConn.SendMessageBinary] 连接已失效，取消发送")
		return ErrCloseed
	}
	// 如果发送数据为空
	if msgbinary == nil {
		this.Debug("[TCPConn.SendMessageBinary] 发送消息为空，取消发送")
		return ErrSendNilData
	}

	// 检查发送channel是否已经关闭
	select {
	case <-this.shutdownChan:
		this.Warn("[TCPConn.SendMessageBinary] 发送Channel已关闭，取消发送")
		return ErrCloseed
	default:
	}

	// 检查等待缓冲区数据是否已满
	if this.waitingSendBufferLength > int64(this.maxWaitingSendBufferLength) {
		this.Error("[TCPConn.SendMessageBinary] 等待发送缓冲区满")
		return ErrBufferFull
	}

	// 确认发送channel是否已经关闭
	select {
	case <-this.shutdownChan:
		this.Warn("[TCPConn.SendMessageBinary] 发送Channel已关闭，取消发送")
		return ErrCloseed
	case this.sendmsgchan <- msgbinary:
		atomic.AddInt64(&this.waitingSendBufferLength, int64(msgbinary.CmdLen))
	default:
		this.Warn("[TCPConn.SendMessageBinary] 发送Channel缓冲区满，阻塞超时")
		return ErrBufferFull
	}
	return nil
}

// 消息发送线程 必须单线程执行
func (this *TCPConn) sendThread() {
	for {
		if this.asyncSendCmd() {
			// 正常退出
			break
		}
	}
	// 用于通知发送线程，发送channel已关闭
	this.Debug("[TCPConn.sendThread] 发送线程已关闭")
	// close(this.stopChan)
	err := this.closeSocket()
	if err != nil {
		this.Error("[TCPConn.sendThread] closeSocket Err[%s]",
			err.Error())
	}
}

//关闭socket 应该在消息尝试发送完之后执行
func (this *TCPConn) closeSocket() error {
	this.state = TCPCONNSTATE_CLOSED
	return this.Conn.Close()
}

// 异步方式发送消息 必须单线程执行
func (this *TCPConn) asyncSendCmd() (normalreturn bool) {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := util.GetPanicInfo(recover()); err != nil {
			this.Error("[TCPConn.asyncSendCmd] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
			normalreturn = false
		}
	}()

	isrunning := true
	for isrunning {
		select {
		case msg, ok := <-this.sendmsgchan:
			if msg == nil || !ok {
				this.Warn("[TCPConn.asyncSendCmd] " +
					"Channle已关闭，发送行为终止")
				break
			}
			this.sendMsgList(msg)
		case <-this.shutdownChan:
			// 线程被主动关闭
			isrunning = false
			break
		}
	}

	// 线程准备退出，执行收尾工作，尝试将未发送的消息发送出去
	waitting := true
	for waitting {
		select {
		case msg, ok := <-this.sendmsgchan:
			// 从发送chan中获取一条消息
			if msg == nil || !ok {
				this.Warn("[TCPConn.asyncSendCmd] " +
					"Channle已关闭，发送行为终止")
				break
			}
			this.sendMsgList(msg)
		default:
			waitting = false
			break
		}
	}

	return true
}

// 发送拼接消息 必须单线程执行
// 	tmsg 首消息，如果没有需要加入的第一个消息，直接给Nil即可
func (this *TCPConn) sendMsgList(tmsg *msg.MessageBinary) {
	// 开始拼包
	msglist := this.joinMsgByFunc(
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
			if this.sendBuffer.TotalSize() < msg.MessageMaxSize ||
				nowpkglen > this.sendBuffer.TotalSize()-msg.MessageMaxSize {
				// 超过最大限制长度，停止拼包
				return nil
			}
			// 遍历消息发送通道
			select {
			case msg, ok := <-this.sendmsgchan:
				// 取到了数据
				if msg == nil || !ok {
					// 通道中的数据不合法
					this.Warn("[TCPConn.sendMsgList] " +
						"Channle已关闭，发送行为终止")
					return nil
				}
				// 返回取到的消息
				return msg
			default:
				// 通道中没有数据了，停止拼包
				return nil
			}
			return nil
		})
	// 拼包总消息长度
	nowpkgsum := len(msglist)
	if nowpkgsum == 0 {
		// 当前没有需要发送的消息
		return
	}

	bs, err := this.sendBuffer.SeekAll()
	if err != nil {
		this.Error("[TCPConn.sendMsgList] "+
			"this.sendBuffer.SeekAll() Err[%s]",
			err.Error())
	} else {
		secn, err := this.Conn.Write(bs)
		if err != nil {
			this.Warn("[TCPConn.sendMsgList] "+
				"缓冲区发送消息异常 Err[%s]",
				err.Error())
		} else {
			this.sendBuffer.MoveStart(secn)
			// 发送
			// 发送缓冲区长度减少
			atomic.AddInt64(&this.waitingSendBufferLength, int64(-secn))
		}
	}

	// 遍历已经发送的消息
	for _, msg := range msglist {
		// 调用发送回调函数
		if msg.OnSendDone != nil {
			msg.OnSendDone(msg.OnSendDoneArg)
		}
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
//  	最大延迟
func (this *TCPConn) joinMsgByFunc(getMsg func(int, int) *msg.MessageBinary) []*msg.MessageBinary {
	// 初始化变量
	var (
		msglist   = make([]*msg.MessageBinary, 0)
		nowpkgsum = int(0)
		nowpkglen = int(0)
	)
	for {
		msg := getMsg(nowpkgsum, nowpkglen)
		if msg == nil {
			break
		}
		// 拼接一个消息

		sendata, sendlen := msg.WriteBinary()
		err := this.sendBuffer.Write(sendata)
		if err != nil {
			this.Error("[TCPConn.joinMsgByFunc] "+
				"this.sendBuffer.Write(sendata) Err:%s", err.Error())
		}
		msglist = append(msglist, msg)

		nowpkgsum++
		nowpkglen += sendlen
	}

	return msglist
}

func (this *TCPConn) recvThread() {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := util.GetPanicInfo(recover()); err != nil {
			this.Error("[TCPConn.recvThread] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
		}
		close(this.recvmsgchan)
	}()
	if this.recvBuffer.TotalSize() == 0 {
		return
	}
	// 消息缓冲
	msgReader := msg.NewMessageBinaryReader(this.recvBuffer)

	for true {
		if this.state != TCPCONNSTATE_LINKED {
			return
		}
		// 设置阻塞过期时间
		derr := this.SetReadDeadline(time.Now().
			Add(time.Duration(time.Millisecond * 250)))
		if derr != nil {
			// 设置阻塞过期时间失败
			return
		}
		// 从连接中读取数据
		_, err := this.recvBuffer.ReadFromReader()
		if err != nil {
			if err == io.EOF {
				// 连接关闭
				return
			} else {
				// 其他错误
				continue
			}
		}
		// 循环读取当前缓冲区中的所有消息
		err = msgReader.RangeMsgBinary(func(msgbinary *msg.MessageBinary) {
			// 解析消息
			this.recvmsgchan <- msgbinary
		})
		if err != nil {
			this.Error("[TCPConn.recvThread] "+
				"RangeMsgBinary读消息失败，断开连接 Err[%s]", err.Error())
			return
		}
	}
}
