package log

import (
	syslog "log"
	"syscall"
)

// sysDup 将进程致命错误转储
func sysDup(fd int) error {
	// Duplicate the stdin/stdout/stderr handles
	files := []uintptr{uintptr(syscall.Stdin), uintptr(syscall.Stdout), uintptr(syscall.Stderr)}
	p, _ := syscall.GetCurrentProcess()
	h := syscall.Handle(fd)
	for i := range files {
		err := syscall.DuplicateHandle(p, syscall.Handle(files[i]), p, &h, 0, true, syscall.DUPLICATE_SAME_ACCESS)
		if err != nil {
			syslog.Println(err)
			return err
		}
	}
	return nil
}
