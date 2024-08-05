package rabbitmq

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}

type RabbitMQ struct {
	msgBrokerUsername     string
	msgBrokerHostPassword string
	msgBrokerHost         string
	msgBrokerPort         int
}

func NewRabbitMQ(msgBrokerUsername string, msgBrokerHostPassword string, msgBrokerHost string, msgBrokerPort int) *RabbitMQ {
	return &RabbitMQ{msgBrokerUsername: msgBrokerUsername, msgBrokerHostPassword: msgBrokerHostPassword, msgBrokerHost: msgBrokerHost, msgBrokerPort: msgBrokerPort}
}

func (r *RabbitMQ) Publish(queueName, msg string) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%v/", r.msgBrokerUsername, r.msgBrokerHostPassword, r.msgBrokerHost, r.msgBrokerPort))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName, // name of the queue
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", queueName)
}
