package connect

import (
	"math/rand"
	"net"

	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/util/pool"
)

// ClientPool 客户端连接池
type ClientPool struct {
	*log.Logger
	allSockets pool.MapPool // 所有连接
	linkSum    int32
}

// Init 初始化Clieng连接池
func (cp *ClientPool) Init() {
	cp.allSockets.Init(mClientPoolGroupSum)
}

// SetLogger 设置客户端连接池的Logger
func (cp *ClientPool) SetLogger(l *log.Logger) {
	cp.Logger = l
}

// NewTCPClient 使用TCP连接新建一个Client
func (cp *ClientPool) NewTCPClient(conn net.Conn,
	connHook IConnectHook) (*Client, error) {
	client := &Client{}
	client.SetLogger(cp.Logger)
	client.InitTCP(conn, connHook)
	cp.Add(client)
	return client, nil
}

// LoadOrStore 加载或存储一个客户端连接
func (cp *ClientPool) LoadOrStore(k string,
	v *Client) (*Client, bool) {
	vi, isLoad := cp.allSockets.LoadOrStore(k, v)
	res := vi.(*Client)
	return res, isLoad
}

// Get 根据连接的 TmpID 获取连接
func (cp *ClientPool) Get(tempid string) *Client {
	if vi, ok := cp.allSockets.Load(tempid); ok {
		return vi.(*Client)
	}
	return nil
}

// Len 当前连接池中的连接数量
func (cp *ClientPool) Len() uint32 {
	if cp.linkSum < 0 {
		return 0
	}
	return uint32(cp.linkSum)
}

// remove 根据连接的 TmpID 从连接池移除一个连接
func (cp *ClientPool) remove(tmpid string) {
	if _, ok := cp.allSockets.Load(tmpid); !ok {
		return
	}
	// 删除连接
	cp.allSockets.Delete(tmpid)
	cp.linkSum--
}

// Add 增加一个连接到连接池中
func (cp *ClientPool) Add(client *Client) {
	tmpid := client.GetTempID()
	_, isLoad := cp.allSockets.LoadOrStore(tmpid, client)
	if !isLoad {
		cp.linkSum++
	} else {
		cp.allSockets.Store(tmpid, client)
	}
}

// Range 遍历连接池中的所有连接，如果 cb() 返回 false 则中止遍历
func (cp *ClientPool) Range(
	cb func(*Client) bool) {
	cp.allSockets.Range(func(tk, tv interface{}) bool {
		if tk == nil || tv == nil {
			return true
		}
		if !cb(tv.(*Client)) {
			return false
		}
		return true
	})
}

// Remove 根据连接的 TmpID 从连接池移除一个连接
func (cp *ClientPool) Remove(tempid string) {
	if vi, ok := cp.allSockets.Load(tempid); ok {
		client := vi.(*Client)
		// 关闭消息发送协程
		client.Shutdown()
		// 删除连接
		cp.remove(tempid)
		client.Debug("[ClientPool.Remove] 断开连接 TmpID[%s] 当前连接数量"+
			" Len[%d]",
			tempid, cp.Len())
		return
	}
}

// GetRandom 随机获取连接池中的一个连接
func (cp *ClientPool) GetRandom() *Client {
	tasklist := make([]string, 0)
	cp.Range(func(client *Client) bool {
		tasklist = append(tasklist, client.GetTempID())
		return true
	})

	length := len(tasklist)
	if length > 0 {
		return cp.Get(tasklist[rand.Intn(length)])
	}
	return nil
}
