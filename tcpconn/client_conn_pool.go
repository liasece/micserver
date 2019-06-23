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

type uint64ToClientConn struct {
	m *util.MapPool
}

func (this *uint64ToClientConn) Init(gsum int) {
	this.m = util.NewMapPool(uint32(gsum))
}

func (this *uint64ToClientConn) Store(k uint64, v *ClientConn) {
	this.m.Push(uint32(k%mClientConnPoolGroupSum), k, v)
}

func (this *uint64ToClientConn) Load(k uint64) (*ClientConn, bool) {
	tv, ok := this.m.Get(uint32(k%mClientConnPoolGroupSum), k)
	if tv == nil || !ok {
		return nil, false
	}
	return tv.(*ClientConn), true
}

func (this *uint64ToClientConn) Delete(k uint64) {
	this.m.Pop(uint32(k%mClientConnPoolGroupSum), k)
}

func (this *uint64ToClientConn) Range(callback func(uint64, *ClientConn)) {
	this.m.RangeAll(func(tk interface{}, tv interface{}) {
		if tk == nil || tv == nil {
			return
		}
		callback(tk.(uint64), tv.(*ClientConn))
	})
}

type stringToClientConn struct {
	m *util.MapPool
}

func (this *stringToClientConn) Store(k string, v *ClientConn) {
	this.m.Push(util.GetStringHash(k)%mClientConnPoolGroupSum, k, v)
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

func (this *stringToClientConn) Init(gsum uint32) {
	this.m = util.NewMapPool(gsum)
}

type ClientConnPool struct {
	allsockets     uint64ToClientConn // 所有连接
	allopenidtasks stringToClientConn // 所有连接 by openid
	alluuidtasks   uint64ToClientConn // 所有连接 by uuid
	linkSum        int32
	groupID        uint16
}

func (this *ClientConnPool) Init(groupID int) {
	this.groupID = uint16(groupID)
	this.allsockets.Init(mClientConnPoolGroupSum)
	this.allopenidtasks.Init(mClientConnPoolGroupSum)
	this.alluuidtasks.Init(mClientConnPoolGroupSum)
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
	this.allsockets.Range(func(key uint64,
		value *ClientConn) {
		callback(value)
	})
}

// 根据连接的 TmpID 获取一个连接
func (this *ClientConnPool) Get(tempid uint64) *ClientConn {
	if tcptask, found := this.allsockets.Load(tempid); found {
		return tcptask
	}
	return nil
}

func (this *ClientConnPool) Remove(tempid uint64) {
	if value, found := this.allsockets.Load(tempid); found {
		if value.IsVertify() {
			if _, openidfound := this.allopenidtasks.
				Load(value.Openid); openidfound {
				this.allopenidtasks.Delete(value.Openid)
			}
			if _, uuidfound := this.alluuidtasks.Load(value.UUID); uuidfound {
				this.alluuidtasks.Delete(value.UUID)
			}
		}
		// 关闭消息发送协程
		value.Conn.Shutdown()
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

func (this *ClientConnPool) remove(tmpid uint64) {
	if _, ok := this.allsockets.Load(tmpid); !ok {
		return
	}
	// 删除连接
	this.allsockets.Delete(tmpid)
	this.linkSum--
}

func (this *ClientConnPool) add(tmpid uint64, value *ClientConn) {
	if _, ok := this.allsockets.Load(tmpid); !ok {
		// 是新增的连接
		this.linkSum++
	}
	this.allsockets.Store(tmpid, value)
}

// 根据 OpenID 索引 Task
func (this *ClientConnPool) AddTaskOpenID(
	task *ClientConn, openid string) {
	this.allopenidtasks.Store(openid, task)
}

// 随机获取指定类型的一个连接
func (this *ClientConnPool) GetRandom() *ClientConn {
	tasklist := make([]uint64, 0)
	this.allsockets.Range(func(key uint64, value *ClientConn) {
		tasklist = append(tasklist, key)
	})

	length := len(tasklist)
	if length > 0 {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		tmpindex := r.Intn(length)
		id := tasklist[tmpindex]

		if tcptask, found := this.allsockets.Load(id); found {
			return tcptask
		}
	}
	return nil
}

// 根据 OpenID 索引 Task
func (this *ClientConnPool) GetTaskByOpenID(
	openid string) *ClientConn {
	if oldtask, found := this.allopenidtasks.Load(openid); found {
		return oldtask
	}
	return nil
}

// 根据 UUID 索引 Task
func (this *ClientConnPool) AddTaskUUID(task *ClientConn,
	uuid uint64) {
	this.alluuidtasks.Store(uuid, task)
}

// 根据 UUID 索引 Task
func (this *ClientConnPool) GetTaskByUUID(
	uuid uint64) *ClientConn {
	if oldtask, found := this.alluuidtasks.Load(uuid); found {
		return oldtask
	}
	return nil
}
