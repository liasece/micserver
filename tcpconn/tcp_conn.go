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
	"github.com/liasece/micserver/util"
	// "github.com/liasece/micserver/functime"
	"bytes"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/msg"
	// "encoding/hex"
	// "fmt"
	"net"
	// "runtime"
	// "errors"
	"io/ioutil"
	"sync/atomic"
	"time"
)

const MsgMaxSize = 64 * 1024
const MaxMsgPackSum = 5000

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
	net.Conn
	sendmsgchan chan *msg.MessageBinary

	// 尝试关闭一个连接
	shutdownChan chan struct{}
	stopChan     chan struct{}
	state        int32

	// 发送缓冲区大小
	sendBufferSize           int
	maxWaitSendMsgBufferSize int

	// 当前已发送数据长度
	nowSendBufferLength int64
}

// 初始化一个TCPConn对象
// 	conn: net.Conn对象
// 	msgBufferSize: 消息Channel缓冲区数量，不宜过大
// 	sendBufferSize: 拼包发送缓冲区大小
// 	maxWaitSendMsgBufferSize: 等待队列中的消息缓冲区大小
func (this *TCPConn) Init(conn net.Conn,
	msgBufferSize uint32, sendBufferSize int,
	maxWaitSendMsgBufferSize int) {
	this.sendmsgchan = make(chan *msg.MessageBinary, msgBufferSize)
	this.shutdownChan = make(chan struct{})
	this.stopChan = make(chan struct{})
	this.Conn = conn
	this.sendBufferSize = sendBufferSize
	this.maxWaitSendMsgBufferSize = maxWaitSendMsgBufferSize
	this.state = TCPCONNSTATE_LINKED

	go this.sendThread()
}

func (this *TCPConn) IsAlive() bool {
	if this.state == TCPCONNSTATE_LINKED {
		return true
	}
	return false
}

// 尝试关闭此连接
func (this *TCPConn) Shutdown() {
	if this.state < TCPCONNSTATE_LINKED {
		return
	}
	sec := atomic.CompareAndSwapInt32(&this.state,
		TCPCONNSTATE_LINKED, TCPCONNSTATE_HOLD)
	if sec {
		go this.shutdownThread()
	}
}

// 异步 停止一个TCP连接的发送过程
// 异步是为了推迟实际关闭连接的时间，给处理剩余未发送的消息的时间
func (this *TCPConn) shutdownThread() {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := util.GetPanicInfo(recover()); err != nil {
			log.Warn("[TCPConn.shutdownThread] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
		}
	}()
	// 延迟两秒发送，否则消息可能处理不完
	// time.Sleep(2 * time.Second)
	close(this.shutdownChan)
}

//关闭socket
func (this *TCPConn) closeSocket() error {
	this.state = TCPCONNSTATE_CLOSED
	return this.Conn.Close()
}

func (this *TCPConn) ReadAll() ([]byte, error) {
	return ioutil.ReadAll(this.Conn)
}

// 异步发送一条消息，不带发送完成回调
func (this *TCPConn) SendCmd(v msg.MsgStruct,
	encryption msg.TEncryptionType) error {
	if this.state >= TCPCONNSTATE_HOLD {
		log.Warn("[TCPConn.SendCmd] 连接已失效，取消发送")
		return ErrCloseed
	}
	msg := msg.MakeMessageByJson(v)
	if encryption != 0 && msg != nil {
		msg.Encryption(encryption)
	}
	return this.SendMessageBinary(msg)
}

// 异步发送一条消息，带发送完成回调
func (this *TCPConn) SendCmdWithCallback(v msg.MsgStruct,
	callback func(interface{}), cbarg interface{},
	encryption msg.TEncryptionType) error {
	if this.state >= TCPCONNSTATE_HOLD {
		log.Warn("[TCPConn.SendCmdWithCallback] 连接已失效，取消发送")
		return ErrCloseed
	}
	msg := msg.MakeMessageByJson(v)

	msg.OnSendDone = callback
	msg.OnSendDoneArg = cbarg

	if encryption != 0 && msg != nil {
		msg.Encryption(encryption)
	}
	return this.SendMessageBinary(msg)
}

// 发送 Bytes
func (this *TCPConn) SendBytes(
	cmdid uint16, protodata []byte, encryption msg.TEncryptionType) error {
	if this.state >= TCPCONNSTATE_HOLD {
		log.Warn("[TCPConn.SendBytes] 连接已失效，取消发送")
		return ErrCloseed
	}
	msgbinary := msg.MakeMessageByBytes(cmdid, protodata)
	if encryption != 0 && msgbinary != nil {
		msgbinary.Encryption(encryption)
	}
	return this.SendMessageBinary(msgbinary)
}

// 发送 MsgBinary
func (this *TCPConn) SendMessageBinary(
	msgbinary *msg.MessageBinary) error {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := util.GetPanicInfo(recover()); err != nil {
			log.Warn("[TCPConn.SendMessageBinary] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
		}
	}()
	// 检查连接是否已死亡
	if this.state >= TCPCONNSTATE_HOLD {
		log.Warn("[TCPConn.SendMessageBinary] 连接已失效，取消发送")
		return ErrCloseed
	}
	// 如果发送数据为空
	if msgbinary == nil {
		log.Debug("[TCPConn.SendMessageBinary] 发送消息为空，取消发送")
		return ErrSendNilData
	}

	// 检查发送channel是否已经关闭
	select {
	case <-this.stopChan:
		log.Warn("[TCPConn.SendMessageBinary] 发送Channel已关闭，取消发送")
		return ErrCloseed
	default:
	}

	// 检查等待缓冲区数据是否已满
	if this.nowSendBufferLength > int64(this.maxWaitSendMsgBufferSize) {
		log.Warn("[TCPConn.SendMessageBinary] 等待发送缓冲区满")
		return ErrBufferFull
	}

	// 确认发送channel是否已经关闭
	select {
	case <-this.stopChan:
		log.Warn("[TCPConn.SendMessageBinary] 发送Channel已关闭，取消发送")
		return ErrCloseed
	case this.sendmsgchan <- msgbinary:
		atomic.AddInt64(&this.nowSendBufferLength, int64(msgbinary.CmdLen))
	default:
		log.Warn("[TCPConn.SendMessageBinary] 发送Channel缓冲区满，阻塞超时")
		return ErrBufferFull
	}
	return nil
}

