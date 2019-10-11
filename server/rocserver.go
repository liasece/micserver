package server

import (
	"errors"
	"github.com/liasece/micserver/roc"
	"github.com/liasece/micserver/servercomm"
	"github.com/liasece/micserver/util"
	"sync"
	"time"
)

type catchServerInfo struct {
	serverid string
}

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

	catchType   sync.Map
	catchObj    sync.Map
	catchServer sync.Map
}

func (this *ROCServer) Init(server *Server) {
	this.server = server

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
	serverid := this.catchGet(objType, util.GetStringHash(objID))
	this.server.Info("ROCServer.ROCCallNR {%s:%s(%s:%s:%d):%v}",
		serverid, callpath, objType, objID, util.GetStringHash(objID), callarg)
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
	serverid := this.catchGet(objType, util.GetStringHash(objID))
	this.server.Info("ROCServer.ROCCallBlock {%s:%s(%s:%s:%d):%v}",
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
			this.server.Info("ROCServer.rocRequestProcess %+v", agent)
			res, err := this.ROCManager.Call(agent.callpath, agent.callarg)
			if err != nil {
				this.server.Error("ROCManager.Call err:%s", err.Error())
			} else {
				this.server.Info("ROC调用成功 res:%+v", res)
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
			this.server.Info("ROCServer.rocResponseProcess %+v", agent)
		case <-tm.C:
			tm.Reset(time.Millisecond * 300)
			break
		}
	}
}

// 当ROC对象发生注册行为时
func (this *ROCServer) onRegROCObj(obj roc.IObj) {
	// 保存本地映射缓存
	this.catch(obj.GetObjType(), util.GetStringHash(obj.GetObjID()),
		this.server.serverid)
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
				ObjIDHashs:   make([]uint32, tmpListI),
			}
			for i, obj := range tmpList {
				if i >= tmpListI {
					break
				}
				hash := util.GetStringHash(obj.GetObjID())
				sendmsg.ObjIDHashs[i] = hash
				tmpList[i] = nil
			}
			this.server.subnetManager.BroadcastCmd(sendmsg)
			tmpListI = 0
		}
	}
}

func (this *ROCServer) onMsgROCBind(msg *servercomm.SROCBind) {
	this.catchSet(msg.ObjType, msg.ObjIDHashs, msg.HostServerID)
}

func (this *ROCServer) catchGetTypeMust(objType string) *sync.Map {
	if vi, ok := this.catchType.Load(objType); !ok {
		vi, _ := this.catchType.LoadOrStore(objType, &sync.Map{})
		return vi.(*sync.Map)
	} else {
		return vi.(*sync.Map)
	}
}

func (this *ROCServer) catchGetServerMust(serverid string) *catchServerInfo {
	if vi, ok := this.catchServer.Load(serverid); !ok {
		vi, _ := this.catchServer.LoadOrStore(serverid, &catchServerInfo{
			serverid: serverid,
		})
		return vi.(*catchServerInfo)
	} else {
		return vi.(*catchServerInfo)
	}
}

func (this *ROCServer) catch(objType string, objIDHash uint32, serverid string) {
	m := this.catchGetTypeMust(objType)
	server := this.catchGetServerMust(serverid)
	m.Store(objIDHash, server)
	this.server.Debug("ROCServer.catchSet [%s] [%d] [%s]",
		objType, objIDHash, serverid)
}

func (this *ROCServer) catchSet(objType string, objIDHashs []uint32, serverid string) {
	m := this.catchGetTypeMust(objType)
	server := this.catchGetServerMust(serverid)
	for _, v := range objIDHashs {
		m.Store(v, server)
	}
	this.server.Debug("ROCServer.catchSet [%s] [%v] [%s]",
		objType, objIDHashs, serverid)
}

func (this *ROCServer) catchGet(objType string, objIDHash uint32) string {
	m := this.catchGetTypeMust(objType)
	if vi, ok := m.Load(objIDHash); ok && vi != nil {
		return (vi.(*catchServerInfo)).serverid
	}
	return ""
}
