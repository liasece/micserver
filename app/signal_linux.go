package app

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

// 监听系统消息
func (this *App) SignalListen() {
	c := make(chan os.Signal, 10)
	signal.Notify(c, syscall.SIGINT,
		syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1,
		syscall.SIGUSR2)
	for {
		s := <-c
		this.Debug("[SubNetManager.SignalListen] "+
			"Get signal Signal[%d]", s)
		// manager.OnSignal(s)
		switch s {
		case syscall.SIGUSR1:
			go this.startTestCpuProfile()
		case syscall.SIGUSR2:
		case syscall.SIGTERM:
		case syscall.SIGINT:
			// kill -2 || Ctrl+c 触发的信号，中断信号，执行正常退出操作
			this.isStoped = true
		case syscall.SIGQUIT:
			// kill -9 || Ctrl+z 触发的信号，强制杀死进程
			// 捕捉到就退不出了
			buf := make([]byte, 1<<20)
			stacklen := runtime.Stack(buf, true)
			this.Debug("[SubNetManager.SignalListen] "+
				"Received SIGQUIT, \n: Stack[%s]", buf[:stacklen])
		}
	}
}
