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

func (c *Consul) RegisterService(serviceHostName, servicePrefixPath, serviceHTTPHealthPath string, serviceHTTPPort int) error {
	consulConfig := consulAPI.DefaultConfig()
	consulConfig.Address = c.Address
	consulClient, err := consulAPI.NewClient(consulConfig)
	if err != nil {
		return err
	}

	HTTPHealthURL := fmt.Sprintf("http://%s:%v/%s", serviceHostName, serviceHTTPPort, serviceHTTPHealthPath)
	// Register service with Consul
	registration := &consulAPI.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-service-id", serviceHostName),
		Name:    serviceHostName,
		Address: serviceHostName,
		Port:    serviceHTTPPort,
		Tags: []string{
			serviceHostName,
			fmt.Sprintf("traefik.http.routers.%s_router.rule=PathPrefix(`%s`)", serviceHostName, servicePrefixPath),
			fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.port=%v", serviceHostName, serviceHTTPPort),
			fmt.Sprintf("traefik.http.middlewares.auth-middleware.forwardauth.address=http://oauth2-proxy:4180/oauth2/auth"),
			fmt.Sprintf("traefik.http.middlewares.auth-middleware.forwardauth.trustForwardHeader=true"),
			fmt.Sprintf("traefik.http.middlewares.auth-middleware.forwardauth.authResponseHeaders=X-Auth-User, X-Secret"),
			fmt.Sprintf("traefik.http.middlewares.auth-middleware.forwardauth.authRequestHeaders=Accept,X-CustomHeader"),
			fmt.Sprintf("traefik.http.middlewares.auth-middleware.forwardauth.addAuthCookiesToResponse=Session-Cookie,State-Cookie"),
			fmt.Sprintf("traefik.http.routers.%s_router.middlewares=auth-middleware", serviceHostName),
		},
		Check: &consulAPI.AgentServiceCheck{
			HTTP:     HTTPHealthURL,
			Interval: "10s",
			Timeout:  "1s",
		},
	}

	err = consulClient.Agent().ServiceRegister(registration)
	if err != nil {
		return err
	}
	return nil
}
