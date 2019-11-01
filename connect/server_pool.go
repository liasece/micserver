package connect

import (
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/util"
)

type ServerPool struct {
	*log.Logger

	allSockets sync.Map // 所有连接
	linkSum    int32
	groupID    uint16
}

func (this *ServerPool) Init(groupID uint16) {
	this.groupID = groupID
}

func (this *ServerPool) NewTCPServer(sctype TServerSCType,
	conn net.Conn, serverid string,
	onRecv func(*Server, *msg.MessageBinary),
	onClose func(*Server)) *Server {
	tcptask := &Server{}
	tcptask.SetLogger(this.Logger)
	tcptask.InitTCP(sctype, conn, onRecv, onClose)
	tcptask.Logger = this.Logger
	if serverid == "" {
		this.AddServerAuto(tcptask)
	} else {
		this.AddServer(tcptask, serverid)
	}
	return tcptask
}

func (this *ServerPool) NewChanServer(sctype TServerSCType,
	sendChan chan *msg.MessageBinary, recvChan chan *msg.MessageBinary,
	serverid string, onRecv func(*Server, *msg.MessageBinary),
	onClose func(*Server)) *Server {
	tcptask := &Server{}
	tcptask.SetLogger(this.Logger)
	tcptask.InitChan(sctype, sendChan, recvChan, onRecv, onClose)
	tcptask.Logger = this.Logger
	if serverid == "" {
		this.AddServerAuto(tcptask)
	} else {
		this.AddServer(tcptask, serverid)
	}
	return tcptask
}

// 遍历连接池中的所有连接
func (this *ServerPool) RangeServer(
	callback func(*Server) bool) {
	this.allSockets.Range(func(tkey interface{},
		tvalue interface{}) bool {
		return callback(tvalue.(*Server))
	})
}

// 将一条消息广播至指定类型的所有连接
func (this *ServerPool) BroadcastByType(servertype string,
	v msg.MsgStruct) {
	this.allSockets.Range(func(tkey interface{},
		tvalue interface{}) bool {
		value := tvalue.(*Server)
		if util.GetServerIDType(value.ServerInfo.ServerID) == servertype {
			value.SendCmd(v)
		}
		return true
	})
}

// 广播消息
func (this *ServerPool) BroadcastCmd(v msg.MsgStruct) {
	this.allSockets.Range(func(tkey interface{},
		tvalue interface{}) bool {
		tvalue.(*Server).SendCmd(v)
		return true
	})
}

// 获取指定类型负载最小的一个连接
func (this *ServerPool) GetMinServer(servertype string) *Server {
	var jobnum uint32 = 0xFFFFFFFF
	var res *Server
	this.allSockets.Range(func(tkey interface{},
		tvalue interface{}) bool {
		value := tvalue.(*Server)
		if util.GetServerIDType(value.ServerInfo.ServerID) == servertype {
			if jobnum >= value.GetJobNum() {
				jobnum = value.GetJobNum()
				res = value
			}
		}
		return true
	})
	return res
}

// 获取指定类型服务器的最新版本
func (this *ServerPool) GetLatestVersionByType(servertype string) uint64 {
	latestVersion := uint64(0)
	this.allSockets.Range(func(tkey interface{},
		tvalue interface{}) bool {
		value := tvalue.(*Server)
		if util.GetServerIDType(value.ServerInfo.ServerID) == servertype &&
			value.ServerInfo.Version > latestVersion {
			latestVersion = value.ServerInfo.Version
		}
		return true
	})
	return latestVersion
}

