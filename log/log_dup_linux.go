package log

import (
	"syscall"
)

func SysDup(fd int) {
	e1 := syscall.Dup2(fd, 1)
	if e1 != nil {
		// return e1
	}
	e2 := syscall.Dup2(fd, 2)
	if e2 != nil {
		// return e2
	}
}
