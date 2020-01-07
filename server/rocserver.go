/*
micserver中的ROC调用发生时，处理调用以及返回值。
*/
package server

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/liasece/micserver/conf"
	"github.com/liasece/micserver/connect"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/msg"
	"github.com/liasece/micserver/process"
	"github.com/liasece/micserver/roc"
	"github.com/liasece/micserver/servercomm"
	"github.com/liasece/micserver/util/hash"
)

// ROC请求信息
type requestAgent struct {
	fromModuleID string
	callpath     string
	callarg      []byte
	seq          int64
	needReturn   bool
}

// ROC响应信息
type responseAgent struct {
	fromModuleID string
	seq          int64
	data         []byte
	err          string
}

// ROC服务
type ROCServer struct {
	*log.Logger
	server *Server

	// 记录在本地的缓存信息
	// 第一层键为ROCObj类型，第二层键为ROCObj的ID
	localObj      map[string]map[string]struct{}
	localObjMutex sync.Mutex

	// 远程对象调用支持
	_ROCManager     roc.ROCManager
	rocAddCacheChan chan roc.IObj
	rocDelCacheChan chan roc.IObj

	rocRequestChan  chan *requestAgent
	rocResponseChan chan *responseAgent
	rocBlockChanMap sync.Map

	seqMutex sync.Mutex
	lastSeq  int64
}

// 初始化ROC服务
func (this *ROCServer) Init(server *Server) {
	this.server = server
	this.Logger = server.Logger.Clone()
	this.Logger.SetTopic("ROCServer")

	this.rocAddCacheChan = make(chan roc.IObj, 10000)
	this.rocDelCacheChan = make(chan roc.IObj, 10000)
	go this.rocObjNoticeProcess(this.rocAddCacheChan, false)
	go this.rocObjNoticeProcess(this.rocDelCacheChan, true)
	this._ROCManager.HookObjEvent(this)

	this.rocRequestChan = make(chan *requestAgent, 10000)
	go this.rocRequestProcess()
	this.rocResponseChan = make(chan *responseAgent, 10000)
	go this.rocResponseProcess()
}

// 生成一个 ROC 调用的序号，在每个模块中应该唯一
func (this *ROCServer) newSeq() (res int64) {
	this.seqMutex.Lock()
	this.lastSeq++
	res = this.lastSeq
	this.seqMutex.Unlock()
	return
}

// 获取指定ROC对象类型的ROC对象
func (this *ROCServer) GetROC(objtype roc.ROCObjType) *roc.ROC {
	return this._ROCManager.GetROC(objtype)
}

// 新建一个指定ROC对象类型的ROC对象
func (this *ROCServer) NewROC(objtype roc.ROCObjType) *roc.ROC {
	return this._ROCManager.NewROC(objtype)
}

// 无返回值的ROC调用
func (this *ROCServer) ROCCallNR(callpath *roc.ROCPath, callarg []byte) error {
	objType := callpath.GetObjType()
	objID := callpath.GetObjID()
	moduleid := roc.GetCache().Get(objType, objID)
	this.Syslog("ROCCallNR {%s:%s(%s:%s):%X}",
		moduleid, callpath, objType, objID, callarg)
	// 构造消息
	sendmsg := &servercomm.SROCRequest{
		FromModuleID: this.server.moduleid,
		Seq:          this.newSeq(),
		CallStr:      callpath.String(),
		CallArg:      callarg,
	}
	if moduleid == this.server.moduleid {
		sendmsg.ToModuleID = moduleid
		this.onMsgROCRequest(sendmsg)
	} else {
		server := this.server.subnetManager.GetServer(moduleid)
		if server != nil {
			sendmsg.ToModuleID = moduleid
			server.SendCmd(sendmsg)
		} else {
			this.Warn("Can't find roc object location %s",
				callpath.String())
			return fmt.Errorf("Can't find roc object location %s",
				callpath.String())
		}
	}
	return nil
}

// 获取ROC缓存中的位置信息
// 返回目标ROC对象所在的moduleid
func (this *ROCServer) GetROCCachedLocation(objType roc.ROCObjType,
	objID string) string {
	moduleid := roc.GetCache().Get(objType, objID)
	return moduleid
}

