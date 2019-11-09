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
	fromServerID string
	callpath     string
	callarg      []byte
	seq          int64
	needReturn   bool
}

type responseAgent struct {
	fromServerID string
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
	go this.rocObjNoticeProcess()
	this._ROCManager.HookObjEvent(this)

	this.rocRequestChan = make(chan *requestAgent, 10000)
	go this.rocRequestProcess()
	this.rocResponseChan = make(chan *responseAgent, 10000)
	go this.rocResponseProcess()
}

func (this *ROCServer) newSeq() (res int64) {
	this.seqMutex.Lock()
	this.lastSeq++
	res = this.lastSeq
	this.seqMutex.Unlock()
	return
}

func (this *ROCServer) GetROC(objtype roc.ROCObjType) *roc.ROC {
	return this._ROCManager.GetROC(objtype)
}

func (this *ROCServer) NewROC(objtype roc.ROCObjType) {
	this._ROCManager.NewROC(objtype)
}

// 无返回值的RPC调用
func (this *ROCServer) ROCCallNR(callpath *roc.ROCPath, callarg []byte) {
	objType := callpath.GetObjType()
	objID := callpath.GetObjID()
	serverid := roc.GetCache().Get(objType, objID)
	this.Info("ROCCallNR {%s:%s(%s:%s):%v}",
		serverid, callpath, objType, objID, callarg)
	// 构造消息
	sendmsg := &servercomm.SROCRequest{
		FromServerID: this.server.serverid,
		Seq:          this.newSeq(),
		CallStr:      callpath.String(),
		CallArg:      callarg,
	}
	if serverid == this.server.serverid {
		sendmsg.ToServerID = serverid
		this.onMsgROCRequest(sendmsg)
	} else {
		server := this.server.subnetManager.GetServer(serverid)
		if server != nil {
			sendmsg.ToServerID = serverid
			server.SendCmd(sendmsg)
		} else {
			this.server.subnetManager.BroadcastCmd(sendmsg)
		}
	}
}

func (this *ROCServer) addBlockChan(seq int64) chan *responseAgent {
	ch := make(chan *responseAgent, 1)
	this.rocBlockChanMap.Store(seq, ch)
	return ch
}

// 无返回值的RPC调用
func (this *ROCServer) ROCCallBlock(callpath *roc.ROCPath,
	callarg []byte) ([]byte, error) {
	objType := callpath.GetObjType()
	objID := callpath.GetObjID()
	serverid := roc.GetCache().Get(objType, objID)
	this.Info("ROCCallBlock {%s:%s(%s:%s:%d):%v}",
		serverid, callpath, objType, objID, hash.GetStringHash(string(objID)),
		callarg)
	// 构造消息
	sendmsg := &servercomm.SROCRequest{
		FromServerID: this.server.serverid,
		Seq:          this.newSeq(),
		CallStr:      callpath.String(),
		CallArg:      callarg,
		NeedReturn:   true,
	}

	ch := this.addBlockChan(sendmsg.Seq)

	if serverid == this.server.serverid {
		sendmsg.ToServerID = serverid
		this.onMsgROCRequest(sendmsg)
	} else {
		server := this.server.subnetManager.GetServer(serverid)
		if server != nil {
			sendmsg.ToServerID = serverid
			server.SendCmd(sendmsg)
		} else {
			this.server.subnetManager.BroadcastCmd(sendmsg)
		}
	}

	agent := <-ch
	return agent.data, errors.New(agent.err)
}

func (this *ROCServer) onMsgROCRequest(msg *servercomm.SROCRequest) {
	agent := &requestAgent{
		callpath:     msg.CallStr,
		callarg:      msg.CallArg,
		seq:          msg.Seq,
		needReturn:   msg.NeedReturn,
		fromServerID: msg.FromServerID,
	}
	this.rocRequestChan <- agent
}

