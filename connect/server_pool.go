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

// 服务器连接池
type ServerPool struct {
	*log.Logger

	allSockets sync.Map // 所有连接
	linkSum    int
	groupID    uint16
}

// 初始化服务器连接池
func (this *ServerPool) Init(groupID uint16) {
	this.groupID = groupID
}

// 使用 TCP 创建一个服务器连接
func (this *ServerPool) NewTCPServer(sctype TServerSCType,
	conn net.Conn, moduleid string,
	onRecv func(*Server, *msg.MessageBinary),
	onClose func(*Server)) *Server {
	tcptask := &Server{}
	tcptask.SetLogger(this.Logger)
	tcptask.InitTCP(sctype, conn, onRecv, onClose)
	tcptask.Logger = this.Logger
	if moduleid == "" {
		this.AddServerAuto(tcptask)
	} else {
		this.AddServer(tcptask, moduleid)
	}
	return tcptask
}

// 使用 chan 创建一个服务器连接
func (this *ServerPool) NewChanServer(sctype TServerSCType,
	sendChan chan *msg.MessageBinary, recvChan chan *msg.MessageBinary,
	moduleid string, onRecv func(*Server, *msg.MessageBinary),
	onClose func(*Server)) *Server {
	tcptask := &Server{}
	tcptask.SetLogger(this.Logger)
	tcptask.InitChan(sctype, sendChan, recvChan, onRecv, onClose)
	tcptask.Logger = this.Logger
	if moduleid == "" {
		this.AddServerAuto(tcptask)
	} else {
		this.AddServer(tcptask, moduleid)
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
		if util.GetModuleIDType(value.ModuleInfo.ModuleID) == servertype {
			value.SendCmd(v)
		}
		return true
	})
}

// 广播消息到本连接池中的所有连接
func (this *ServerPool) BroadcastCmd(v msg.MsgStruct) {
	this.allSockets.Range(func(tkey interface{},
		tvalue interface{}) bool {
		tvalue.(*Server).SendCmd(v)
		return true
	})
}

// 获取指定类型负载最小的一个连接
func (this *ServerPool) GetMinLoadServer(servertype string) *Server {
	var jobnum uint32 = 0xFFFFFFFF
	var res *Server
	this.allSockets.Range(func(tkey interface{},
		tvalue interface{}) bool {
		value := tvalue.(*Server)
		if util.GetModuleIDType(value.ModuleInfo.ModuleID) == servertype {
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
		if util.GetModuleIDType(value.ModuleInfo.ModuleID) == servertype &&
			value.ModuleInfo.Version > latestVersion {
			latestVersion = value.ModuleInfo.Version
		}
		return true
	})
	return latestVersion
}

// 获取指定类型服务器的最新版本负载最小的一个连接
func (this *ServerPool) GetMinLoadServerLatestVersion(
	servertype string) *Server {
	var jobnum uint32 = 0xFFFFFFFF
	var moduleid uint64 = 0

	latestVersion := this.GetLatestVersionByType(servertype)

	this.allSockets.Range(func(tkey interface{},
		tvalue interface{}) bool {
		value := tvalue.(*Server)
		key := tkey.(uint64)
		if util.GetModuleIDType(value.ModuleInfo.ModuleID) == servertype &&
			value.ModuleInfo.Version == latestVersion {
			if jobnum >= value.GetJobNum() {
				jobnum = value.GetJobNum()
				moduleid = key
			}
		}
		return true
	})
	if moduleid == 0 {
		return nil
	}
	if tcptask, found := this.allSockets.Load(moduleid); found {
		return tcptask.(*Server)
	}
	return nil
}

// 随机获取指定类型的一个连接
func (this *ServerPool) GetRandomServer(servertype string) *Server {
	tasklist := make([]string, 0)
	this.allSockets.Range(func(tkey interface{},
		tvalue interface{}) bool {
		value := tvalue.(*Server)
		key := tkey.(string)
		if util.GetModuleIDType(value.ModuleInfo.ModuleID) == servertype {
			tasklist = append(tasklist, key)
		}
		return true
	})

	length := len(tasklist)
	if length > 0 {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		tmpindex := r.Intn(length)
		moduleid := tasklist[tmpindex]

		if tcptask, found := this.allSockets.Load(moduleid); found {
			return tcptask.(*Server)
		}
	}
	return nil
}

// 根据连接的 TmpID 获取一个连接
func (this *ServerPool) GetServer(tmpid string) *Server {
	if tcptask, found := this.allSockets.Load(tmpid); found {
		return tcptask.(*Server)
	}
	return nil
}

// 通过连接的 TmpID 从该连接池移除一个服务器连接
func (this *ServerPool) RemoveServer(tmpid string) {
	if tvalue, found := this.allSockets.Load(tmpid); found {
		value := tvalue.(*Server)
		// 关闭连接
		value.Shutdown()
		// 删除连接
		this.remove(tmpid)
		this.Syslog("[ServerPool] 断开连接 TmpID[%s] 当前连接数量"+
			" LinkSum[%d] ModuleID[%s]",
			tmpid, this.Len(), value.ModuleInfo.ModuleID)
		return
	}
}

// 在本连接池中新增一个服务器连接，并且指定该连接的 TmpID
func (this *ServerPool) AddServer(connct *Server, tmpid string) {
	connct.setTempID(tmpid)
	this.add(tmpid, connct)
	this.Syslog("[ServerPool] 增加连接 TmpID[%s] 当前连接数量"+
		" LinkSum[%d] ModuleID[%s]",
		connct.GetTempID(), this.Len(), connct.ModuleInfo.ModuleID)
}

// 在本连接池中新增一个服务器连接
func (this *ServerPool) AddServerAuto(connct *Server) {
	this.add(connct.GetTempID(), connct)
	this.Syslog("[ServerPool] 增加连接 TmpID[%s] 当前连接数量"+
		" LinkSum[%d] ModuleID[%s]",
		connct.GetTempID(), this.Len(), connct.ModuleInfo.ModuleID)
}

// 修改链接的 TmpID ，目标 TmpID 不可已存在于该连接池，否则返回 error
func (this *ServerPool) ChangeServerTempid(tcptask *Server,
	newTmpID string) error {
	afterI, isLoad := this.allSockets.LoadOrStore(newTmpID, tcptask)
	if isLoad {
		return fmt.Errorf("目标连接已存在:%s", newTmpID)
	} else {
		after := afterI.(*Server)
		oldTmpID := after.GetTempID()
		// 修改连接内的唯一ID标识
		after.setTempID(newTmpID)
		// 删除旧ID的索引，注意，如果你的ID生成规则不是唯一的，这里会有并发问题
		this.remove(oldTmpID)
		this.linkSum++
		this.Syslog("[ServerPool]修改连接tmpid Old[%s] -->> New[%s]",
			oldTmpID, newTmpID)
	}
	return nil
}

// 当前连接池的连接数量
func (this *ServerPool) Len() int {
	if this.linkSum < 0 {
		return 0
	}
	return this.linkSum
}

// 移除一个指定 TmpID 的连接
func (this *ServerPool) remove(tmpid string) {
	if _, ok := this.allSockets.Load(tmpid); !ok {
		return
	}
	// 删除连接
	this.allSockets.Delete(tmpid)
	this.linkSum--
}

// 增加一个服务器连接
func (this *ServerPool) add(tmpid string, value *Server) {
	_, isLoad := this.allSockets.LoadOrStore(tmpid, value)
	if !isLoad {
		this.linkSum++
	} else {
		this.allSockets.Store(tmpid, value)
	}
}
