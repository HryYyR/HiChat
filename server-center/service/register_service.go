package service

import (
	"bufio"
	"fmt"
	"net"
	cconfig "server-center/config"
	"server-center/models"
	"strconv"

	consulapi "github.com/hashicorp/consul/api"
)

func RegisterService(dis models.DiscoveryConfig) error {
	config := consulapi.DefaultConfig()
	config.Address = cconfig.ConsulAddress
	client, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Printf("create consul client : %v\n", err.Error())
	}
	registration := &consulapi.AgentServiceRegistration{
		ID:      dis.ID,
		Name:    dis.Name,
		Port:    dis.Port,
		Tags:    dis.Tags,
		Address: dis.Address,
	}
	// 启动tcp的健康检测，注意address不能使用127.0.0.1或者localhost，因为consul-agent在docker容器里，如果用这个的话，
	// consul会访问容器里的port就会出错，一直检查不到实例
	check := &consulapi.AgentServiceCheck{}
	check.TCP = fmt.Sprintf("%s:%d", registration.Address, registration.Port)
	check.Timeout = "5s"
	check.Interval = "5s"
	check.DeregisterCriticalServiceAfter = "60s"
	registration.Check = check

	if err := client.Agent().ServiceRegister(registration); err != nil {
		fmt.Printf("register to consul error: %v\n", err.Error())
		return err
	}
	return nil
}

func StartTcp(tcpaddress int) {
	ls, err := net.Listen("tcp", fmt.Sprintf(":%s", strconv.Itoa(tcpaddress)))
	if err != nil {
		fmt.Printf("start tcp listener error: %v\n", err.Error())
		return
	}
	for {
		conn, err := ls.Accept()
		if err != nil {
			fmt.Printf("connect error: %v\n", err.Error())
		}
		go func(conn net.Conn) {
			_, err := bufio.NewWriter(conn).WriteString("hello consul")
			if err != nil {
				fmt.Printf("write conn error: %v\n", err)
			}
		}(conn)
	}
}
