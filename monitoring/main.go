package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	return origin == "http://localhost:1337"
}}

var connections = []*websocket.Conn{}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	connections = append(connections, c)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	// defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
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

	//forever := make(chan bool)
	eventChan := make(chan string)
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
	go publishEventToWsConnections(eventChan)
	http.HandleFunc("/", echo)

	log.Fatal(http.ListenAndServe(":8080", nil))

	//<-forever
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

func processMessage(msg Message, eventChan chan string) {

	msgName, ok := msg["name"].(string)
	if !ok {
		log.Printf("msg.name is not a string")
		return
	}
	log.Printf("Received message: %s", msgName)
	eventChan <- msgName
}
func publishEventToWsConnections(eventChan chan string) {
	for {
		event := <-eventChan
		log.Printf("Publishing event: %s", event)
		for _, c := range connections {
			err := c.WriteMessage(websocket.TextMessage, []byte(event))
			if err != nil {
				log.Printf("Error writing message to websocket: %v", err)
			}
		}
	}
}
