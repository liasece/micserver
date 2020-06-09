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

// ServerPool 服务器连接池
type ServerPool struct {
	*log.Logger

	allSockets sync.Map // 所有连接
	linkSum    int
	groupID    uint16
}

// Init 初始化服务器连接池
func (sp *ServerPool) Init(groupID uint16) {
	sp.groupID = groupID
}

// NewTCPServer 使用 TCP 创建一个服务器连接
func (sp *ServerPool) NewTCPServer(sctype TServerSCType,
	conn net.Conn, moduleid string,
	onRecv func(*Server, *msg.MessageBinary),
	onClose func(*Server)) *Server {
	tcptask := &Server{}
	tcptask.SetLogger(sp.Logger)
	tcptask.InitTCP(sctype, conn, onRecv, onClose)
	tcptask.Logger = sp.Logger
	if moduleid == "" {
		sp.AddServerAuto(tcptask)
	} else {
		sp.AddServer(tcptask, moduleid)
	}
	return tcptask
}

// NewChanServer 使用 chan 创建一个服务器连接
func (sp *ServerPool) NewChanServer(sctype TServerSCType,
	sendChan chan *msg.MessageBinary, recvChan chan *msg.MessageBinary,
	moduleid string, onRecv func(*Server, *msg.MessageBinary),
	onClose func(*Server)) *Server {
	tcptask := &Server{}
	tcptask.SetLogger(sp.Logger)
	tcptask.InitChan(sctype, sendChan, recvChan, onRecv, onClose)
	tcptask.Logger = sp.Logger
	if moduleid == "" {
		sp.AddServerAuto(tcptask)
	} else {
		sp.AddServer(tcptask, moduleid)
	}
	return tcptask
}

// RangeServer 遍历连接池中的所有连接
func (sp *ServerPool) RangeServer(
	callback func(*Server) bool) {
	sp.allSockets.Range(func(tkey interface{},
		tvalue interface{}) bool {
		return callback(tvalue.(*Server))
	})
}

// BroadcastByType 将一条消息广播至指定类型的所有连接
func (sp *ServerPool) BroadcastByType(servertype string,
	v msg.IMsgStruct) {
	sp.allSockets.Range(func(tkey interface{},
		tvalue interface{}) bool {
		value := tvalue.(*Server)
		if util.GetModuleIDType(value.ModuleInfo.ModuleID) == servertype {
			value.SendCmd(v)
		}
		return true
	})
}

// BroadcastCmd 广播消息到本连接池中的所有连接
func (sp *ServerPool) BroadcastCmd(v msg.IMsgStruct) {
	sp.allSockets.Range(func(tkey interface{},
		tvalue interface{}) bool {
		tvalue.(*Server).SendCmd(v)
		return true
	})
}

// GetMinLoadServer 获取指定类型负载最小的一个连接
func (sp *ServerPool) GetMinLoadServer(servertype string) *Server {
	var jobnum uint32 = 0xFFFFFFFF
	var res *Server
	sp.allSockets.Range(func(tkey interface{},
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

// GetLatestVersionByType 获取指定类型服务器的最新版本
func (sp *ServerPool) GetLatestVersionByType(servertype string) uint64 {
	latestVersion := uint64(0)
	sp.allSockets.Range(func(tkey interface{},
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

// GetMinLoadServerLatestVersion 获取指定类型服务器的最新版本负载最小的一个连接
func (sp *ServerPool) GetMinLoadServerLatestVersion(
	servertype string) *Server {
	var jobnum uint32 = 0xFFFFFFFF
	var moduleid uint64 = 0

	latestVersion := sp.GetLatestVersionByType(servertype)

	sp.allSockets.Range(func(tkey interface{},
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
	if tcptask, found := sp.allSockets.Load(moduleid); found {
		return tcptask.(*Server)
	}
	return nil
}

// GetRandomServer 随机获取指定类型的一个连接
func (sp *ServerPool) GetRandomServer(servertype string) *Server {
	tasklist := make([]string, 0)
	sp.allSockets.Range(func(tkey interface{},
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

		if tcptask, found := sp.allSockets.Load(moduleid); found {
			return tcptask.(*Server)
		}
	}
	return nil
}

// GetServer 根据连接的 TmpID 获取一个连接
func (sp *ServerPool) GetServer(tmpid string) *Server {
	if tcptask, found := sp.allSockets.Load(tmpid); found {
		return tcptask.(*Server)
	}
	return nil
}

// RemoveServer 通过连接的 TmpID 从该连接池移除一个服务器连接
func (sp *ServerPool) RemoveServer(tmpid string) {
	if tvalue, found := sp.allSockets.Load(tmpid); found {
		value := tvalue.(*Server)
		// 关闭连接
		value.Shutdown()
		// 删除连接
		sp.remove(tmpid)
		sp.Syslog("[ServerPool] 断开连接 TmpID[%s] 当前连接数量"+
			" LinkSum[%d] ModuleID[%s]",
			tmpid, sp.Len(), value.ModuleInfo.ModuleID)
		return
	}
}

// AddServer 在本连接池中新增一个服务器连接，并且指定该连接的 TmpID
func (sp *ServerPool) AddServer(connct *Server, tmpid string) {
	connct.setTempID(tmpid)
	sp.add(tmpid, connct)
	sp.Syslog("[ServerPool] 增加连接 TmpID[%s] 当前连接数量"+
		" LinkSum[%d] ModuleID[%s]",
		connct.GetTempID(), sp.Len(), connct.ModuleInfo.ModuleID)
}

// AddServerAuto 在本连接池中新增一个服务器连接
func (sp *ServerPool) AddServerAuto(connct *Server) {
	sp.add(connct.GetTempID(), connct)
	sp.Syslog("[ServerPool] 增加连接 TmpID[%s] 当前连接数量"+
		" LinkSum[%d] ModuleID[%s]",
		connct.GetTempID(), sp.Len(), connct.ModuleInfo.ModuleID)
}

// ChangeServerTempid 修改链接的 TmpID ，目标 TmpID 不可已存在于该连接池，否则返回 error
func (sp *ServerPool) ChangeServerTempid(tcptask *Server,
	newTmpID string) error {
	afterI, isLoad := sp.allSockets.LoadOrStore(newTmpID, tcptask)
	if isLoad {
		return fmt.Errorf("目标连接已存在:%s", newTmpID)
	}
	after := afterI.(*Server)
	oldTmpID := after.GetTempID()
	// 修改连接内的唯一ID标识
	after.setTempID(newTmpID)
	// 删除旧ID的索引，注意，如果你的ID生成规则不是唯一的，这里会有并发问题
	sp.remove(oldTmpID)
	sp.linkSum++
	sp.Syslog("[ServerPool]修改连接tmpid Old[%s] -->> New[%s]",
		oldTmpID, newTmpID)
	return nil
}

// Len 当前连接池的连接数量
func (sp *ServerPool) Len() int {
	if sp.linkSum < 0 {
		return 0
	}
	return sp.linkSum
}

// remove 移除一个指定 TmpID 的连接
func (sp *ServerPool) remove(tmpid string) {
	if _, ok := sp.allSockets.Load(tmpid); !ok {
		return
	}
	// 删除连接
	sp.allSockets.Delete(tmpid)
	sp.linkSum--
}

// add 增加一个服务器连接
func (sp *ServerPool) add(tmpid string, value *Server) {
	_, isLoad := sp.allSockets.LoadOrStore(tmpid, value)
	if !isLoad {
		sp.linkSum++
	} else {
		sp.allSockets.Store(tmpid, value)
	}
}
