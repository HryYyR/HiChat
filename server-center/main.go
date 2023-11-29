package main

import (
	"fmt"
	"os"
	"server-center/models"
	"server-center/service"
	"strconv"
)

func main() {
	var ServerPort int
	Args := os.Args[1:]
	if len(Args) == 0 {
		panic("please input ServerPort\nexmple:go run . 10111")
	}
	ServerPort, err := strconv.Atoi(Args[0])
	if err != nil {
		panic("ServerPort must be a number")
	}

	ch := make(chan error)
	dis := models.DiscoveryConfig{
		ID:      Args[0],
		Name:    " ",
		Tags:    []string{"a", "b"},
		Port:    ServerPort,
		Address: "100.98.64.254", //通过ifconfig查看本机的eth0的ipv4地址
	}
	go service.StartTcp(ServerPort)
	service.RegisterService(dis)
	fmt.Println("Service started to:", ServerPort)
	// 阻塞等待
	<-ch
}
