package server

import (
	"errors"
	"sync"
	"time"

	"github.com/liasece/micserver/connect"
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/process"
	"github.com/liasece/micserver/roc"
	"github.com/liasece/micserver/servercomm"
	"github.com/liasece/micserver/util/hash"
)

type requestAgent struct {
	fromModuleID string
	callpath     string
	callarg      []byte
	seq          int64
	needReturn   bool
}

type responseAgent struct {
	fromModuleID string
	seq          int64
	data         []byte
	err          string
}

type ROCServer struct {
	*log.Logger
	server *Server

	// 远程对象调用支持
	_ROCManager     roc.ROCManager
	rocAddCacheChan chan roc.IObj
	rocDelCacheChan chan roc.IObj
	rocCatchList    sync.Map

	rocRequestChan  chan *requestAgent
	rocResponseChan chan *responseAgent
	rocBlockChanMap sync.Map

	seqMutex sync.Mutex
	lastSeq  int64
}

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
func (this *ROCServer) ROCCallNR(callpath *roc.ROCPath, callarg []byte) {
	objType := callpath.GetObjType()
	objID := callpath.GetObjID()
	moduleid := roc.GetCache().Get(objType, objID)
	this.Info("ROCCallNR {%s:%s(%s:%s):%X}",
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
		}
	}
}

// 获取ROC缓存中的位置信息
// 返回目标ROC对象所在的moduleid
func (this *ROCServer) GetROCObjCacheLocation(path *roc.ROCPath) string {
	objType := path.GetObjType()
	objID := path.GetObjID()
	moduleid := roc.GetCache().Get(objType, objID)
	return moduleid
}

// 遍历指定类型的ROC缓存
func (this *ROCServer) RangeROCObjIDByType(objType roc.ROCObjType,
	f func(id string, location string) bool) {
	roc.GetCache().RangeByType(objType, f, nil)
}

// 遍历指定类型的ROC缓存，限制目标对象必须本module可以访问
func (this *ROCServer) RangeMyROCObjIDByType(objType roc.ROCObjType,
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

// 遍历指定类型的ROC缓存，限制目标对象必须本module可以访问
func (this *ROCServer) RandomMyROCObjIDByType(objType roc.ROCObjType) string {
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
	this.Info("ROCCallBlock {%s:%s(%s:%s:%d):%X}",
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

func (this *ROCServer) rocRequestProcess() {
	tm := time.NewTimer(time.Millisecond * 300)
	for !this.server.isStop {
		select {
		case <-this.server.stopChan:
			break
		case agent := <-this.rocRequestChan:
			// 处理ROC请求
			this.Info("ROC Request[%s]", agent.callpath)
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

func (this *ROCServer) rocResponseProcess() {
	tm := time.NewTimer(time.Millisecond * 300)
	for !this.server.isStop {
		select {
		case <-this.server.stopChan:
			break
		case agent := <-this.rocResponseChan:
			// 处理ROC相应
			this.Info("rocResponseProcess %+v", agent)
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

// 当ROC对象发生注册行为时
func (this *ROCServer) OnROCObjAdd(obj roc.IObj) {
	// 保存本地映射缓存
	roc.GetCache().Set(obj.GetROCObjType(), obj.GetROCObjID(),
		this.server.moduleid)
	this.Debug("OnROCObjAdd roc cache set type[%s] "+
		"id[%s] host[%s]",
		obj.GetROCObjType(), obj.GetROCObjID(), this.server.moduleid)
	this.rocAddCacheChan <- obj
}

// 当ROC对象发生注册行为时
func (this *ROCServer) OnROCObjDel(obj roc.IObj) {
	// 保存本地映射缓存
	roc.GetCache().Del(obj.GetROCObjType(), obj.GetROCObjID(),
		this.server.moduleid)
	this.Debug("OnROCObjDel roc cache del type[%s] "+
		"id[%s] host[%s]",
		obj.GetROCObjType(), obj.GetROCObjID(), this.server.moduleid)
	this.rocDelCacheChan <- obj
}

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
	this.Debug("onMsgROCBind roc cache setm type[%s] "+
		"ids%+v host[%s]",
		msg.ObjType, msg.ObjIDs, msg.HostModuleID)
}
