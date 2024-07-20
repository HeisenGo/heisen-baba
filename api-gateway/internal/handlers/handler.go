package handlers

import (
	"api-gateway/pkg/consul"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io/ioutil"
	"net/http"
)

type Handler struct {
	ConsulClient *consul.Client
}

func NewHandler(Client *consul.Client) *Handler {
	return &Handler{ConsulClient: Client}
}

func (h *Handler) ProxyRequest(c *fiber.Ctx) error {
	serviceName := c.Params("service")
	//servicePath := c.Params("*")

	services, err := h.ConsulClient.GetService(serviceName)
	if err != nil || len(services) == 0 {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "service not available",
		})
	}

	service := services[0].Service
	targetURL := fmt.Sprintf("http://%v:%v/%v", service.Address, service.Port, serviceName)

	// maybe gRPC
	resp, err := http.Get(targetURL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}
	defer resp.Body.Close()

	c.Set("Content-Type", resp.Header.Get("Content-Type"))
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	return c.JSON(map[string]string{"body": string(body)})
}
