package server

import (
	"github.com/liasece/micserver/roc"
	"github.com/liasece/micserver/servercomm"
	"github.com/liasece/micserver/util"
	"sync"
	"time"
)

type catchServerInfo struct {
	serverid string
}

type callAgent struct {
	callpath   string
	callarg    []byte
	seq        int64
	needReturn bool
}

type ROCServer struct {
	server *Server

	// 远程对象调用支持
	ROCManager   roc.ROCManager
	rocCatchChan chan roc.IObj
	rocCatchList sync.Map

	rocCallChan chan *callAgent

	catchType   sync.Map
	catchObj    sync.Map
	catchServer sync.Map
}

func (this *ROCServer) Init(server *Server) {
	this.server = server

	this.rocCatchChan = make(chan roc.IObj, 10000)
	go this.rocObjNoticeProcess()
	this.ROCManager.RegOnRegObj(this.onRegROCObj)

	this.rocCallChan = make(chan *callAgent, 10000)
	go this.rocCallProcess()
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
		Seq:          100,
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

func (this *ROCServer) onMsgROCRequest(msg *servercomm.SROCRequest) {
	agent := &callAgent{
		callpath:   msg.CallStr,
		callarg:    msg.CallArg,
		seq:        msg.Seq,
		needReturn: msg.NeedReturn,
	}
	this.rocCallChan <- agent
}

func (this *ROCServer) rocCallProcess() {
	tm := time.NewTimer(time.Millisecond * 300)
	for !this.server.isStop {
		select {
		case <-this.server.stopChan:
			break
		case agent := <-this.rocCallChan:
			// 处理ROC请求
			this.server.Info("ROCServer.rocCallProcess %+v", agent)
			res, err := this.ROCManager.Call(agent.callpath, agent.callarg)
			if err != nil {
				this.server.Error("ROCManager.Call err:%s", err.Error())
			} else {
				this.server.Info("ROC调用成功 res:%+v", res)
			}
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
