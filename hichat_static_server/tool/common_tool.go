package tool

import (
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net"
	"time"
)

var tstr = "2006-01-02 15:04:05"

func FormatTime(t time.Time) string {

	return t.Format(tstr)
}

func ParseTime(t string) (time.Time, error) {
	parse, err := time.Parse(tstr, t)
	if err != nil {
		return time.Time{}, err
	}
	return parse, nil
}

func FormatTampTime(tamptime *timestamppb.Timestamp) string {
	return tamptime.AsTime().Format(tstr)
}

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
