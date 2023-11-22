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

	// Just for learning purposes, let's create a custom topic exchange and bind it to the queue
	ch, err := rabbitMQConn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare("ecom-exchange", "topic", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare an exchange: %v", err)
	}

	queue, err := ch.QueueDeclare("ecom-queue", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %v", err)
	}

	// We're gonna define a queue with its binding being ecom.* which means it will receive all messages
	// if the routing key starts with ecom.
	err = ch.QueueBind(queue.Name, "ecom.*", "ecom-exchange", false, nil)
	if err != nil {
		return fmt.Errorf("failed to bind the queue to the exchange: %v", err)
	}

	return nil
}
func publishToRabbitMQ(routingKey string, message []byte) {
	ch, err := rabbitMQConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	err = ch.Publish(
		"ecom-exchange", // exchange name
		routingKey,      // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}
}
