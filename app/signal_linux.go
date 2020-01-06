package app

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

// 监听系统信号，对 Ctrl+c 等命令做出反应，通知App阻塞处执行退出逻辑
func (this *App) SignalListen() {
	c := make(chan os.Signal, 10)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT)
	for {
		s := <-c
		this.Debug("[App] "+
			"Get signal Signal[%d]", s)
		// manager.OnSignal(s)
		switch s {
		case syscall.SIGINT:
			// kill -2 || Ctrl+c 触发的信号，中断信号，执行正常退出操作
			this.isStoped <- struct{}{}
		case syscall.SIGQUIT:
			// kill -9 || Ctrl+z 触发的信号，强制杀死进程
			// 捕捉到就退不出了
			buf := make([]byte, 1<<20)
			stacklen := runtime.Stack(buf, true)
			this.Debug("[App] "+
				"Received SIGQUIT, \n: Stack[%s]", buf[:stacklen])
		}
	}
}
