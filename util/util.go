/*
Package util micserver中的一些算法及实用工具
*/
package util

import (
	"fmt"
	"strconv"
)

// GetModuleIDType 获取模块ID中的模块类型部分，如 gate002 类型就是 gate
func GetModuleIDType(id string) string {
	res := ""
	for _, k := range id {
		if k >= '0' && k <= '9' {
			return res
		}
		res = res + fmt.Sprintf("%c", k)
	}
	return res
}

// GetModuleIDNum 获取模块ID中的模块序号，如 gate002 序号就是 2
func GetModuleIDNum(id string) int {
	numstr := ""
	for _, k := range id {
		if k >= '0' && k <= '9' {
			// 数字
			numstr = numstr + fmt.Sprintf("%c", k)
		} else {
			if numstr != "" {
				break
			}
		}
	}
	if numstr == "" {
		return 0
	}
	num, _ := strconv.Atoi(numstr)
	return num
}
