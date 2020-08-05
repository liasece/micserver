/*
Package server micserver中的ROC调用发生时，处理调用以及返回值。
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

// ROCServer ROC服务
type ROCServer struct {
	*log.Logger
	server *Server

	// 记录在本地的缓存信息
	// 第一层键为ROCObj类型，第二层键为ROCObj的ID
	localObj      map[string]map[string]struct{}
	localObjMutex sync.Mutex

	// 远程对象调用支持
	_ROCManager     roc.Manager
	rocAddCacheChan chan roc.IObj
	rocDelCacheChan chan roc.IObj

	rocRequestChan  chan *requestAgent
	rocResponseChan chan *responseAgent
	rocBlockChanMap sync.Map

	seqMutex sync.Mutex
	lastSeq  int64
}

// Init 初始化ROC服务
func (rocServer *ROCServer) Init(server *Server) {
	rocServer.server = server
	rocServer.Logger = server.Logger.Clone()
	rocServer.Logger.SetTopic("ROCServer")

	rocServer.rocAddCacheChan = make(chan roc.IObj, 10000)
	rocServer.rocDelCacheChan = make(chan roc.IObj, 10000)
	go rocServer.rocObjNoticeProcess(rocServer.rocAddCacheChan, false)
	go rocServer.rocObjNoticeProcess(rocServer.rocDelCacheChan, true)
	rocServer._ROCManager.HookObjEvent(rocServer)

	rocServer.rocRequestChan = make(chan *requestAgent, 10000)
	go rocServer.rocRequestProcess()
	rocServer.rocResponseChan = make(chan *responseAgent, 10000)
	go rocServer.rocResponseProcess()
}

// newSeq 生成一个 ROC 调用的序号，在每个模块中应该唯一
func (rocServer *ROCServer) newSeq() (res int64) {
	rocServer.seqMutex.Lock()
	rocServer.lastSeq++
	res = rocServer.lastSeq
	rocServer.seqMutex.Unlock()
	return
}

// GetROC 获取指定ROC对象类型的ROC对象
func (rocServer *ROCServer) GetROC(objtype roc.ObjType) *roc.ROC {
	return rocServer._ROCManager.GetROC(objtype)
}

// NewROC 新建一个指定ROC对象类型的ROC对象
func (rocServer *ROCServer) NewROC(objtype roc.ObjType) *roc.ROC {
	return rocServer._ROCManager.NewROC(objtype)
}

// ROCCallNR 无返回值的ROC调用
func (rocServer *ROCServer) ROCCallNR(callpath *roc.Path, callarg []byte) error {
	objType := callpath.GetObjType()
	objID := callpath.GetObjID()
	moduleid := roc.GetCache().Get(objType, objID)
	rocServer.Syslog("[ROCServer.ROCCallNR] ROCCallNR", log.String("TargetModuleID", moduleid), log.String("CallPath", callpath.String()),
		log.String("ObjType", string(objType)), log.String("ObjID", objID), log.Reflect("CallArg", callarg))
	// 构造消息
	sendmsg := &servercomm.SROCRequest{
		FromModuleID: rocServer.server.moduleid,
		Seq:          rocServer.newSeq(),
		CallStr:      callpath.String(),
		CallArg:      callarg,
	}
	if moduleid == rocServer.server.moduleid {
		sendmsg.ToModuleID = moduleid
		rocServer.onMsgROCRequest(sendmsg)
	} else {
		server := rocServer.server.subnetManager.GetServer(moduleid)
		if server != nil {
			sendmsg.ToModuleID = moduleid
			server.SendCmd(sendmsg)
		} else {
			rocServer.Warn("[ROCServer.ROCCallNR] Can't find roc object location", log.String("CallPath", callpath.String()))
			return fmt.Errorf("Can't find roc object location %s", callpath.String())
		}
	}
	return nil
}

// GetROCCachedLocation 获取ROC缓存中的位置信息
// 返回目标ROC对象所在的moduleid
func (rocServer *ROCServer) GetROCCachedLocation(objType roc.ObjType, objID string) string {
	moduleid := roc.GetCache().Get(objType, objID)
	return moduleid
}

// RangeROCCachedByType 遍历指定类型的ROC缓存，限制目标对象必须本module可以访问
func (rocServer *ROCServer) RangeROCCachedByType(objType roc.ObjType, f func(id string, location string) bool) {
	connecedModuleIDs := make(map[string]bool)
	rocServer.server.subnetManager.RangeServer(func(server *connect.Server) bool {
		if server.ModuleInfo != nil {
			connecedModuleIDs[server.ModuleInfo.ModuleID] = true
		}
		return true
	})
	roc.GetCache().RangeByType(objType, f, connecedModuleIDs)
}

// RandomROCCachedByType 随机获取本地缓存的ROC对象，返回该对象的ID，限制目标对象必须本module可以访问
func (rocServer *ROCServer) RandomROCCachedByType(objType roc.ObjType) string {
	connecedModuleIDs := make(map[string]bool)
	rocServer.server.subnetManager.RangeServer(func(server *connect.Server) bool {
		if server.ModuleInfo != nil {
			connecedModuleIDs[server.ModuleInfo.ModuleID] = true
		}
		return true
	})
	return roc.GetCache().RandomObjIDByType(objType, connecedModuleIDs)
}

// addBlockChan 根据ROC请求的序号，生成一个用于阻塞等待ROC返回的chan
func (rocServer *ROCServer) addBlockChan(seq int64) chan *responseAgent {
	ch := make(chan *responseAgent, 1)
	rocServer.rocBlockChanMap.Store(seq, ch)
	return ch
}

// ROCCallBlock 有返回值的RPC调用
func (rocServer *ROCServer) ROCCallBlock(callpath *roc.Path, callarg []byte) ([]byte, error) {
	objType := callpath.GetObjType()
	objID := callpath.GetObjID()
	moduleid := roc.GetCache().Get(objType, objID)
	rocServer.Syslog("[ROCServer.ROCCallBlock] ROCCallBlock", log.String("TargetModuleID", moduleid), log.String("CallPath", callpath.String()),
		log.String("ObjType", string(objType)), log.String("ObjID", objID), log.Uint32("ObjIDHash", hash.GetStringHash(string(objID))), log.Reflect("CallArg", callarg))
	// 构造消息
	sendmsg := &servercomm.SROCRequest{
		FromModuleID: rocServer.server.moduleid,
		Seq:          rocServer.newSeq(),
		CallStr:      callpath.String(),
		CallArg:      callarg,
		NeedReturn:   true,
	}

	ch := rocServer.addBlockChan(sendmsg.Seq)

	if moduleid == rocServer.server.moduleid {
		sendmsg.ToModuleID = moduleid
		rocServer.onMsgROCRequest(sendmsg)
	} else {
		server := rocServer.server.subnetManager.GetServer(moduleid)
		if server != nil {
			sendmsg.ToModuleID = moduleid
			server.SendCmd(sendmsg)
		} else {
			rocServer.server.subnetManager.BroadcastCmd(sendmsg)
		}
	}

	// 等待返回值
	agent := <-ch
	return agent.data, errors.New(agent.err)
}

// onMsgROCRequest 当收到ROC调用请求时
func (rocServer *ROCServer) onMsgROCRequest(msg *servercomm.SROCRequest) {
	agent := &requestAgent{
		callpath:     msg.CallStr,
		callarg:      msg.CallArg,
		seq:          msg.Seq,
		needReturn:   msg.NeedReturn,
		fromModuleID: msg.FromModuleID,
	}
	rocServer.rocRequestChan <- agent
}

// onMsgROCResponse 当收到ROC调用返回时
func (rocServer *ROCServer) onMsgROCResponse(msg *servercomm.SROCResponse) {
	agent := &responseAgent{
		fromModuleID: msg.FromModuleID,
		seq:          msg.ReqSeq,
		data:         msg.ResData,
		err:          msg.Error,
	}
	rocServer.rocResponseChan <- agent
}

// rocRequestProcess 处理带返回值的ROC调用返回的线程
func (rocServer *ROCServer) rocRequestProcess() {
	tm := time.NewTimer(time.Millisecond * 300)
	for !rocServer.server.isStop {
		select {
		case <-rocServer.server.stopChan:
			break
		case agent := <-rocServer.rocRequestChan:
			// 处理ROC请求
			rocServer.Syslog("[ROCServer.rocRequestProcess] ROC Request", log.String("CallPath", agent.callpath))
			res, err := rocServer._ROCManager.Call(agent.callpath, agent.callarg)
			if err != nil {
				if err != roc.ErrUnknowObj {
					rocServer.Error("[ROCServer.rocRequestProcess] Manager.Call", log.ErrorField(err))
				} else {
					rocServer.Syslog("[ROCServer.rocRequestProcess] Manager.Call", log.String("CallPath", agent.callpath), log.ErrorField(err))
				}
			} else {
				// rocServer.Debug("ROC调用成功 res:%+v", res)
			}
			if agent.needReturn {
				if agent.fromModuleID == rocServer.server.moduleid {
					// 本地服务器调用结果

				} else {
					server := rocServer.server.subnetManager.GetServer(agent.fromModuleID)
					if server != nil {
						// 返回执行结果
						sendmsg := &servercomm.SROCResponse{
							FromModuleID: rocServer.server.moduleid,
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

// rocResponseProcess 处理ROC带返回值电泳调用的线程
func (rocServer *ROCServer) rocResponseProcess() {
	tm := time.NewTimer(time.Millisecond * 300)
	for !rocServer.server.isStop {
		select {
		case <-rocServer.server.stopChan:
			break
		case agent := <-rocServer.rocResponseChan:
			// 处理ROC相应
			rocServer.Syslog("[ROCServer.rocResponseProcess] rocResponseProcess", log.Reflect("Agent", agent))
			chi, ok := rocServer.rocBlockChanMap.Load(agent.seq)
			if ok {
				if ch, ok := chi.(chan *responseAgent); ok {
					// 写入返回值
					ch <- agent
				} else {
					rocServer.Error("[ROCServer.rocResponseProcess] ROC returns chi.(chan *responseAgent) error", log.Reflect("Agent", agent))
				}
			} else {
				rocServer.Error("[ROCServer.rocResponseProcess] ROC return No target ROC request exists error", log.Reflect("Agent", agent))
			}
		case <-tm.C:
			tm.Reset(time.Millisecond * 300)
			break
		}
	}
}

// onServerJoinSubnet 当一个服务器加入了子网时
func (rocServer *ROCServer) onServerJoinSubnet(server *connect.Server) {
	if process.HasModule(server.ModuleInfo.ModuleID) {
		return
	}
	rocServer.localObjMutex.Lock()
	defer rocServer.localObjMutex.Unlock()

	if rocServer.localObj == nil {
		return
	}
	// 遍历所有类型的ROCObj
	for objtype, typemap := range rocServer.localObj {
		if typemap == nil || len(typemap) == 0 {
			continue
		}
		leftnum := len(typemap)
		// 临时变量
		tmplist := make([]string, 0)
		tmpsize := 0
		// 遍历所有对象
		for objid := range typemap {
			leftnum--
			tmplist = append(tmplist, objid)
			tmpsize += len(objid) + 4
			// 分包发送
			if leftnum == 0 ||
				(tmpsize > (msg.MessageMaxSize/2) && tmpsize > 32*1024) {
				sendmsg := &servercomm.SROCBind{
					HostModuleID: rocServer.server.moduleid,
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

// recordLocalObj 记录本地的ROC对象
func (rocServer *ROCServer) recordLocalObj(objtype string, objid string,
	isDelete bool) {
	rocServer.localObjMutex.Lock()
	defer rocServer.localObjMutex.Unlock()

	// 初始化内存
	if rocServer.localObj == nil {
		rocServer.localObj = make(map[string]map[string]struct{})
	}
	if v, ok := rocServer.localObj[objtype]; !ok || v == nil {
		rocServer.localObj[objtype] = make(map[string]struct{})
	}

	// 记录
	if isDelete {
		// 删除
		if _, ok := rocServer.localObj[objtype][objid]; ok {
			delete(rocServer.localObj[objtype], objid)
		}
	} else {
		rocServer.localObj[objtype][objid] = struct{}{}
	}
}

// OnROCObjAdd 当ROC对象发生注册行为时
func (rocServer *ROCServer) OnROCObjAdd(obj roc.IObj) {
	// 保存本地映射缓存
	roc.GetCache().Set(obj.GetROCObjType(), obj.GetROCObjID(),
		rocServer.server.moduleid)
	rocServer.Syslog("[ROCServer.OnROCObjAdd] Roc cache set", log.String("ObjType", string(obj.GetROCObjType())), log.String("ObjID", obj.GetROCObjID()),
		log.String("HostModuleID", rocServer.server.moduleid))

	// 由于ROC绑定消息与ROC调用之间存在异步问题，除非经过设置，
	// 否则使用同步方式同步ROC对象绑定
	if rocServer.server.moduleConfig.GetBool(conf.AsynchronousSyncRocbind) {
		rocServer.rocAddCacheChan <- obj
	} else {
		sendmsg := &servercomm.SROCBind{
			HostModuleID: rocServer.server.moduleid,
			IsDelete:     false,
			ObjType:      string(obj.GetROCObjType()),
			ObjIDs:       []string{obj.GetROCObjID()},
		}
		rocServer.sendROCBindMsg(sendmsg)
	}
	rocServer.recordLocalObj(string(obj.GetROCObjType()), obj.GetROCObjID(), false)
}

// OnROCObjDel 当ROC对象发生注册行为时
func (rocServer *ROCServer) OnROCObjDel(obj roc.IObj) {
	// 保存本地映射缓存
	roc.GetCache().Del(obj.GetROCObjType(), obj.GetROCObjID(), rocServer.server.moduleid)
	rocServer.Syslog("[ROCServer.OnROCObjDel] Roc cache del", log.String("ObjType", string(obj.GetROCObjType())), log.String("ObjID", obj.GetROCObjID()),
		log.String("HostModuleID", rocServer.server.moduleid))

	// 由于ROC绑定消息与ROC调用之间存在异步问题，除非经过设置，
	// 否则使用同步方式同步ROC对象绑定
	if rocServer.server.moduleConfig.GetBool(conf.AsynchronousSyncRocbind) {
		rocServer.rocDelCacheChan <- obj
	} else {
		sendmsg := &servercomm.SROCBind{
			HostModuleID: rocServer.server.moduleid,
			IsDelete:     true,
			ObjType:      string(obj.GetROCObjType()),
			ObjIDs:       []string{obj.GetROCObjID()},
		}
		rocServer.sendROCBindMsg(sendmsg)
	}
	rocServer.recordLocalObj(string(obj.GetROCObjType()), obj.GetROCObjID(), true)
}

// rocObjNoticeProcess 向其他模块通知ROC注册信息的线程
func (rocServer *ROCServer) rocObjNoticeProcess(rocCacheChan chan roc.IObj, isDelete bool) {
	tm := time.NewTimer(time.Millisecond * 300)
	tmpList := make([]roc.IObj, 100)
	lenTmpList := len(tmpList)
	tmpListI := 0
	var leftObj roc.IObj
	for !rocServer.server.isStop {
		if tmpListI == 0 && leftObj != nil {
			tmpList[0] = leftObj
			tmpListI++
			leftObj = nil
		}
		wait := true
		for wait {
			wait = false
			select {
			case <-rocServer.server.stopChan:
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

		if !rocServer.server.isStop {
			if tmpListI > 0 {
				sendmsg := &servercomm.SROCBind{
					HostModuleID: rocServer.server.moduleid,
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
				rocServer.sendROCBindMsg(sendmsg)
				tmpListI = 0
			}
		}
	}
}

// sendROCBindMsg 发送ROC对象绑定信息
func (rocServer *ROCServer) sendROCBindMsg(sendmsg *servercomm.SROCBind) {
	rocServer.server.subnetManager.RangeServer(
		func(s *connect.Server) bool {
			if !process.HasModule(s.ModuleInfo.ModuleID) {
				s.SendCmd(sendmsg)
			}
			return true
		})
}

// onMsgROCBind 当收到ROC绑定信息时
func (rocServer *ROCServer) onMsgROCBind(msg *servercomm.SROCBind) {
	if !msg.IsDelete {
		roc.GetCache().SetM(roc.ObjType(msg.ObjType), msg.ObjIDs, msg.HostModuleID)
	} else {
		roc.GetCache().DelM(roc.ObjType(msg.ObjType), msg.ObjIDs, msg.HostModuleID)
	}
	rocServer.Syslog("[ROCServer.onMsgROCBind] Roc cache setm", log.String("ObjType", string(msg.ObjType)), log.Reflect("ObjIDs", msg.ObjIDs),
		log.String("HostModuleID", msg.HostModuleID))
}
