package connect

import (
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

	allSockets sync.Map // 所有连接
	linkSum    int32
	groupID    uint16
}

func (this *ServerConnPool) Init(groupID uint16) {
	this.groupID = groupID
}

func (this *ServerConnPool) NewServerConn(sctype TServerSCType,
	conn net.Conn, serverid string,
	onRecv func(*ServerConn, *msg.MessageBinary),
	onClose func(*ServerConn)) *ServerConn {
	tcptask := NewServerConn(sctype, conn, onRecv, onClose)
	tcptask.Logger = this.Logger
	if serverid == "" {
		this.AddServerConnAuto(tcptask)
	} else {
		this.AddServerConn(tcptask, serverid)
	}
	return tcptask
}

// 遍历连接池中的所有连接
func (this *ServerConnPool) RangeServerConn(
	callback func(*ServerConn) bool) {
	this.allSockets.Range(func(tkey interface{},
		tvalue interface{}) bool {
		return callback(tvalue.(*ServerConn))
	})
}

// 将一条消息广播至指定类型的所有连接
func (this *ServerConnPool) BroadcastByType(servertype string,
	v msg.MsgStruct) {
	this.allSockets.Range(func(tkey interface{},
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
	this.allSockets.Range(func(tkey interface{},
		tvalue interface{}) bool {
		tvalue.(*ServerConn).SendCmd(v)
		return true
	})
}

// 获取指定类型负载最小的一个连接
func (this *ServerConnPool) GetMinServerConn(servertype string) *ServerConn {
	var jobnum uint32 = 0xFFFFFFFF
	var res *ServerConn
	this.allSockets.Range(func(tkey interface{},
		tvalue interface{}) bool {
		value := tvalue.(*ServerConn)
		if util.GetServerIDType(value.Serverinfo.ServerID) == servertype {
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
func (this *ServerConnPool) GetLatestVersionByType(servertype string) uint64 {
	latestVersion := uint64(0)
	this.allSockets.Range(func(tkey interface{},
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
func (this *ServerConnPool) GetMinServerConnLatestVersion(
	servertype string) *ServerConn {
	var jobnum uint32 = 0xFFFFFFFF
	var serverid uint64 = 0

	latestVersion := this.GetLatestVersionByType(servertype)

	this.allSockets.Range(func(tkey interface{},
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
	if tcptask, found := this.allSockets.Load(serverid); found {
		return tcptask.(*ServerConn)
	}
	return nil
}

// 随机获取指定类型的一个连接
func (this *ServerConnPool) GetRandomServerConn(
	servertype string) *ServerConn {
	tasklist := make([]string, 0)
	this.allSockets.Range(func(tkey interface{},
		tvalue interface{}) bool {
		value := tvalue.(*ServerConn)
		key := tkey.(string)
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

		if tcptask, found := this.allSockets.Load(serverid); found {
			return tcptask.(*ServerConn)
		}
	}
	return nil
}

// 根据连接的 TmpID 获取一个连接
func (this *ServerConnPool) GetServerConn(tempid string) *ServerConn {
	if tcptask, found := this.allSockets.Load(tempid); found {
		return tcptask.(*ServerConn)
	}
	return nil
}

func (this *ServerConnPool) RemoveServerConn(tempid string) {
	if tvalue, found := this.allSockets.Load(tempid); found {
		value := tvalue.(*ServerConn)
		// 关闭连接
		value.Shutdown()
		// 删除连接
		this.remove(tempid)
		this.Debug("[ServerConnPool] 删除连接 TmpID[%s] 当前连接数量"+
			" LinkSum[%d] ServerID[%s]",
			tempid, this.ServerConnSum(), value.Serverinfo.ServerID)
		return
	}
}

func (this *ServerConnPool) AddServerConn(connct *ServerConn, tmpid string) {
	connct.Tempid = tmpid
	this.add(tmpid, connct)
	this.Debug("[ServerConnPool] 增加连接 TmpID[%s] 当前连接数量"+
		" LinkSum[%d] ServerID[%s]",
		connct.Tempid, this.ServerConnSum(), connct.Serverinfo.ServerID)
}

func (this *ServerConnPool) AddServerConnAuto(connct *ServerConn) {
	tmpid, err := util.NewUniqueID(this.groupID)
	if err != nil {
		this.Error("[ServerConnPool.AddAuto] 生成UUID出错 Error[%s]",
			err.Error())
		return
	}
	connct.Tempid = tmpid
	this.add(connct.Tempid, connct)
	this.Debug("[ServerConnPool] 增加连接 TmpID[%s] 当前连接数量"+
		" LinkSum[%d] ServerID[%s]",
		connct.Tempid, this.ServerConnSum(), connct.Serverinfo.ServerID)
}

// 修改链接的 tempip
func (this *ServerConnPool) ChangeServerConnTempid(tcptask *ServerConn,
	newTempID string) error {
	afterI, isLoad := this.allSockets.LoadOrStore(newTempID, tcptask)
	if isLoad {
		return fmt.Errorf("目标连接已存在:%s", newTempID)
	} else {
		after := afterI.(*ServerConn)
		oldTmpID := after.Tempid
		// 修改连接内的唯一ID标识
		after.Tempid = newTempID
		// 删除旧ID的索引，注意，如果你的ID生成规则不是唯一的，这里会有并发问题
		this.remove(oldTmpID)
		this.linkSum++
		this.Debug("[ServerConnPool]修改连接tempid Old[%s] -->> New[%s]",
			oldTmpID, newTempID)
	}
	return nil
}

func (this *ServerConnPool) ServerConnSum() uint32 {
	if this.linkSum < 0 {
		return 0
	}
	return uint32(this.linkSum)
}

func (this *ServerConnPool) remove(tmpid string) {
	if _, ok := this.allSockets.Load(tmpid); !ok {
		return
	}
	// 删除连接
	this.allSockets.Delete(tmpid)
	this.linkSum--
}

func (this *ServerConnPool) add(tmpid string, value *ServerConn) {
	_, isLoad := this.allSockets.LoadOrStore(tmpid, value)
	if !isLoad {
		this.linkSum++
	} else {
		this.allSockets.Store(tmpid, value)
	}
}
