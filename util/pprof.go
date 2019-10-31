package util

import (
	"fmt"
	"github.com/liasece/micserver/log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func BindPprof(ip string, port uint32) error {
	log.Debug("[startPprofThread] pprof正在启动 IPPort[%s:%d]", ip, port)
	go startPprofThread(ip, port)
	return nil
}

func startPprofThread(ip string, port uint32) {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := GetPanicInfo(recover()); err != nil {
			log.Error("[startPprofThread] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
		}
	}()
	time.Sleep(500 * time.Millisecond)
	log.Debug("[startPprofThread] pprof开始监听 IPPort[%s:%d]", ip, port)
	pprofportstr := fmt.Sprintf("%s:%d", ip, port)
	err := http.ListenAndServe(pprofportstr, nil)
	if err == nil {
		log.Debug("[startPprofThread] pprof启动成功 Addr[%s]",
			pprofportstr)
	} else {
		log.Error("[startPprofThread] pprof启动失败 Addr[%s] Err[%s]",
			pprofportstr, err.Error())
	}
}
