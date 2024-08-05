package ports

type IMessageBroker interface {
	Publish(queueName, msg string)
}
