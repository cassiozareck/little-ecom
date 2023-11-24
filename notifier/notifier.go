package main

import (
	"encoding/json"
	"fmt"
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

func consumeFromRabbitMQ(queueName string, routingKey string, handler func([]byte)) {
	ch, err := rabbitMQConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	// Declare the queue
	queue, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Bind the queue to the exchange
	err = ch.QueueBind(
		queue.Name,      // queue name
		routingKey,      // routing key
		"ecom-exchange", // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to bind the queue: %v", err)
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

	go func() {
		for d := range msgs {
			log.Println("Received a message from: ", d.RoutingKey)
			handler(d.Body)
		}
	}()
}

type NotificationItem struct {
	Email string  `json:"email"`
	Name  string  `json:"item"`
	Price float64 `json:"price"`
}

func handlerAddedItems(message []byte) {
	notificationItem := NotificationItem{}
	err := json.Unmarshal(message, &notificationItem)
	if err != nil {
		log.Println("Failed to unmarshal message: ", err)
		return
	}

	// Email content.
	subject := "New Item Added: " + notificationItem.Name
	body := "A new item has been added with a price of $" + fmt.Sprintf("%.2f", notificationItem.Price)

	// Since emails services like gmail can block emails sent repeatedly from the same account, I'll comment
	// the lines and just print the email content to the console. But keep in mind that the code is correctly
	// implemented if you uncomment and put your own email credentials under notifier/smtp.go

	// Send the email.
	//err = sendEmail(notificationItem.Email, subject, body)
	//if err != nil {
	//	log.Printf("Failed to send email: %v\n", err)
	//}

	log.Printf("Email sent successfully \n%v \n%v\n", subject, body)
}

func handlerBoughtItems(message []byte) {
	notificationItem := NotificationItem{}
	err := json.Unmarshal(message, &notificationItem)
	if err != nil {
		log.Println("Failed to unmarshal message: ", err)
		return
	}

	// Email content.
	subject := "Item Bought: " + notificationItem.Name
	body := "You bought an item with a price of $" + fmt.Sprintf("%.2f", notificationItem.Price)

	log.Printf("Email sent successfully \n%v \n%v\n", subject, body)
}

func main() {
	setupRabbitMQ()
	defer rabbitMQConn.Close()

	consumeFromRabbitMQ("ecom-queue-item-added", "ecom.item.add", handlerAddedItems)
	consumeFromRabbitMQ("ecom-queue-item-bought", "ecom.item.buy", handlerBoughtItems)

	select {} // Block
}
