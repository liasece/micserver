package sysutil

import (
	"fmt"
	"runtime"
)

func GetPanicInfo(err interface{}) (error, string) {
	if err != nil {
		stackInfo := ""
		buf := make([]byte, 4*1024)
		n := runtime.Stack(buf, false)
		stackInfo += fmt.Sprintf("%s", buf[:n])
		return err.(error), stackInfo
	}
	return nil, ""
}
