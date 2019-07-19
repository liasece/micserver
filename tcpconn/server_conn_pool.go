package tcpconn

import (
	//	"os"
	"fmt"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/util"
	"math/rand"
	"net"
	"sync"
	"time"
)

type ServerConnPool struct {
	*log.Logger

	allsockets sync.Map // 所有连接
	linkSum    int32
	groupID    uint16

	TCPConnectPoolMutex sync.Mutex
	linkSumMutex        sync.Mutex
}

func (this *ServerConnPool) Init(groupID uint16) {
	this.groupID = groupID
}

func (this *ServerConnPool) NewServerConn(sctype TServerSCType,
	conn net.Conn,
	serverid string) *ServerConn {
	tcptask := NewServerConn(sctype, conn)
	if serverid == "" {
		this.AddAuto(tcptask)
	} else {
		this.Add(tcptask, serverid)
	}
	return tcptask
}

// 遍历连接池中的所有连接
func (this *ServerConnPool) Range(
	callback func(*ServerConn) bool) {
	this.allsockets.Range(func(tkey interface{},
		tvalue interface{}) bool {
		return callback(tvalue.(*ServerConn))
	})
}

// 将一条消息广播至指定类型的所有连接
func (this *ServerConnPool) BroadcastByType(servertype string,
	v msg.MsgStruct) {
	this.allsockets.Range(func(tkey interface{},
		tvalue interface{}) bool {
		value := tvalue.(*ServerConn)
		if util.GetServerIDType(value.Serverinfo.ServerID) == servertype {
			value.SendCmd(v)
		}
		return true
	})
}

// 广播消息
func (this *ServerConnPool) BroadcastCmd(v msg.MsgStruct) {
	this.allsockets.Range(func(tkey interface{},
		tvalue interface{}) bool {
		tvalue.(*ServerConn).SendCmd(v)
		return true
	})
}

// 获取指定类型负载最小的一个连接
func (this *ServerConnPool) GetMinClient(
	servertype string) *ServerConn {
	var jobnum uint32 = 0xFFFFFFFF
	var serverid uint64 = 0
	this.allsockets.Range(func(tkey interface{},
		tvalue interface{}) bool {
		value := tvalue.(*ServerConn)
		key := tkey.(uint64)
		if util.GetServerIDType(value.Serverinfo.ServerID) == servertype {
			if jobnum >= value.GetJobNum() {
				jobnum = value.GetJobNum()
				serverid = key
			}
		}
		return true
	})
	if serverid == 0 {
		return nil
	}
	if tcptask, found := this.allsockets.Load(serverid); found {
		return tcptask.(*ServerConn)
	}
	return nil
}

// 获取指定类型服务器的最新版本
func (this *ServerConnPool) GetLatestVersionByType(servertype string) uint64 {
	latestVersion := uint64(0)
	this.allsockets.Range(func(tkey interface{},
		tvalue interface{}) bool {
		value := tvalue.(*ServerConn)
		if util.GetServerIDType(value.Serverinfo.ServerID) == servertype &&
			value.Serverinfo.Version > latestVersion {
			latestVersion = value.Serverinfo.Version
		}
		return true
	})
	return latestVersion
}

// 获取指定类型负载最小的一个连接
func (this *ServerConnPool) GetMinClientLatestVersion(
	servertype string) *ServerConn {
	var jobnum uint32 = 0xFFFFFFFF
	var serverid uint64 = 0

	latestVersion := this.GetLatestVersionByType(servertype)

	this.allsockets.Range(func(tkey interface{},
		tvalue interface{}) bool {
		value := tvalue.(*ServerConn)
		key := tkey.(uint64)
		if util.GetServerIDType(value.Serverinfo.ServerID) == servertype &&
			value.Serverinfo.Version == latestVersion {
			if jobnum >= value.GetJobNum() {
				jobnum = value.GetJobNum()
				serverid = key
			}
		}
		return true
	})
	if serverid == 0 {
		return nil
	}
	if tcptask, found := this.allsockets.Load(serverid); found {
		return tcptask.(*ServerConn)
	}
	return nil
}

// 随机获取指定类型的一个连接
func (this *ServerConnPool) GetRandom(
	servertype string) *ServerConn {
	tasklist := make([]uint64, 0)
	this.allsockets.Range(func(tkey interface{},
		tvalue interface{}) bool {
		value := tvalue.(*ServerConn)
		key := tkey.(uint64)
		if util.GetServerIDType(value.Serverinfo.ServerID) == servertype {
			tasklist = append(tasklist, key)
		}
		return true
	})

	length := len(tasklist)
	if length > 0 {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		tmpindex := r.Intn(length)
		serverid := tasklist[tmpindex]

		if tcptask, found := this.allsockets.Load(serverid); found {
			return tcptask.(*ServerConn)
		}
	}
	return nil
}

// 根据连接的 TmpID 获取一个连接
func (this *ServerConnPool) Get(tempid string) *ServerConn {
	if tcptask, found := this.allsockets.Load(tempid); found {
		return tcptask.(*ServerConn)
	}
	return nil
}

func (this *ServerConnPool) Remove(tempid string) {
	if tvalue, found := this.allsockets.Load(tempid); found {
		value := tvalue.(*ServerConn)
		value.isAlive = false
		// 关闭消息发送协程
		value.Conn.Shutdown()
		// 删除连接
		this.remove(tempid)
		this.Debug("[ServerConnPool] 删除连接 TmpID[%s] 当前连接数量"+
			" Len[%d] ServerID[%s]",
			tempid, this.Len(), value.Serverinfo.ServerID)
		return
	}
}

func (this *ServerConnPool) Add(connct *ServerConn, tmpid string) {
	connct.Tempid = tmpid
	this.add(tmpid, connct)
}

func (this *ServerConnPool) AddAuto(connct *ServerConn) {
	tmpid, err := util.NewUniqueID(this.groupID)
	if err != nil {
		this.Error("[ServerConnPool.AddAuto] 生成UUID出错 Error[%s]",
			err.Error())
		return
	}
	connct.Tempid = tmpid
	this.add(connct.Tempid, connct)
}

// 修改链接的 tempip
func (this *ServerConnPool) ChangeTempid(tcptask *ServerConn,
	newtempid string) error {
	this.TCPConnectPoolMutex.Lock()
	defer this.TCPConnectPoolMutex.Unlock()
	if _, found := this.allsockets.Load(newtempid); found {
		return fmt.Errorf("目标连接已存在:%s", newtempid)
	}
	if ttcptask, found := this.allsockets.Load(tcptask.Tempid); found {
		tcptask := ttcptask.(*ServerConn)
		this.Debug("[ServerConnPool]修改连接tempid Old[%s] -->> New[%s]",
			tcptask.Tempid, newtempid)
		this.remove(tcptask.Tempid)
		tcptask.Tempid = newtempid
		this.add(tcptask.Tempid, tcptask)
	}
	return nil
}

func (this *ServerConnPool) Len() uint32 {
	if this.linkSum < 0 {
		return 0
	}
	return uint32(this.linkSum)
}

func (this *ServerConnPool) remove(tmpid string) {
	if _, ok := this.allsockets.Load(tmpid); !ok {
		return
	}
	// 删除连接
	this.allsockets.Delete(tmpid)
	this.linkSum--
}

func (this *ServerConnPool) add(tmpid string, value *ServerConn) {
	this.linkSumMutex.Lock()
	defer this.linkSumMutex.Unlock()
	if _, ok := this.allsockets.Load(tmpid); !ok {
		// 是新增的连接
		this.linkSum++
	}
	this.allsockets.Store(tmpid, value)
}
