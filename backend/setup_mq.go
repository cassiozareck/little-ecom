package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

var rabbitMQConn *amqp.Connection

func setupRabbitMQ() error {
	conn, err := amqp.Dial("amqp://user:123@rabbitmq.default.svc.cluster.local:5672/")
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	rabbitMQConn = conn
	return nil
}

func publishToRabbitMQ(message []byte) {
	ch, err := rabbitMQConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		"ecom-queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	err = ch.Publish(
		"",         // exchange
		queue.Name, // routing key (queue name)
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}
}
