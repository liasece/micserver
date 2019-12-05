package util

import (
	"fmt"
	"strconv"
)

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
