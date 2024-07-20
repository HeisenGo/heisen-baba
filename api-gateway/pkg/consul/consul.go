package consul

import (
	"github.com/hashicorp/consul/api"
)

type Client struct {
	client *api.Client
}

func NewConsulClient(address string) (*Client, error) {
	config := api.DefaultConfig()
	config.Address = address

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Client{client: client}, nil
}

func (c *Client) GetService(serviceName string) ([]*api.ServiceEntry, error) {
	services, _, err := c.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}

	return services, nil
}
