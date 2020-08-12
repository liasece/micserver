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
func (sp *ServerPool) NewTCPServer(scType TServerSCType, conn net.Conn, moduleID string, onRecv func(*Server, *msg.MessageBinary), onClose func(*Server)) *Server {
	tcptask := &Server{}
	tcptask.SetLogger(sp.Logger)
	tcptask.InitTCP(scType, conn, onRecv, onClose)
	tcptask.Logger = sp.Logger
	if moduleID == "" {
		sp.AddServerAuto(tcptask)
	} else {
		sp.AddServer(tcptask, moduleID)
	}
	return tcptask
}

// NewChanServer 使用 chan 创建一个服务器连接
func (sp *ServerPool) NewChanServer(scType TServerSCType, sendChan chan *msg.MessageBinary, recvChan chan *msg.MessageBinary, moduleID string, onRecv func(*Server, *msg.MessageBinary), onClose func(*Server)) *Server {
	tcptask := &Server{}
	tcptask.SetLogger(sp.Logger)
	tcptask.InitChan(scType, sendChan, recvChan, onRecv, onClose)
	tcptask.Logger = sp.Logger
	if moduleID == "" {
		sp.AddServerAuto(tcptask)
	} else {
		sp.AddServer(tcptask, moduleID)
	}
	return tcptask
}

// RangeServer 遍历连接池中的所有连接
func (sp *ServerPool) RangeServer(
	callback func(*Server) bool) {
	sp.allSockets.Range(func(tKey interface{},
		tValue interface{}) bool {
		return callback(tValue.(*Server))
	})
}

// BroadcastByType 将一条消息广播至指定类型的所有连接
func (sp *ServerPool) BroadcastByType(servertype string, v msg.IMsgStruct) {
	sp.allSockets.Range(func(tKey interface{},
		tValue interface{}) bool {
		value := tValue.(*Server)
		if util.GetModuleIDType(value.ModuleInfo.ModuleID) == servertype {
			value.SendCmd(v)
		}
		return true
	})
}

// BroadcastCmd 广播消息到本连接池中的所有连接
func (sp *ServerPool) BroadcastCmd(v msg.IMsgStruct) {
	sp.allSockets.Range(func(tKey interface{},
		tValue interface{}) bool {
		tValue.(*Server).SendCmd(v)
		return true
	})
}

