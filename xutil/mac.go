// Package xutil
/*
 * @Date: 2023-07-20 10:23:00
 * @LastEditTime: 2023-07-20 10:23:01
 * @Description:
 */
package xutil

import (
	"fmt"
	"net"
	"runtime"
)

// GetLocalMacAddr
// @Description: get local mac address
// @return string
// @return error
// @author cx
func GetLocalMacAddr() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	switch runtime.GOOS {
	case "darwin":
		name, err := net.InterfaceByName("en0")
		if err != nil {
			return "", err
		}
		return name.HardwareAddr.String(), nil
	case "windows":
		for _, i := range interfaces {
			if i.Flags&net.FlagLoopback == 0 && len(i.HardwareAddr) > 0 {
				return i.HardwareAddr.String(), nil
			}
		}
	case "linux":
		var macAddr string
		for _, i := range interfaces {
			macAddr = i.HardwareAddr.String()
		}
		if macAddr != "" {
			return macAddr, nil
		}
	}
	return "", fmt.Errorf("no valid network interface found")
}
