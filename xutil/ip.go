/*
 * @Date: 2023-07-20 09:34:46
 * @LastEditTime: 2023-07-20 10:26:35
 * @Description:
 */
package xutil

import (
	"errors"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
)

// GetIntranetIP 获取内网 IP
func GetIntranetIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// GetExternalIP 获取公网 IP
func GetExternalIP() string {
	// http://icanhazip.com/
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return string(b)
}

// ConvertToIntIP 转换ip为int
func ConvertToIntIP(ip string) (int, error) {
	ips := strings.Split(ip, ".")
	E := errors.New("not A IP")
	if len(ips) != 4 {
		return 0, E
	}
	var intIP int
	for k, v := range ips {
		i, err := strconv.Atoi(v)
		if err != nil || i > 255 {
			return 0, E
		}
		intIP = intIP | i<<uint(8*(3-k))
	}
	return intIP, nil
}

// GetLocalIpToInt 获取本机IP转成int
func GetLocalIpToInt() (int, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return 0, err
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ConvertToIntIP(ipnet.IP.String())
			}
		}
	}
	return 0, errors.New("can not find the client ip address")
}
