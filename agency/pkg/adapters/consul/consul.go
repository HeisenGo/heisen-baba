package consul

import (
	"fmt"
	consulAPI "github.com/hashicorp/consul/api"
)

type Consul struct {
	Address string
}

func NewConsul(address string) *Consul {
	return &Consul{Address: address}
}

func (c *Consul) RegisterService(serviceName, serviceHostAddress, servicePrefixPath, serviceHTTPHealthPath string, serviceHTTPPort int) error {
	consulConfig := consulAPI.DefaultConfig()
	consulConfig.Address = c.Address
	consulClient, err := consulAPI.NewClient(consulConfig)
	if err != nil {
		return err
	}

	HTTPHealthURL := fmt.Sprintf("http://%s:%v/%s", serviceHostAddress, serviceHTTPPort, serviceHTTPHealthPath)
	// Register service with Consul
	registration := &consulAPI.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-service-id", serviceName),
		Name:    serviceName,
		Address: serviceHostAddress,
		Port:    serviceHTTPPort,
		Tags: []string{
			serviceName,
			fmt.Sprintf("traefik.http.routers.%s_router.rule=PathPrefix(`%s`)", serviceName, servicePrefixPath),
			fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.port=%v", serviceName, serviceHTTPPort),
		},
		Checks: []*consulAPI.AgentServiceCheck{
			{
				HTTP:     HTTPHealthURL,
				Interval: "10s",
				Timeout:  "1s",
			},
		},
	}

	err = consulClient.Agent().ServiceRegister(registration)
	if err != nil {
		return err
	}
	return nil
}

func (c *Consul) DiscoverService(serviceName string) (port int, ip string, err error) {
	consulConfig := consulAPI.DefaultConfig()
	consulConfig.Address = c.Address
	consulClient, err := consulAPI.NewClient(consulConfig)
	if err != nil {
		return 0, "", err
	}
	services, _, err := consulClient.Catalog().Service(serviceName, "", nil)
	if err != nil {
		return 0, "", err
	}
	service := services[0]
	return service.ServicePort, service.ServiceAddress, nil
}
