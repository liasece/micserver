package tcpconn

import (
	//	"os"
	// "base"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/util"
	"math/rand"
	"net"
	// "servercomm"
	"errors"
	// "sync"
	"time"
)

const mClientConnPoolGroupSum = 10

type stringToClientConn struct {
	m *util.MapPool
}

func (this *stringToClientConn) Store(k string, v *ClientConn) {
	this.m.Push(util.GetStringHash(k)%mClientConnPoolGroupSum, k, v)
}

func (this *stringToClientConn) LoadOrStore(k string, v *ClientConn) (*ClientConn, bool) {
	vi, isLoad := this.m.LoadOfStroe(util.GetStringHash(k)%mClientConnPoolGroupSum, k, v)
	res := vi.(*ClientConn)
	return res, isLoad
}

func (this *stringToClientConn) Load(k string) (*ClientConn, bool) {
	tv, ok := this.m.Get(util.GetStringHash(k)%mClientConnPoolGroupSum, k)
	if tv == nil || !ok {
		return nil, false
	}
	return tv.(*ClientConn), true
}

func (this *stringToClientConn) Delete(k string) {
	this.m.Pop(util.GetStringHash(k)%mClientConnPoolGroupSum, k)
}

func (this *stringToClientConn) Range(callback func(string, *ClientConn)) {
	this.m.RangeAll(func(tk interface{}, tv interface{}) {
		if tk == nil || tv == nil {
			return
		}
		callback(tk.(string), tv.(*ClientConn))
	})
}

func (this *stringToClientConn) InitMapPool(gsum uint32) {
	this.m = util.NewMapPool(gsum)
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

func (this *ClientConnPool) NewClientConn(conn net.Conn) (*ClientConn, error) {
	tcptask := NewClientConn(conn)
	err := this.AddAuto(tcptask)
	if err != nil {
		return nil, err
	}
	return tcptask, nil
}

// 遍历连接池中的所有连接
func (this *ClientConnPool) Range(
	callback func(*ClientConn)) {
	this.allSockets.Range(func(key string,
		value *ClientConn) {
		callback(value)
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
		value.Debug("[ClientConnPool.Remove] 删除连接 当前连接数量"+
			" Len[%d]",
			this.Len())
		// 删除连接
		this.remove(tempid)
		return
	}
}

func (this *ClientConnPool) AddAuto(connct *ClientConn) error {
	tmpid, err := util.NewUniqueID(this.groupID)
	if err != nil {
		log.Error("[ClientConnPool.AddAuto] 生成UUID出错 Error[%s]",
			err.Error())
		return errors.New("unique id create error: " + err.Error())
	}
	connct.Tempid = tmpid
	this.add(connct.Tempid, connct)
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
	this.allSockets.Range(func(key string, value *ClientConn) {
		tasklist = append(tasklist, key)
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
