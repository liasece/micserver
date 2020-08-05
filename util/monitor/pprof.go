package monitor

import (
	"fmt"
	"net/http"
	"time"

	// pprof
	_ "net/http/pprof"

	"github.com/liasece/micserver/log"
	"github.com/liasece/micserver/util/sysutil"
)

// BindPprof func
func BindPprof(ip string, port uint32) error {
	log.Syslog("[BindPprof] Pprof is starting", log.String("IP", ip), log.Uint32("Port", port))
	go startPprofThread(ip, port)
	return nil
}

func startPprofThread(ip string, port uint32) {
	defer func() {
		// 必须要先声明defer，否则不能捕获到panic异常
		if stackInfo, err := sysutil.GetPanicInfo(recover()); err != nil {
			log.Error("[startPprofThread] Panic", log.ErrorField(err), log.String("Stack", stackInfo))
		}
	}()
	time.Sleep(500 * time.Millisecond)
	log.Syslog("[startPprofThread] Pprof is starting", log.String("IP", ip), log.Uint32("Port", port))
	pprofportstr := fmt.Sprintf("%s:%d", ip, port)
	err := http.ListenAndServe(pprofportstr, nil)
	if err == nil {
		log.Syslog("[startPprofThread] Pprof started successfully", log.String("Addr", pprofportstr))
	} else {
		log.Error("[startPprofThread] Pprof started error", log.String("Addr", pprofportstr), log.ErrorField(err))
	}
}
