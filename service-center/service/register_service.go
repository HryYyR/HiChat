package service

import (
	cconfig "HiChat/service-center/config"
	"HiChat/service-center/models"
	"fmt"
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
