package base

import (
	"base/logger"
	"base/util"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func BindPprof(ip string, port uint32) error {
	logger.Debug("[startPprofThread] pprof正在启动 IPPort[%s:%d]", ip, port)
	go startPprofThread(ip, port)
	return nil
}

func startPprofThread(ip string, port uint32) {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if err, stackInfo := util.GetPanicInfo(recover()); err != nil {
			logger.Error("[startPprofThread] "+
				"Panic: Err[%v] \n Stack[%s]", err, stackInfo)
		}
	}()
	time.Sleep(1 * time.Second)
	logger.Debug("[startPprofThread] pprof开始监听 IPPort[%s:%d]", ip, port)
	pprofportstr := fmt.Sprintf("%s:%d", ip, port)
	err := http.ListenAndServe(pprofportstr, nil)
	if err == nil {
		logger.Debug("[startPprofThread] pprof启动成功 Addr[%s]",
			pprofportstr)
	} else {
		logger.Error("[startPprofThread] pprof启动失败 Addr[%s] Err[%s]",
			pprofportstr, err.Error())
		GetGBServerConfigM().hasConfigPprof = false
	}
}
