package main

import (
	"bytes"
	"encoding/json"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

type PaymentHandler interface {
	ChargeCustomer(orderPayment OrderPayment) bool
}
type InfraHandler interface {
	TryMarkMessageAsProcessed(messageId string) (bool, error)
}

type rabbit struct {
	PaymentHandler PaymentHandler
	InfraHandler   InfraHandler
	logger         *Logger
	conn           *amqp.Connection
}

func NewRabbitWorker(paymentHandler PaymentHandler, infraHandler InfraHandler, logger *Logger) *rabbit {
	return &rabbit{PaymentHandler: paymentHandler, InfraHandler: infraHandler, logger: logger, conn: nil}
}

type Message map[string]interface{}
type PaymentEvent struct {
	CorrelationId string  `json:"correlationId"`
	Name          string  `json:"name"`
	MessageId     string  `json:"messageId"`
	OrderId       float64 `json:"orderId"`
}

func deserialize(b []byte) (Message, error) {
	var msg Message
	buf := bytes.NewBuffer(b)
	decoder := json.NewDecoder(buf)
	err := decoder.Decode(&msg)
	return msg, err
}
func (r *rabbit) TryConnectToRabbit(connectionAttempt int) *amqp.Connection {
	conn, err := amqp.Dial("amqp://guest:guest@rabbit:5672/")
	if err != nil {
		r.logger.error.Printf("Unable to connect to rabbit: %v\n", err)
		if connectionAttempt < 5 {
			connectionAttempt++
			r.logger.info.Printf("Trying again in 4 seconds attempt %v of 5\n", connectionAttempt)
			time.Sleep(4 * time.Second)
			return r.TryConnectToRabbit(connectionAttempt)
		}
		os.Exit(1)
	}
	r.logger.info.Println("Successfully connected to rabbit")
	return conn
}

func (r *rabbit) failOnError(err error, msg string) {
	if err != nil {
		r.logger.error.Fatalf("%s: %s", msg, err)
	}
}
func (r *rabbit) StartListen() {
	r.conn = r.TryConnectToRabbit(1)
	defer r.conn.Close()

	ch, err := r.conn.Channel()
	r.failOnError(err, "Failed to open a channel")
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
	r.failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"payment_queue", // name
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	r.failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,                 // queue name
		"order.items.reserved", // routing key
		"order.topics",         // exchange
		false,
		nil)
	r.failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	r.failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			r.logger.info.Printf(" [x] %s", d.Body)
			var msg, err = deserialize(d.Body)
			if err != nil {
				r.logger.error.Printf("Error deserializing message: %v", err)
				continue
			}
			alreadyProcessed, err := r.InfraHandler.TryMarkMessageAsProcessed(msg["messageId"].(string))
			r.failOnError(err, "Error marking message as processed")
			if alreadyProcessed {
				r.logger.info.Printf("Message %v already processed", msg["messageId"])
			} else {
				go r.processMessage(msg)
			}
		}
	}()

	r.logger.info.Printf(" [*] Listening on exchange %s. To exit press CTRL+C", exchange)
	<-forever
}
func (r *rabbit) processMessage(msg Message) {
	r.logger.info.Println("Processing message")
	if msg["name"].(string) == "items reserved" {
		reservation := msg["reservation"].(map[string]interface{})
		orderPayment := OrderPayment{
			Amount:   reservation["price"].(float64),
			Item:     reservation["item"].(string),
			Quantity: reservation["quantity"].(float64),
		}
		paymentSucceded := r.PaymentHandler.ChargeCustomer(orderPayment)
		if paymentSucceded {
			paymentSucceededEvent := PaymentEvent{
				CorrelationId: msg["correlationId"].(string),
				Name:          "payment succeeded",
				MessageId:     uuid.New().String(),
				OrderId:       reservation["orderId"].(float64),
			}
			r.publishMessage(paymentSucceededEvent)
		}
	} else {
		r.logger.info.Println("Unknown message type")
	}
}
func (r *rabbit) publishMessage(paymentEvent PaymentEvent) {
	ch, err := r.conn.Channel()
	r.failOnError(err, "Failed to open a channel")
	defer ch.Close()
	body, err := json.Marshal(paymentEvent)
	if err != nil {
		r.logger.error.Println(err)
		return
	}
	err = ch.Publish(
		"order.topics",            // exchange
		"order.payment.completed", // routing key
		false,                     // mandatory
		false,                     // immediate
		amqp.Publishing{
			Body: []byte(body),
		})
	r.failOnError(err, "Failed to publish a message")
	r.logger.info.Printf(" [x] Sent %s", body)
}