// 获取指定类型负载最小的一个连接
func (this *ServerPool) GetMinServerLatestVersion(
	servertype string) *Server {
	var jobnum uint32 = 0xFFFFFFFF
	var serverid uint64 = 0

	latestVersion := this.GetLatestVersionByType(servertype)

	this.allSockets.Range(func(tkey interface{},
		tvalue interface{}) bool {
		value := tvalue.(*Server)
		key := tkey.(uint64)
		if util.GetServerIDType(value.ServerInfo.ServerID) == servertype &&
			value.ServerInfo.Version == latestVersion {
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
	if tcptask, found := this.allSockets.Load(serverid); found {
		return tcptask.(*Server)
	}
	return nil
}

// 随机获取指定类型的一个连接
func (this *ServerPool) GetRandomServer(
	servertype string) *Server {
	tasklist := make([]string, 0)
	this.allSockets.Range(func(tkey interface{},
		tvalue interface{}) bool {
		value := tvalue.(*Server)
		key := tkey.(string)
		if util.GetServerIDType(value.ServerInfo.ServerID) == servertype {
			tasklist = append(tasklist, key)
		}
		return true
	})

	length := len(tasklist)
	if length > 0 {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		tmpindex := r.Intn(length)
		serverid := tasklist[tmpindex]

		if tcptask, found := this.allSockets.Load(serverid); found {
			return tcptask.(*Server)
		}
	}
	return nil
}

// 根据连接的 TmpID 获取一个连接
func (this *ServerPool) GetServer(tempid string) *Server {
	if tcptask, found := this.allSockets.Load(tempid); found {
		return tcptask.(*Server)
	}
	return nil
}

func (this *ServerPool) RemoveServer(tempid string) {
	if tvalue, found := this.allSockets.Load(tempid); found {
		value := tvalue.(*Server)
		// 关闭连接
		value.Shutdown()
		// 删除连接
		this.remove(tempid)
		this.Debug("[ServerPool] 删除连接 TmpID[%s] 当前连接数量"+
			" LinkSum[%d] ServerID[%s]",
			tempid, this.ServerSum(), value.ServerInfo.ServerID)
		return
	}
}

func (this *ServerPool) AddServer(connct *Server, tmpid string) {
	connct.SetTempID(tmpid)
	this.add(tmpid, connct)
	this.Debug("[ServerPool] 增加连接 TmpID[%s] 当前连接数量"+
		" LinkSum[%d] ServerID[%s]",
		connct.GetTempID(), this.ServerSum(), connct.ServerInfo.ServerID)
}

func (this *ServerPool) AddServerAuto(connct *Server) {
	this.add(connct.GetTempID(), connct)
	this.Debug("[ServerPool] 增加连接 TmpID[%s] 当前连接数量"+
		" LinkSum[%d] ServerID[%s]",
		connct.GetTempID(), this.ServerSum(), connct.ServerInfo.ServerID)
}

// 修改链接的 tempip
func (this *ServerPool) ChangeServerTempid(tcptask *Server,
	newTempID string) error {
	afterI, isLoad := this.allSockets.LoadOrStore(newTempID, tcptask)
	if isLoad {
		return fmt.Errorf("目标连接已存在:%s", newTempID)
	} else {
		after := afterI.(*Server)
		oldTmpID := after.GetTempID()
		// 修改连接内的唯一ID标识
		after.SetTempID(newTempID)
		// 删除旧ID的索引，注意，如果你的ID生成规则不是唯一的，这里会有并发问题
		this.remove(oldTmpID)
		this.linkSum++
		this.Debug("[ServerPool]修改连接tempid Old[%s] -->> New[%s]",
			oldTmpID, newTempID)
	}
	return nil
}

func (this *ServerPool) ServerSum() uint32 {
	if this.linkSum < 0 {
		return 0
	}
	return uint32(this.linkSum)
}

func (this *ServerPool) remove(tmpid string) {
	if _, ok := this.allSockets.Load(tmpid); !ok {
		return
	}
	// 删除连接
	this.allSockets.Delete(tmpid)
	this.linkSum--
}

func (this *ServerPool) add(tmpid string, value *Server) {
	_, isLoad := this.allSockets.LoadOrStore(tmpid, value)
	if !isLoad {
		this.linkSum++
	} else {
		this.allSockets.Store(tmpid, value)
	}
}
