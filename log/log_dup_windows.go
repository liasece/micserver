package log

import (
	"syscall"
)

const (
	kernel32dll = "kernel32.dll"
)

// SysDup 将进程致命错误转储
func SysDup(fd int) {
	kernel32 := syscall.NewLazyDLL(kernel32dll)
	setStdHandle := kernel32.NewProc("SetStdHandle")
	// 把错误重定向到日志文件来
	_, _, e1 := setStdHandle.Call(uintptr(1), uintptr(fd))
	if e1 != nil {
		// return e1
	}
	_, _, e2 := setStdHandle.Call(uintptr(2), uintptr(fd))
	if e2 != nil {
		// return e2
	}
}
