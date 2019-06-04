package tcpconn

import (
	//	"os"
	"github.com/liasece/micserver"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/util"
	"math/rand"
	"net"
	// "servercomm"
	"sync"
	"time"
)

type ServerConnPool struct {
	sctype TServerSCType // 连接池 服务器/客户端 类型

	allsockets sync.Map // 所有连接
	linkSum    int32
	groupID    uint16

	TCPConnectPoolMutex sync.Mutex
	linkSumMutex        sync.Mutex
}

func (this *ServerConnPool) Init(sctype TServerSCType, groupID uint16) {
	this.groupID = groupID
}

func (this *ServerConnPool) NewServerConn(conn net.Conn,
	serverid uint64) *ServerConn {
	tcptask := NewServerConn(this.sctype, conn)
	if serverid == 0 {
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

// 广播消息
func (this *ServerConnPool) BroadcastCmd(v base.MsgStruct) {
	this.allsockets.Range(func(tkey interface{},
		tvalue interface{}) bool {
		tvalue.(*ServerConn).SendCmd(v)
		return true
	})
}

// 将一条消息广播至指定类型的所有连接
func (this *ServerConnPool) BroadcastByType(servertype uint32,
	v base.MsgStruct) {
	this.allsockets.Range(func(tkey interface{},
		tvalue interface{}) bool {
		value := tvalue.(*ServerConn)
		if value.Serverinfo.Servertype == servertype {
			value.SendCmd(v)
		}
		return true
	})
}

// 获取指定类型负载最小的一个连接
func (this *ServerConnPool) GetMinClient(
	servertype uint32) *ServerConn {
	var jobnum uint32 = 0xFFFFFFFF
	var serverid uint64 = 0
	this.allsockets.Range(func(tkey interface{},
		tvalue interface{}) bool {
		value := tvalue.(*ServerConn)
		key := tkey.(uint64)
		if value.Serverinfo.Servertype == servertype {
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
func (this *ServerConnPool) GetLatestVersionByType(servertype uint32) uint64 {
	latestVersion := uint64(0)
	this.allsockets.Range(func(tkey interface{},
		tvalue interface{}) bool {
		value := tvalue.(*ServerConn)
		if value.Serverinfo.Servertype == servertype &&
			value.Serverinfo.Version > latestVersion {
			latestVersion = value.Serverinfo.Version
		}
		return true
	})
	return latestVersion
}

// 获取指定类型负载最小的一个连接
func (this *ServerConnPool) GetMinClientLatestVersion(
	servertype uint32) *ServerConn {
	var jobnum uint32 = 0xFFFFFFFF
	var serverid uint64 = 0

	latestVersion := this.GetLatestVersionByType(servertype)

	this.allsockets.Range(func(tkey interface{},
		tvalue interface{}) bool {
		value := tvalue.(*ServerConn)
		key := tkey.(uint64)
		if value.Serverinfo.Servertype == servertype &&
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
	servertype uint32) *ServerConn {
	tasklist := make([]uint64, 0)
	this.allsockets.Range(func(tkey interface{},
		tvalue interface{}) bool {
		value := tvalue.(*ServerConn)
		key := tkey.(uint64)
		if value.Serverinfo.Servertype == servertype {
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
func (this *ServerConnPool) Get(tempid uint64) *ServerConn {
	if tcptask, found := this.allsockets.Load(tempid); found {
		return tcptask.(*ServerConn)
	}
	return nil
}

func (this *ServerConnPool) Remove(tempid uint64) {
	if tvalue, found := this.allsockets.Load(tempid); found {
		value := tvalue.(*ServerConn)
		value.isAlive = false
		// 关闭消息发送协程
		value.Conn.Shutdown()
		log.Debug("[ServerConn] 删除连接 TmpID[%d] 当前连接数量"+
			" Len[%d] ServerID[%d] ServerName[%s]",
			tempid, this.Len(), value.Serverinfo.Serverid,
			value.Serverinfo.Servername)
		// 删除连接
		this.remove(tempid)
		return
	}
}

func (this *ServerConnPool) Add(connct *ServerConn, tmpid uint64) {
	connct.Tempid = tmpid
	this.add(tmpid, connct)
}

func (this *ServerConnPool) AddAuto(connct *ServerConn) {
	tmpid, err := util.NewUniqueID(this.groupID)
	if err != nil {
		log.Error("[ServerConnPool.AddAuto] 生成UUID出错 Error[%s]",
			err.Error())
		return
	}
	connct.Tempid = tmpid
	this.add(connct.Tempid, connct)
}

// 修改链接的 tempip
func (this *ServerConnPool) ChangeTempid(tcptask *ServerConn,
	newtempid uint64) {
	this.TCPConnectPoolMutex.Lock()
	defer this.TCPConnectPoolMutex.Unlock()
	if ttcptask, found := this.allsockets.Load(tcptask.Tempid); found {
		tcptask := ttcptask.(*ServerConn)
		log.Debug("[ServerConn]修改连接tempid Old[%d] -->> New[%d]",
			tcptask.Tempid, newtempid)
		this.remove(tcptask.Tempid)
		tcptask.Tempid = newtempid
		this.add(tcptask.Tempid, tcptask)
		return
	}
}

func (this *ServerConnPool) Len() uint32 {
	if this.linkSum < 0 {
		return 0
	}
	return uint32(this.linkSum)
}

func (this *ServerConnPool) remove(tmpid uint64) {
	if _, ok := this.allsockets.Load(tmpid); !ok {
		return
	}
	// 删除连接
	this.allsockets.Delete(tmpid)
	this.linkSum--
}

func (this *ServerConnPool) add(tmpid uint64, value *ServerConn) {
	this.linkSumMutex.Lock()
	defer this.linkSumMutex.Unlock()
	if _, ok := this.allsockets.Load(tmpid); !ok {
		// 是新增的连接
		this.linkSum++
	}
	this.allsockets.Store(tmpid, value)
}
