/*
Package sysutil 系统当前的网络信息
*/
package sysutil

import (
	"net"
)

// GetIPv4ByInterface GetIPv4ByInterface return IPv4 address from a specific interface IPv4 addresses
func GetIPv4ByInterface(name string) string {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return ""
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return ""
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}

	return ""
}
