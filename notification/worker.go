package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	conn := TryConnectToRabbit(1)
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	exchange := "order.topics"
	err = ch.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		false,    // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"notification_queue", // name
		false,                // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,         // queue name
		"#",            // routing key = all topics
		"order.topics", // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var msg, err = deserialize(d.Body)
			if err != nil {
				log.Printf("Error deserializing message: %v", err)
				continue
			}
			go processMessage(msg)
		}
	}()

	log.Printf(" [*] Listening on exchange %s ALL topics. To exit press CTRL+C", exchange)
	<-forever
}

type Message map[string]interface{}

func deserialize(b []byte) (Message, error) {
	var msg Message
	buf := bytes.NewBuffer(b)
	decoder := json.NewDecoder(buf)
	err := decoder.Decode(&msg)
	return msg, err
}
func TryConnectToRabbit(connectionAttempt int) *amqp.Connection {
	conn, err := amqp.Dial("amqp://guest:guest@rabbit:5672/")
	if err != nil {
		log.Printf("Unable to connect to rabbit: %v\n", err)
		if connectionAttempt < 5 {
			connectionAttempt++
			log.Printf("Trying again in 4 seconds attempt %v of 5\n", connectionAttempt)
			time.Sleep(4 * time.Second)
			return TryConnectToRabbit(connectionAttempt)
		}
		os.Exit(1)
	}
	log.Println("Successfully connected to rabbit")
	return conn
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func processMessage(msg Message) {

	var customerMessage string
	switch msg["name"].(string) {
	case "order created":
		order := msg["order"].(map[string]interface{})
		customerMessage = fmt.Sprintf("Your order (#%v) has been created", order["id"].(float64))
	case "order cancelled":
		order := msg["order"].(map[string]interface{})
		reason := msg["reason"].(string)
		customerMessage = fmt.Sprintf("Your order (#%v) has been cancelled. Reason: %s", order["id"].(float64), reason)
	// case "items reserved": not communicated to customer
	// case "items not in stock": not communicated to customer
	case "payment succeeded":
		orderId := msg["orderId"].(float64)
		customerMessage = fmt.Sprintf("Your payment has been completed for order #%v", orderId)
	case "order shipped":
		orderId := msg["orderId"].(float64)
		customerMessage = fmt.Sprintf("Your order (#%v) has been shipped", orderId)
	// case "order completed": not communicated to customer
	default:
	}
	if customerMessage != "" {
		log.Printf("Sent the following to the customer:\n%s", customerMessage)
	}
}