// 遍历指定类型的ROC缓存，限制目标对象必须本module可以访问
func (this *ROCServer) RangeROCCachedByType(objType roc.ROCObjType,
	f func(id string, location string) bool) {
	connecedModuleIDs := make(map[string]bool)
	this.server.subnetManager.RangeServer(func(server *connect.Server) bool {
		if server.ModuleInfo != nil {
			connecedModuleIDs[server.ModuleInfo.ModuleID] = true
		}
		return true
	})
	roc.GetCache().RangeByType(objType, f, connecedModuleIDs)
}

// 随机获取本地缓存的ROC对象，返回该对象的ID，限制目标对象必须本module可以访问
func (this *ROCServer) RandomROCCachedByType(objType roc.ROCObjType) string {
	connecedModuleIDs := make(map[string]bool)
	this.server.subnetManager.RangeServer(func(server *connect.Server) bool {
		if server.ModuleInfo != nil {
			connecedModuleIDs[server.ModuleInfo.ModuleID] = true
		}
		return true
	})
	return roc.GetCache().RandomObjIDByType(objType, connecedModuleIDs)
}

// 根据ROC请求的序号，生成一个用于阻塞等待ROC返回的chan
func (this *ROCServer) addBlockChan(seq int64) chan *responseAgent {
	ch := make(chan *responseAgent, 1)
	this.rocBlockChanMap.Store(seq, ch)
	return ch
}

// 有返回值的RPC调用
func (this *ROCServer) ROCCallBlock(callpath *roc.ROCPath,
	callarg []byte) ([]byte, error) {
	objType := callpath.GetObjType()
	objID := callpath.GetObjID()
	moduleid := roc.GetCache().Get(objType, objID)
	this.Syslog("ROCCallBlock {%s:%s(%s:%s:%d):%X}",
		moduleid, callpath, objType, objID, hash.GetStringHash(string(objID)),
		callarg)
	// 构造消息
	sendmsg := &servercomm.SROCRequest{
		FromModuleID: this.server.moduleid,
		Seq:          this.newSeq(),
		CallStr:      callpath.String(),
		CallArg:      callarg,
		NeedReturn:   true,
	}

	ch := this.addBlockChan(sendmsg.Seq)

	if moduleid == this.server.moduleid {
		sendmsg.ToModuleID = moduleid
		this.onMsgROCRequest(sendmsg)
	} else {
		server := this.server.subnetManager.GetServer(moduleid)
		if server != nil {
			sendmsg.ToModuleID = moduleid
			server.SendCmd(sendmsg)
		} else {
			this.server.subnetManager.BroadcastCmd(sendmsg)
		}
	}

	// 等待返回值
	agent := <-ch
	return agent.data, errors.New(agent.err)
}

// 当收到ROC调用请求时
func (this *ROCServer) onMsgROCRequest(msg *servercomm.SROCRequest) {
	agent := &requestAgent{
		callpath:     msg.CallStr,
		callarg:      msg.CallArg,
		seq:          msg.Seq,
		needReturn:   msg.NeedReturn,
		fromModuleID: msg.FromModuleID,
	}
	this.rocRequestChan <- agent
}

// 当收到ROC调用返回时
func (this *ROCServer) onMsgROCResponse(msg *servercomm.SROCResponse) {
	agent := &responseAgent{
		fromModuleID: msg.FromModuleID,
		seq:          msg.ReqSeq,
		data:         msg.ResData,
		err:          msg.Error,
	}
	this.rocResponseChan <- agent
}

// 处理带返回值的ROC调用返回的线程
func (this *ROCServer) rocRequestProcess() {
	tm := time.NewTimer(time.Millisecond * 300)
	for !this.server.isStop {
		select {
		case <-this.server.stopChan:
			break
		case agent := <-this.rocRequestChan:
			// 处理ROC请求
			this.Syslog("ROC Request[%s]", agent.callpath)
			res, err := this._ROCManager.Call(agent.callpath, agent.callarg)
			if err != nil {
				if err != roc.ErrUnknowObj {
					this.Error("ROCManager.Call err:%s", err.Error())
				} else {
					this.Syslog("ROCManager.Call Path[%s] Err[%s]",
						agent.callpath, err.Error())
				}
			} else {
				// this.Debug("ROC调用成功 res:%+v", res)
			}
			if agent.needReturn {
				if agent.fromModuleID == this.server.moduleid {
					// 本地服务器调用结果

				} else {
					server := this.server.subnetManager.GetServer(
						agent.fromModuleID)
					if server != nil {
						// 返回执行结果
						sendmsg := &servercomm.SROCResponse{
							FromModuleID: this.server.moduleid,
							ToModuleID:   agent.fromModuleID,
							ReqSeq:       agent.seq,
							ResData:      res,
							Error:        err.Error(),
						}
						server.SendCmd(sendmsg)
					}
				}
			}
		case <-tm.C:
			tm.Reset(time.Millisecond * 300)
			break
		}
	}
}

