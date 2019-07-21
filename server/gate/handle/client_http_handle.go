package handle

import (
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/util"
	"net/http"
)

type ClientHttpHandler struct {
	*log.Logger
}

func IDIPHandler(writer http.ResponseWriter, request *http.Request) {
	functiontime := util.FunctionTime{}
	functiontime.Start("IDIPHandler")
	defer functiontime.Stop()

	writer.Header().Add("Access-Control-Allow-Origin", "*")
	writer.WriteHeader(200)

	// bcontent, _ := ioutil.ReadAll(request.Body)
	// content := string(bcontent)
	// msg := &gmmsg.GMReqOriginal{}
	// json.Unmarshal(bcontent, msg)
	// msg.Decode()
	// this.Debug("[GM] Recv Content[%s] Msg[%s]", content, msg.Body)

	// switch msg.Head.CmdID {
	// case 0x102b:
	// 	recvmsg := &gmmsg.UserInfoReq{}
	// 	json.Unmarshal([]byte(msg.Body), recvmsg)
	// 	// forward.ForwardGMMsgToUserServer(writer, recvmsg.UUID, recvmsg.OpenID,
	// 	// msg.Head.CmdID, msg.Body)
	// case 0x1013:
	// 	recvmsg := &gmmsg.SysMsgReq{}
	// 	json.Unmarshal([]byte(msg.Body), recvmsg)

	// 	sysmsg := &jsonmsg.USystemMessage{}
	// 	sysmsg.Systemtype = jsonmsg.TypeSysInfoText
	// 	sysmsg.Systeminfo = recvmsg.NoticeContent
	// 	// forward.BroadcastToUser(sysmsg)
	// 	subnet.WriterReturnHttpStr(writer,
	// 		gmmsg.ErrorMsgJson(1, "操作成功"))
	// case 0x1040:
	// 	recvmsg := &gmmsg.SetGMLevelReq{}
	// 	json.Unmarshal([]byte(msg.Body), recvmsg)
	// 	// forward.ForwardGMMsgToUserServer(writer, recvmsg.UUID, recvmsg.OpenID,
	// 	// 	msg.Head.CmdID, msg.Body)
	// default:
	// 	subnet.WriterReturnHttpStr(writer, gmmsg.ErrorMsgJson(-4,
	// 		fmt.Sprintf("服务器未能处理的指令:0x%X", msg.Head.CmdID)))
	// }
}

func (this *ClientHttpHandler) StartAddHttpHandle(addr string) {
	this.Debug("Gateway 对外GM HTTP服务启动成功 IPPort[%s]", addr)
	// gm接口
	http.HandleFunc("/idip", IDIPHandler)

	// 开始监听
	go http.ListenAndServe(addr, nil)
}
