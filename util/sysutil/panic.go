package sysutil

import (
	"fmt"
	"runtime"
)

func GetPanicInfo(erri interface{}) (error, string) {
	if erri != nil {
		var err error
		switch erri.(type) {
		case error:
			err = erri.(error)
		case string:
			err = fmt.Errorf(erri.(string))
		default:
			err = fmt.Errorf("%+v", erri)
		}
		stackInfo := ""
		buf := make([]byte, 4*1024)
		n := runtime.Stack(buf, false)
		stackInfo += fmt.Sprintf("%s", buf[:n])
		return err, stackInfo
	}
	return nil, ""
}
