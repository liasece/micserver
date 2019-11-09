package util

import (
	"fmt"
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
