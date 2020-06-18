package log

import (
	syslog "log"
	"syscall"
)

// sysDup 将进程致命错误转储
func sysDup(fd int) error {
	// {
	// 	err := syscall.Dup2(fd, 1)
	// 	if err != nil {
	// 		syslog.Println(err.Error())
	// 		// return err
	// 	}
	// }
	{
		err := syscall.Dup2(fd, 2)
		if err != nil {
			syslog.Println(err.Error())
			// return err
		}
	}
	return nil
}
