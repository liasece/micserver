/*
Package sysutil 系统panic恢复
*/
package sysutil

import (
	"fmt"
	"runtime"
)

// GetPanicInfo 获取系统的panic信息
func GetPanicInfo(erri interface{}) (string, error) {
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
		return stackInfo, err
	}
	return "", nil
}
