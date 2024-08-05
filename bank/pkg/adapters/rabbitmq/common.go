package rabbitmq

type RabbitMQ struct {
	msgBrokerUsername     string
	msgBrokerHostPassword string
	msgBrokerHost         string
	msgBrokerPort         int
}

func NewRabbitMQ(msgBrokerUsername string, msgBrokerHostPassword string, msgBrokerHost string, msgBrokerPort int) *RabbitMQ {
	return &RabbitMQ{msgBrokerUsername: msgBrokerUsername, msgBrokerHostPassword: msgBrokerHostPassword, msgBrokerHost: msgBrokerHost, msgBrokerPort: msgBrokerPort}
}
