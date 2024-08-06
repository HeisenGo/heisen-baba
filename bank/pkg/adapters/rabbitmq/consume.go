package rabbitmq

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnConsumeError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}

func (r *RabbitMQ) Consume(queueName string, execute func(msg string)) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%v/", r.msgBrokerUsername, r.msgBrokerPassword, r.msgBrokerHost, r.msgBrokerPort))
	failOnConsumeError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnConsumeError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName, // name of the queue
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnConsumeError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnConsumeError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			userID := string(d.Body)
			execute(userID)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
