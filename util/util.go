/*
Package util micserver中的一些算法及实用工具
*/
package util

import (
	"fmt"
	"net"
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

// GetFreePort asks the kernel for free open ports that are ready to use.
func GetFreePort() (int, error) {
	l, err := GetFreePorts(1)
	if err != nil || len(l) < 1 {
		return 0, err
	}
	return l[0], err
}

// GetFreePorts asks the kernel for free open ports that are ready to use.
func GetFreePorts(count int) ([]int, error) {
	var ports []int
	for i := 0; i < count; i++ {
		addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
		if err != nil {
			return nil, err
		}

		l, err := net.ListenTCP("tcp", addr)
		if err != nil {
			return nil, err
		}
		defer l.Close()
		ports = append(ports, l.Addr().(*net.TCPAddr).Port)
	}
	return ports, nil
}
