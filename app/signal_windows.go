package app

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/liasece/micserver/log"
)

// SignalListen 监听系统消息
func (a *App) SignalListen() {
	c := make(chan os.Signal, 10)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT)
	for {
		s := <-c
		a.Debug("[App] Received os signal", log.Reflect("Signal", s))
		// manager.OnSignal(s)
		switch s {
		case syscall.SIGINT:
			// kill -2 || Ctrl+c 触发的信号，中断信号，执行正常退出操作
			a.isStoped <- struct{}{}
		case syscall.SIGQUIT:
			// kill -9 || Ctrl+z 触发的信号，强制杀死进程
			// 捕捉到就退不出了
			buf := make([]byte, 1<<20)
			stacklen := runtime.Stack(buf, true)
			a.Debug("[App] Received SIGQUIT", log.ByteString("Stack", buf[:stacklen]))
		}
	}
}
