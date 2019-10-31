package connect

import (
	"errors"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/util"
	"math/rand"
	"net"
	"time"
)

const mClientPoolGroupSum = 10

type stringToClient struct {
	util.MapPool
}

// 加载或存储
func (this *stringToClient) LoadOrStore(k string,
	v *Client) (*Client, bool) {
	vi, isLoad := this.MapPool.LoadOrStore(k, v)
	res := vi.(*Client)
	return res, isLoad
}

// 加载
func (this *stringToClient) Load(k string) (*Client, bool) {
	tv, ok := this.MapPool.Load(k)
	if tv == nil || !ok {
		return nil, false
	}
	return tv.(*Client), true
}

// 遍历所有
func (this *stringToClient) Range(callback func(string, *Client) bool) {
	this.MapPool.Range(func(tk interface{}, tv interface{}) bool {
		if tk == nil || tv == nil {
			return true
		}
		if !callback(tk.(string), tv.(*Client)) {
			return false
		}
		return true
	})
}

// 初始化
func (this *stringToClient) InitMapPool(gsum int) {
	this.MapPool.Init(gsum)
}

type ClientPool struct {
	*log.Logger
	allSockets stringToClient // 所有连接
	linkSum    int32
	groupID    uint16
}

func (this *ClientPool) Init(groupID int32) {
	this.groupID = uint16(groupID)
	this.allSockets.InitMapPool(mClientPoolGroupSum)
}

func (this *ClientPool) SetLogger(l *log.Logger) {
	this.Logger = l
}

func (this *ClientPool) NewTCPClient(conn net.Conn,
	connHook ConnectHook) (*Client, error) {
	tcptask := &Client{}
	tcptask.SetLogger(this.Logger)
	tcptask.InitTCP(conn, connHook)
	err := this.AddAuto(tcptask)
	if err != nil {
		return nil, err
	}
	return tcptask, nil
}

// 遍历连接池中的所有连接
func (this *ClientPool) Range(
	callback func(*Client) bool) {
	this.allSockets.Range(func(key string,
		value *Client) bool {
		if !callback(value) {
			return false
		}
		return true
	})
}

// 根据连接的 TmpID 获取一个连接
func (this *ClientPool) Get(tempid string) *Client {
	if tcptask, found := this.allSockets.Load(tempid); found {
		return tcptask
	}
	return nil
}

func (this *ClientPool) Remove(tempid string) {
	if value, found := this.allSockets.Load(tempid); found {
		// 关闭消息发送协程
		value.Shutdown()
		// 删除连接
		this.remove(tempid)
		value.Debug("[ClientPool.Remove] 删除连接 当前连接数量"+
			" Len[%d]",
			this.Len())
		return
	}
}

func (this *ClientPool) AddAuto(connct *Client) error {
	tmpid, err := util.NewUniqueID(this.groupID)
	if err != nil {
		log.Error("[ClientPool.AddAuto] 生成ConnectID出错 Error[%s]",
			err.Error())
		return errors.New("unique id create error: " + err.Error())
	}
	connct.SetConnectID(tmpid)
	this.add(connct.GetConnectID(), connct)
	return nil
}

func (this *ClientPool) Len() uint32 {
	if this.linkSum < 0 {
		return 0
	}
	return uint32(this.linkSum)
}

func (this *ClientPool) remove(tmpid string) {
	if _, ok := this.allSockets.Load(tmpid); !ok {
		return
	}
	// 删除连接
	this.allSockets.Delete(tmpid)
	this.linkSum--
}

func (this *ClientPool) add(tmpid string, value *Client) {
	_, isLoad := this.allSockets.LoadOrStore(tmpid, value)
	if !isLoad {
		this.linkSum++
	} else {
		this.allSockets.Store(tmpid, value)
	}
}

// 随机获取指定类型的一个连接
func (this *ClientPool) GetRandom() *Client {
	tasklist := make([]string, 0)
	this.allSockets.Range(func(key string, value *Client) bool {
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