// 消息发送进程
func (this *TCPConn) sendThread() {
	for {
		if this.asyncSendCmd() {
			// 正常退出
			break
		}
	}
	// 用于通知发送线程，发送channel已关闭
	log.Debug("[TCPConn.sendThread] 发送线程已关闭")
	close(this.stopChan)
	err := this.closeSocket()
	if err != nil {
		log.Error("[TCPConn.sendThread] closeSocket Err[%s]",
			err.Error())
	}
}

func (this *TCPConn) asyncSendCmd() (normalreturn bool) {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := util.GetPanicInfo(recover()); err != nil {
			log.Error("[TCPConn.asyncSendCmd] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
			normalreturn = false
		}
	}()

	isrunning := true
	for isrunning {
		select {
		case msg, ok := <-this.sendmsgchan:
			if msg == nil || !ok {
				log.Warn("[TCPConn.asyncSendCmd] " +
					"Channle已关闭，发送行为终止")
				break
			}
			this.sendMsgList(msg)
		case <-this.shutdownChan:
			isrunning = false
			break
		}
	}

	waitting := true
	for waitting {
		select {
		case msg, ok := <-this.sendmsgchan:
			if msg == nil || !ok {
				log.Warn("[TCPConn.asyncSendCmd] " +
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

// 发送拼接消息
// 	tmsg 首消息，如果没有需要加入的第一个消息，直接给Nil即可
func (this *TCPConn) sendMsgList(tmsg *msg.MessageBinary) {
	// 开始拼包
	msglist, buflist, nowpkglen, maxlow := this.joinMsgByFunc(
		func(nowpkgsum uint32, nowpkglen int) *msg.MessageBinary {
			if nowpkgsum == 0 && tmsg != nil {
				// 如果这是第一个包，且包含首包
				return tmsg
			}
			if nowpkgsum >= MaxMsgPackSum {
				// 如果当前拼包消息数量已大到最大
				return nil
			}
			// 单次最大发送长度
			if this.sendBufferSize < MsgMaxSize ||
				nowpkglen > this.sendBufferSize-MsgMaxSize {
				// 超过最大限制长度，停止拼包
				return nil
			}
			// 遍历消息发送通道
			select {
			case msg, ok := <-this.sendmsgchan:
				// 取到了数据
				if msg == nil || !ok {
					// 通道中的数据不合法
					log.Warn("[TCPConn.asyncSendCmd] " +
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
	nowpkgsum := uint32(len(msglist))
	if nowpkgsum == 0 {
		// 当前没有需要发送的消息
		return
	}

	// 拼包完成
	// 发送缓冲区长度减少
	atomic.AddInt64(&this.nowSendBufferLength, int64(-nowpkglen))

	// 检查消息包的超时时间
	if maxlow > 2 {
		// 发送时间超时
		// 由于可能时客户端网络原因，只需要Info等级log
		log.Info("[TCPConn.asyncSendCmd] "+
			"发送消息延迟[%d]s PkgSum[%d] AllSize[%d]",
			maxlow, nowpkgsum, nowpkglen)
	}

	msgbuf := bytes.NewBuffer(make([]byte, 0, nowpkglen))
	// 此处存在大量拷贝，待优化
	for _, buf := range buflist {
		msgbuf.Write(buf)
	}

	_, err := this.Conn.Write(msgbuf.Bytes())
	if err != nil {
		log.Warn("[TCPConn.asyncSendCmd] "+
			"缓冲区发送消息异常 Err[%s]",
			err.Error())
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

// 从指定接口中拼接消息
// 	回调参数：
// 		当前消息数量
// 		当前消息总大小
//
//  返回：
//  	拼接的	消息列表
//  			二进制列表
//  	总长度
//  	最大延迟
func (this *TCPConn) joinMsgByFunc(getMsg func(uint32, int) *msg.MessageBinary) (
	[]*msg.MessageBinary, [][]byte, int, uint32) {
	// 初始化变量
	var (
		buflist   = make([][]byte, 0)
		msglist   = make([]*msg.MessageBinary, 0)
		nowpkgsum = uint32(0)
		nowpkglen = int(0)
		curtime   = uint32(time.Now().Unix())
		maxlow    = uint32(0)
	)
	for {
		msg := getMsg(nowpkgsum, nowpkglen)
		if msg == nil {
			break
		}
		// 拼接一个消息

		// 用于计算发送延迟
		tmplow := curtime - msg.TimeStamp
		if tmplow > maxlow {
			maxlow = tmplow
		}

		sendata, sendlen := msg.WriteBinary()

		buflist = append(buflist, sendata)
		msglist = append(msglist, msg)

		nowpkgsum++
		nowpkglen += sendlen
	}

	return msglist, buflist, nowpkglen, maxlow
}