func (this *ROCServer) onMsgROCResponse(msg *servercomm.SROCResponse) {
	agent := &responseAgent{
		fromServerID: msg.FromServerID,
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
				this.Error("ROCManager.Call err:%s", err.Error())
			} else {
				// this.Debug("ROC调用成功 res:%+v", res)
			}
			if agent.needReturn {
				if agent.fromServerID == this.server.serverid {
					// 本地服务器调用结果

				} else {
					server := this.server.subnetManager.GetServer(
						agent.fromServerID)
					if server != nil {
						// 返回执行结果
						sendmsg := &servercomm.SROCResponse{
							FromServerID: this.server.serverid,
							ToServerID:   agent.fromServerID,
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
		this.server.serverid)
	this.Debug("OnROCObjAdd roc cache set type[%s] "+
		"id[%s] host[%s]",
		obj.GetROCObjType(), obj.GetROCObjID(), this.server.serverid)
	this.rocAddCacheChan <- obj
}

// 当ROC对象发生注册行为时
func (this *ROCServer) OnROCObjDel(obj roc.IObj) {
	// 保存本地映射缓存
	roc.GetCache().Del(obj.GetROCObjType(), obj.GetROCObjID(),
		this.server.serverid)
	this.Debug("OnROCObjDel roc cache del type[%s] "+
		"id[%s] host[%s]",
		obj.GetROCObjType(), obj.GetROCObjID(), this.server.serverid)
	this.rocDelCacheChan <- obj
}

func (this *ROCServer) rocObjNoticeProcess() {
	tm := time.NewTimer(time.Millisecond * 300)
	tmpAddList := make([]roc.IObj, 100)
	tmpDelList := make([]roc.IObj, 100)
	lenTmpList := len(tmpAddList)
	tmpAddListI := 0
	tmpDelListI := 0
	var leftAddObj roc.IObj
	var leftDelObj roc.IObj
	for !this.server.isStop {
		if tmpAddListI == 0 && leftAddObj != nil {
			tmpAddList[0] = leftAddObj
			tmpAddListI++
			leftAddObj = nil
		}
		if tmpDelListI == 0 && leftDelObj != nil {
			tmpDelList[0] = leftDelObj
			tmpDelListI++
			leftDelObj = nil
		}
		wait := true
		for wait {
			wait = false
			select {
			case <-this.server.stopChan:
				break
			case rocObj := <-this.rocAddCacheChan:
				if tmpAddListI == 0 ||
					rocObj.GetROCObjType() == tmpAddList[0].GetROCObjType() {
					tmpAddList[tmpAddListI] = rocObj
					tmpAddListI++
					if tmpAddListI < lenTmpList {
						wait = true
					}
				} else {
					leftAddObj = rocObj
					break
				}
			case rocObj := <-this.rocDelCacheChan:
				if tmpDelListI == 0 ||
					rocObj.GetROCObjType() == tmpDelList[0].GetROCObjType() {
					tmpDelList[tmpAddListI] = rocObj
					tmpDelListI++
					if tmpDelListI < lenTmpList {
						wait = true
					}
				} else {
					leftDelObj = rocObj
					break
				}
			case <-tm.C:
				tm.Reset(time.Millisecond * 300)
				break
			}
		}

		if !this.server.isStop {
			if tmpAddListI > 0 {
				sendmsg := &servercomm.SROCBind{
					HostServerID: this.server.serverid,
					IsDelete:     false,
					ObjType:      string(tmpAddList[0].GetROCObjType()),
					ObjIDs:       make([]string, tmpAddListI),
				}
				for i, obj := range tmpAddList {
					if i >= tmpAddListI {
						break
					}
					sendmsg.ObjIDs[i] = string(obj.GetROCObjID())
					tmpAddList[i] = nil
				}
				this.sendROCBindMsg(sendmsg)
				tmpAddListI = 0
			}
			if tmpDelListI > 0 {
				sendmsg := &servercomm.SROCBind{
					HostServerID: this.server.serverid,
					IsDelete:     true,
					ObjType:      string(tmpDelList[0].GetROCObjType()),
					ObjIDs:       make([]string, tmpDelListI),
				}
				for i, obj := range tmpDelList {
					if i >= tmpDelListI {
						break
					}
					sendmsg.ObjIDs[i] = string(obj.GetROCObjID())
					tmpDelList[i] = nil
				}
				this.sendROCBindMsg(sendmsg)
				tmpDelListI = 0
			}
		}
	}
}

// 发送ROC对象绑定信息
func (this *ROCServer) sendROCBindMsg(sendmsg *servercomm.SROCBind) {

	this.server.subnetManager.RangeServer(
		func(s *connect.Server) bool {
			if !process.HasModule(s.ServerInfo.ServerID) {
				s.SendCmd(sendmsg)
			}
			return true
		})
}

// 当收到ROC绑定信息时
func (this *ROCServer) onMsgROCBind(msg *servercomm.SROCBind) {
	if !msg.IsDelete {
		roc.GetCache().SetM(roc.ROCObjType(msg.ObjType), msg.ObjIDs, msg.HostServerID)
	} else {
		roc.GetCache().DelM(roc.ROCObjType(msg.ObjType), msg.ObjIDs, msg.HostServerID)
	}
	this.Debug("onMsgROCBind roc cache setm type[%s] "+
		"ids%+v host[%s]",
		msg.ObjType, msg.ObjIDs, msg.HostServerID)
}
