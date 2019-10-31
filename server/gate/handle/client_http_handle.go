package handle

import (
	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/util/monitor"
	"net/http"
)

type ClientHttpHandler struct {
	*log.Logger
}

func PingHandler(writer http.ResponseWriter, request *http.Request) {
	functiontime := monitor.FunctionTime{}
	functiontime.Start("IDIPHandler")
	defer functiontime.Stop()

	writer.Header().Add("Access-Control-Allow-Origin", "*")
	writer.WriteHeader(200)
}

func (this *ClientHttpHandler) StartAddHttpHandle(addr string) {
	this.Debug("ClientHttpHandler.StartAddHttpHandle Addr[%s]", addr)
	// gm接口
	http.HandleFunc("/ping", PingHandler)

	// 开始监听
	go http.ListenAndServe(addr, nil)
}