// 处理ROC带返回值电泳调用的线程
func (this *ROCServer) rocResponseProcess() {
	tm := time.NewTimer(time.Millisecond * 300)
	for !this.server.isStop {
		select {
		case <-this.server.stopChan:
			break
		case agent := <-this.rocResponseChan:
			// 处理ROC相应
			this.Syslog("rocResponseProcess %+v", agent)
			chi, ok := this.rocBlockChanMap.Load(agent.seq)
			if ok {
				if ch, ok := chi.(chan *responseAgent); ok {
					// 写入返回值
					ch <- agent
				} else {
					this.Error("ROC返回 chi.(chan *responseAgent) 错误")
				}
			} else {
				this.Error("ROC返回 不存在目标ROC请求 %+v", agent)
			}
		case <-tm.C:
			tm.Reset(time.Millisecond * 300)
			break
		}
	}
}

// 当一个服务器加入了子网时
func (this *ROCServer) onServerJoinSubnet(server *connect.Server) {
	if process.HasModule(server.ModuleInfo.ModuleID) {
		return
	}
	this.localObjMutex.Lock()
	defer this.localObjMutex.Unlock()

	if this.localObj == nil {
		return
	}
	// 遍历所有类型的ROCObj
	for objtype, typemap := range this.localObj {
		if typemap == nil || len(typemap) == 0 {
			continue
		}
		leftnum := len(typemap)
		// 临时变量
		tmplist := make([]string, 0)
		tmpsize := 0
		// 遍历所有对象
		for objid, _ := range typemap {
			leftnum--
			tmplist = append(tmplist, objid)
			tmpsize += len(objid) + 4
			// 分包发送
			if leftnum == 0 ||
				(tmpsize > (msg.MessageMaxSize/2) && tmpsize > 32*1024) {
				sendmsg := &servercomm.SROCBind{
					HostModuleID: this.server.moduleid,
					IsDelete:     false,
					ObjType:      objtype,
					ObjIDs:       tmplist,
				}
				server.SendCmd(sendmsg)
				if leftnum > 0 {
					tmplist = make([]string, 0)
					tmpsize = 0
				}
			}
		}
	}
}

// 记录本地的ROC对象
func (this *ROCServer) recordLocalObj(objtype string, objid string,
	isDelete bool) {
	this.localObjMutex.Lock()
	defer this.localObjMutex.Unlock()

	// 初始化内存
	if this.localObj == nil {
		this.localObj = make(map[string]map[string]struct{})
	}
	if v, ok := this.localObj[objtype]; !ok || v == nil {
		this.localObj[objtype] = make(map[string]struct{})
	}

	// 记录
	if isDelete {
		// 删除
		if _, ok := this.localObj[objtype][objid]; ok {
			delete(this.localObj[objtype], objid)
		}
	} else {
		this.localObj[objtype][objid] = struct{}{}
	}
}

// 当ROC对象发生注册行为时
func (this *ROCServer) OnROCObjAdd(obj roc.IObj) {
	// 保存本地映射缓存
	roc.GetCache().Set(obj.GetROCObjType(), obj.GetROCObjID(),
		this.server.moduleid)
	this.Syslog("OnROCObjAdd roc cache set type[%s] "+
		"id[%s] host[%s]",
		obj.GetROCObjType(), obj.GetROCObjID(), this.server.moduleid)

	// 由于ROC绑定消息与ROC调用之间存在异步问题，除非经过设置，
	// 否则使用同步方式同步ROC对象绑定
	if this.server.moduleConfig.GetBool(conf.AsynchronousSyncRocbind) {
		this.rocAddCacheChan <- obj
	} else {
		sendmsg := &servercomm.SROCBind{
			HostModuleID: this.server.moduleid,
			IsDelete:     false,
			ObjType:      string(obj.GetROCObjType()),
			ObjIDs:       []string{obj.GetROCObjID()},
		}
		this.sendROCBindMsg(sendmsg)
	}
	this.recordLocalObj(string(obj.GetROCObjType()), obj.GetROCObjID(),
		false)
}

