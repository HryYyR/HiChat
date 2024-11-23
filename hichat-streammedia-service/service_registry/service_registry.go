package service_registry

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"hichat-streammedia-service/config"
)

type DiscoveryConfig struct {
	ID      string
	Name    string
	Tags    []string
	Port    int
	Address string
}

func RegisterService(dis DiscoveryConfig) error {
	clientconfig := consulapi.DefaultConfig()
	clientconfig.Address = config.ConsulAddress
	client, err := consulapi.NewClient(clientconfig)
	if err != nil {
		fmt.Printf("create consul client : %v\n", err.Error())
	}

	// 查询旧的服务实例并注销
	services, _, err := client.Catalog().Service(dis.Name, "", nil)
	if err != nil {
		return fmt.Errorf("query services from consul: %v", err)
	}
	for _, service := range services {
		if service.ServiceAddress == dis.Address && service.ServicePort == dis.Port {
			// 注销旧的服务实例
			err := client.Agent().ServiceDeregister(service.ServiceID)
			if err != nil {
				return fmt.Errorf("deregister service from consul: %v", err)
			}
			fmt.Printf("注销旧的服务实例: %s\n", service.ServiceID)
			break
		}
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
