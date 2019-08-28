package connect

import (
	"errors"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/util"
	"math/rand"
	"net"
	"time"
)

const mClientConnPoolGroupSum = 10

type stringToClientConn struct {
	*util.MapPool
}

// 加载或存储
func (this *stringToClientConn) LoadOrStore(k string,
	v *ClientConn) (*ClientConn, bool) {
	vi, isLoad := this.MapPool.LoadOrStore(k, v)
	res := vi.(*ClientConn)
	return res, isLoad
}

// 加载
func (this *stringToClientConn) Load(k string) (*ClientConn, bool) {
	tv, ok := this.MapPool.Load(k)
	if tv == nil || !ok {
		return nil, false
	}
	return tv.(*ClientConn), true
}

// 遍历所有
func (this *stringToClientConn) Range(callback func(string, *ClientConn) bool) {
	this.MapPool.RangeAll(func(tk interface{}, tv interface{}) bool {
		if tk == nil || tv == nil {
			return true
		}
		if !callback(tk.(string), tv.(*ClientConn)) {
			return false
		}
		return true
	})
}

// 初始化
func (this *stringToClientConn) InitMapPool(gsum uint32) {
	this.MapPool = util.NewMapPool(gsum)
}

type ClientConnPool struct {
	allSockets stringToClientConn // 所有连接
	linkSum    int32
	groupID    uint16
}

func (this *ClientConnPool) Init(groupID int32) {
	this.groupID = uint16(groupID)
	this.allSockets.InitMapPool(mClientConnPoolGroupSum)
}

func (this *ClientConnPool) NewClientConn(conn net.Conn,
	onRecv func(*ClientConn, *msg.MessageBinary),
	onClose func(*ClientConn)) (*ClientConn, error) {
	tcptask := NewClientConn(conn, onRecv, onClose)
	err := this.AddAuto(tcptask)
	if err != nil {
		return nil, err
	}
	return tcptask, nil
}

// 遍历连接池中的所有连接
func (this *ClientConnPool) Range(
	callback func(*ClientConn) bool) {
	this.allSockets.Range(func(key string,
		value *ClientConn) bool {
		if !callback(value) {
			return false
		}
		return true
	})
}

// 根据连接的 TmpID 获取一个连接
func (this *ClientConnPool) Get(tempid string) *ClientConn {
	if tcptask, found := this.allSockets.Load(tempid); found {
		return tcptask
	}
	return nil
}

func (this *ClientConnPool) Remove(tempid string) {
	if value, found := this.allSockets.Load(tempid); found {
		// 关闭消息发送协程
		value.Shutdown()
		// 删除连接
		this.remove(tempid)
		value.Debug("[ClientConnPool.Remove] 删除连接 当前连接数量"+
			" Len[%d]",
			this.Len())
		return
	}
}

func (this *ClientConnPool) AddAuto(connct *ClientConn) error {
	tmpid, err := util.NewUniqueID(this.groupID)
	if err != nil {
		log.Error("[ClientConnPool.AddAuto] 生成ConnectID出错 Error[%s]",
			err.Error())
		return errors.New("unique id create error: " + err.Error())
	}
	connct.SetConnectID(tmpid)
	this.add(connct.GetConnectID(), connct)
	return nil
}

func (this *ClientConnPool) Len() uint32 {
	if this.linkSum < 0 {
		return 0
	}
	return uint32(this.linkSum)
}

func (this *ClientConnPool) remove(tmpid string) {
	if _, ok := this.allSockets.Load(tmpid); !ok {
		return
	}
	// 删除连接
	this.allSockets.Delete(tmpid)
	this.linkSum--
}

func (this *ClientConnPool) add(tmpid string, value *ClientConn) {
	_, isLoad := this.allSockets.LoadOrStore(tmpid, value)
	if !isLoad {
		this.linkSum++
	} else {
		this.allSockets.Store(tmpid, value)
	}
}

// 随机获取指定类型的一个连接
func (this *ClientConnPool) GetRandom() *ClientConn {
	tasklist := make([]string, 0)
	this.allSockets.Range(func(key string, value *ClientConn) bool {
		tasklist = append(tasklist, key)
		return true
	})

	length := len(tasklist)
	if length > 0 {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		tmpindex := r.Intn(length)
		id := tasklist[tmpindex]

		if tcptask, found := this.allSockets.Load(id); found {
			return tcptask
		}
	}
	return nil
}
