package ports

type IMessageBroker interface {
	Publish(queueName, msg string)
	Consume(queueName string, execute func(msg string))
}
