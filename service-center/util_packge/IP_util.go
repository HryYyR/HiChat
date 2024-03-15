package util_packge

import (
	"fmt"
	"net"
)

func GetIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		if ipnet.IP.To4() != nil && !ipnet.IP.IsLoopback() && ipnet.IP.String()[:3] != "169" { // IPv4 address
			ip := ipnet.IP.String()
			return ip
		}
	}
	return ""
}