// 当ROC对象发生注册行为时
func (this *ROCServer) OnROCObjDel(obj roc.IObj) {
	// 保存本地映射缓存
	roc.GetCache().Del(obj.GetROCObjType(), obj.GetROCObjID(),
		this.server.moduleid)
	this.Syslog("OnROCObjDel roc cache del type[%s] "+
		"id[%s] host[%s]",
		obj.GetROCObjType(), obj.GetROCObjID(), this.server.moduleid)

	// 由于ROC绑定消息与ROC调用之间存在异步问题，除非经过设置，
	// 否则使用同步方式同步ROC对象绑定
	if this.server.moduleConfig.GetBool(conf.AsynchronousSyncRocbind) {
		this.rocDelCacheChan <- obj
	} else {
		sendmsg := &servercomm.SROCBind{
			HostModuleID: this.server.moduleid,
			IsDelete:     true,
			ObjType:      string(obj.GetROCObjType()),
			ObjIDs:       []string{obj.GetROCObjID()},
		}
		this.sendROCBindMsg(sendmsg)
	}
	this.recordLocalObj(string(obj.GetROCObjType()), obj.GetROCObjID(),
		true)
}

// 向其他模块通知ROC注册信息的线程
func (this *ROCServer) rocObjNoticeProcess(rocCacheChan chan roc.IObj,
	isDelete bool) {
	tm := time.NewTimer(time.Millisecond * 300)
	tmpList := make([]roc.IObj, 100)
	lenTmpList := len(tmpList)
	tmpListI := 0
	var leftObj roc.IObj
	for !this.server.isStop {
		if tmpListI == 0 && leftObj != nil {
			tmpList[0] = leftObj
			tmpListI++
			leftObj = nil
		}
		wait := true
		for wait {
			wait = false
			select {
			case <-this.server.stopChan:
				break
			case rocObj := <-rocCacheChan:
				if tmpListI == 0 ||
					rocObj.GetROCObjType() == tmpList[0].GetROCObjType() {
					tmpList[tmpListI] = rocObj
					tmpListI++
					if tmpListI < lenTmpList {
						wait = true
					}
				} else {
					leftObj = rocObj
					break
				}
			case <-tm.C:
				tm.Reset(time.Millisecond * 300)
				break
			}
		}

		if !this.server.isStop {
			if tmpListI > 0 {
				sendmsg := &servercomm.SROCBind{
					HostModuleID: this.server.moduleid,
					IsDelete:     isDelete,
					ObjType:      string(tmpList[0].GetROCObjType()),
					ObjIDs:       make([]string, tmpListI),
				}
				for i, obj := range tmpList {
					if i >= tmpListI {
						break
					}
					sendmsg.ObjIDs[i] = string(obj.GetROCObjID())
					tmpList[i] = nil
				}
				this.sendROCBindMsg(sendmsg)
				tmpListI = 0
			}
		}
	}
}

// 发送ROC对象绑定信息
func (this *ROCServer) sendROCBindMsg(sendmsg *servercomm.SROCBind) {
	this.server.subnetManager.RangeServer(
		func(s *connect.Server) bool {
			if !process.HasModule(s.ModuleInfo.ModuleID) {
				s.SendCmd(sendmsg)
			}
			return true
		})
}

// 当收到ROC绑定信息时
func (this *ROCServer) onMsgROCBind(msg *servercomm.SROCBind) {
	if !msg.IsDelete {
		roc.GetCache().SetM(roc.ROCObjType(msg.ObjType), msg.ObjIDs, msg.HostModuleID)
	} else {
		roc.GetCache().DelM(roc.ROCObjType(msg.ObjType), msg.ObjIDs, msg.HostModuleID)
	}
	this.Syslog("onMsgROCBind roc cache setm type[%s] "+
		"ids%+v host[%s]",
		msg.ObjType, msg.ObjIDs, msg.HostModuleID)
}
