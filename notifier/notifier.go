package main

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

var rabbitMQConn *amqp.Connection

func setupRabbitMQ() {
	conn, err := amqp.Dial("amqp://user:123@rabbitmq.default.svc.cluster.local:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	rabbitMQConn = conn
}

func consumeFromRabbitMQ() {
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

	// Start a consumer
	msgs, err := ch.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Println("Received a message from: ", d.RoutingKey)
			handleMessage(d.Body)
		}

	}()

	<-forever // Block main thread to keep it running
}

type Item struct {
	ID    string  `json:"id"`
	Owner string  `json:"owner"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func handleMessage(message []byte) {
	var order Item
	err := json.Unmarshal(message, &order)
	if err != nil {
		log.Println("Failed to unmarshal message: ", err)
		return
	}
}

func main() {
	setupRabbitMQ()
	defer rabbitMQConn.Close()

	consumeFromRabbitMQ()
}
