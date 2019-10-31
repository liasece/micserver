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
	"github.com/liasece/micserver/util"
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
	ROCManager   roc.ROCManager
	rocCatchChan chan roc.IObj
	rocCatchList sync.Map

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

	this.rocCatchChan = make(chan roc.IObj, 10000)
	go this.rocObjNoticeProcess()
	this.ROCManager.RegOnRegObj(this.onRegROCObj)

	this.rocRequestChan = make(chan *requestAgent, 10000)
	go this.rocRequestProcess()
	this.rocResponseChan = make(chan *responseAgent, 10000)
	go this.rocResponseProcess()
}

func (this *ROCServer) NewSeq() (res int64) {
	this.seqMutex.Lock()
	this.lastSeq++
	res = this.lastSeq
	this.seqMutex.Unlock()
	return
}

// 无返回值的RPC调用
func (this *ROCServer) ROCCallNR(callpath string, callarg []byte) {
	objType, objID := this.ROCManager.CallPathDecode(callpath)
	serverid := roc.GetCache().Get(objType, objID)
	this.Info("ROCCallNR {%s:%s(%s:%s):%v}",
		serverid, callpath, objType, objID, callarg)
	// 构造消息
	sendmsg := &servercomm.SROCRequest{
		FromServerID: this.server.serverid,
		Seq:          this.NewSeq(),
		CallStr:      callpath,
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

func (this *ROCServer) AddBlockChan(seq int64) chan *responseAgent {
	ch := make(chan *responseAgent, 1)
	this.rocBlockChanMap.Store(seq, ch)
	return ch
}

// 无返回值的RPC调用
func (this *ROCServer) ROCCallBlock(callpath string,
	callarg []byte) ([]byte, error) {
	objType, objID := this.ROCManager.CallPathDecode(callpath)
	serverid := roc.GetCache().Get(objType, objID)
	this.Info("ROCCallBlock {%s:%s(%s:%s:%d):%v}",
		serverid, callpath, objType, objID, util.GetStringHash(objID), callarg)
	// 构造消息
	sendmsg := &servercomm.SROCRequest{
		FromServerID: this.server.serverid,
		Seq:          this.NewSeq(),
		CallStr:      callpath,
		CallArg:      callarg,
		NeedReturn:   true,
	}

	ch := this.AddBlockChan(sendmsg.Seq)

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
			this.Info("rocRequestProcess %+v", agent)
			res, err := this.ROCManager.Call(agent.callpath, agent.callarg)
			if err != nil {
				this.Error("ROCManager.Call err:%s", err.Error())
			} else {
				this.Info("ROC调用成功 res:%+v", res)
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
func (this *ROCServer) onRegROCObj(obj roc.IObj) {
	// 保存本地映射缓存
	roc.GetCache().Set(obj.GetObjType(), obj.GetObjID(),
		this.server.serverid)
	this.Debug("onRegROCObj roc cache set type[%s] "+
		"id[%s] host[%s]",
		obj.GetObjType(), obj.GetObjID(), this.server.serverid)
	this.rocCatchChan <- obj
}

func (this *ROCServer) rocObjNoticeProcess() {
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
			case rocObj := <-this.rocCatchChan:
				if tmpListI == 0 || rocObj.GetObjType() == tmpList[0].GetObjType() {
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
		if !this.server.isStop && tmpListI > 0 {
			sendmsg := &servercomm.SROCBind{
				HostServerID: this.server.serverid,
				IsDelete:     false,
				ObjType:      tmpList[0].GetObjType(),
				ObjIDs:       make([]string, tmpListI),
			}
			for i, obj := range tmpList {
				if i >= tmpListI {
					break
				}
				sendmsg.ObjIDs[i] = obj.GetObjID()
				tmpList[i] = nil
			}
			this.server.subnetManager.RangeServer(func(s *connect.Server) bool {
				if !process.HasModule(s.ServerInfo.ServerID) {
					s.SendCmd(sendmsg)
				}
				return true
			})
			tmpListI = 0
		}
	}
}

func (this *ROCServer) onMsgROCBind(msg *servercomm.SROCBind) {
	roc.GetCache().SetM(msg.ObjType, msg.ObjIDs, msg.HostServerID)
	this.Debug("onMsgROCBind roc cache setm type[%s] "+
		"ids%+v host[%s]",
		msg.ObjType, msg.ObjIDs, msg.HostServerID)
}
