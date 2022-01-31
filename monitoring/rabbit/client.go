package rabbit

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
)

func StartListen(eventChan chan string) {
	conn := tryConnectToRabbit(1)
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
		"monitoring_queue", // name
		false,              // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
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
			go processMessage(msg, eventChan)
		}
	}()
	log.Printf(" [*] Listening on exchange %s ALL topics. To exit press CTRL+C", exchange)

	<-forever
}

type message map[string]interface{}

func deserialize(b []byte) (message, error) {
	var msg message
	buf := bytes.NewBuffer(b)
	decoder := json.NewDecoder(buf)
	err := decoder.Decode(&msg)
	return msg, err
}
func tryConnectToRabbit(connectionAttempt int) *amqp.Connection {
	conn, err := amqp.Dial("amqp://guest:guest@rabbit:5672/")
	if err != nil {
		log.Printf("Unable to connect to rabbit: %v\n", err)
		if connectionAttempt < 5 {
			connectionAttempt++
			log.Printf("Trying again in 4 seconds attempt %v of 5\n", connectionAttempt)
			time.Sleep(4 * time.Second)
			return tryConnectToRabbit(connectionAttempt)
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

func processMessage(msg message, eventChan chan string) {

	msgName, ok := msg["name"].(string)
	if !ok {
		log.Printf("msg.name is not a string")
		return
	}
	log.Printf("Received message: %s", msgName)
	eventChan <- msgName
}