// GetMinLoadServer 获取指定类型负载最小的一个连接
func (sp *ServerPool) GetMinLoadServer(servertype string) *Server {
	var jobNum uint32 = 0xFFFFFFFF
	var res *Server
	sp.allSockets.Range(func(tKey interface{},
		tValue interface{}) bool {
		value := tValue.(*Server)
		if util.GetModuleIDType(value.ModuleInfo.ModuleID) == servertype {
			if jobNum >= value.GetJobNum() {
				jobNum = value.GetJobNum()
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
	sp.allSockets.Range(func(tKey interface{},
		tValue interface{}) bool {
		value := tValue.(*Server)
		if util.GetModuleIDType(value.ModuleInfo.ModuleID) == servertype && value.ModuleInfo.Version > latestVersion {
			latestVersion = value.ModuleInfo.Version
		}
		return true
	})
	return latestVersion
}

// GetMinLoadServerLatestVersion 获取指定类型服务器的最新版本负载最小的一个连接
func (sp *ServerPool) GetMinLoadServerLatestVersion(servertype string) *Server {
	var jobNum uint32 = 0xFFFFFFFF
	var moduleID uint64 = 0

	latestVersion := sp.GetLatestVersionByType(servertype)

	sp.allSockets.Range(func(tKey interface{}, tValue interface{}) bool {
		value := tValue.(*Server)
		key := tKey.(uint64)
		if util.GetModuleIDType(value.ModuleInfo.ModuleID) == servertype && value.ModuleInfo.Version == latestVersion {
			if jobNum >= value.GetJobNum() {
				jobNum = value.GetJobNum()
				moduleID = key
			}
		}
		return true
	})
	if moduleID == 0 {
		return nil
	}
	if tcptask, found := sp.allSockets.Load(moduleID); found {
		return tcptask.(*Server)
	}
	return nil
}

// GetRandomServer 随机获取指定类型的一个连接
func (sp *ServerPool) GetRandomServer(servertype string) *Server {
	tasklist := make([]string, 0)
	sp.allSockets.Range(func(tKey interface{},
		tValue interface{}) bool {
		value := tValue.(*Server)
		key := tKey.(string)
		if util.GetModuleIDType(value.ModuleInfo.ModuleID) == servertype {
			tasklist = append(tasklist, key)
		}
		return true
	})

	length := len(tasklist)
	if length > 0 {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		tmpIndex := r.Intn(length)
		moduleID := tasklist[tmpIndex]

		if tcptask, found := sp.allSockets.Load(moduleID); found {
			return tcptask.(*Server)
		}
	}
	return nil
}

// GetServer 根据连接的 TempID 获取一个连接
func (sp *ServerPool) GetServer(tempID string) *Server {
	if tcptask, found := sp.allSockets.Load(tempID); found {
		return tcptask.(*Server)
	}
	return nil
}

// RemoveServer 通过连接的 TempID 从该连接池移除一个服务器连接
func (sp *ServerPool) RemoveServer(tempID string) {
	if tValue, found := sp.allSockets.Load(tempID); found {
		value := tValue.(*Server)
		// 关闭连接
		value.Shutdown()
		// 删除连接
		sp.remove(tempID)
		sp.Syslog("[ServerPool] RemoveServer connection", log.String("TempID", tempID), log.String("ModuleID", value.ModuleInfo.ModuleID), log.Int("CurrentNum", sp.Len()))
		return
	}
}

// AddServer 在本连接池中新增一个服务器连接，并且指定该连接的 TempID
func (sp *ServerPool) AddServer(svr *Server, tempID string) {
	svr.setTempID(tempID)
	sp.add(tempID, svr)
	sp.Syslog("[ServerPool] AddServer connection", log.String("TempID", svr.GetTempID()), log.String("ModuleID", svr.ModuleInfo.ModuleID), log.Int("CurrentNum", sp.Len()))
}

// AddServerAuto 在本连接池中新增一个服务器连接
func (sp *ServerPool) AddServerAuto(svr *Server) {
	sp.add(svr.GetTempID(), svr)
	sp.Syslog("[ServerPool] AddServerAuto connection", log.String("TempID", svr.GetTempID()), log.String("ModuleID", svr.ModuleInfo.ModuleID), log.Int("CurrentNum", sp.Len()))
}

// ChangeServerTempID 修改链接的 TempID ，目标 TempID 不可已存在于该连接池，否则返回 error
func (sp *ServerPool) ChangeServerTempID(tcptask *Server,
	newTmpID string) error {
	afterI, isLoad := sp.allSockets.LoadOrStore(newTmpID, tcptask)
	if isLoad {
		return fmt.Errorf("Target connection tempID already exists:%s", newTmpID)
	}
	after := afterI.(*Server)
	oldTmpID := after.GetTempID()
	// 修改连接内的唯一ID标识
	after.setTempID(newTmpID)
	// 删除旧ID的索引，注意，如果你的ID生成规则不是唯一的，这里会有并发问题
	sp.remove(oldTmpID)
	sp.linkSum++
	sp.Syslog("[ServerPool] Change target server tempID", log.String("Old", oldTmpID), log.String("New", newTmpID))
	return nil
}

// Len 当前连接池的连接数量
func (sp *ServerPool) Len() int {
	if sp.linkSum < 0 {
		return 0
	}
	return sp.linkSum
}

// remove 移除一个指定 TempID 的连接
func (sp *ServerPool) remove(tempID string) {
	if _, ok := sp.allSockets.Load(tempID); !ok {
		return
	}
	// 删除连接
	sp.allSockets.Delete(tempID)
	sp.linkSum--
}

// add 增加一个服务器连接
func (sp *ServerPool) add(tempID string, value *Server) {
	_, isLoad := sp.allSockets.LoadOrStore(tempID, value)
	if !isLoad {
		sp.linkSum++
	} else {
		sp.allSockets.Store(tempID, value)
	}
}
