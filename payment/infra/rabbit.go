package infra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/hellgrenj/sagas/payment/events/inbound"
	"github.com/hellgrenj/sagas/payment/events/outbound"
	"github.com/hellgrenj/sagas/payment/logic"
	"github.com/hellgrenj/sagas/payment/models"

	"github.com/streadway/amqp"
)

type rabbit struct {
	PaymentHandler logic.PaymentHandler
	InfraHandler   InfraHandler
	logger         Logger
	conn           *amqp.Connection
}

func NewRabbitWorker(paymentHandler logic.PaymentHandler, infraHandler InfraHandler, logger Logger) *rabbit {
	return &rabbit{PaymentHandler: paymentHandler, InfraHandler: infraHandler, logger: logger, conn: nil}
}

type Message map[string]interface{}

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
		r.logger.Error(fmt.Sprintf("Unable to connect to rabbit: %v\n", err))
		if connectionAttempt < 5 {
			connectionAttempt++
			r.logger.Info(fmt.Sprintf("Trying again in 4 seconds attempt %v of 5\n", connectionAttempt))
			time.Sleep(4 * time.Second)
			return r.TryConnectToRabbit(connectionAttempt)
		}
		os.Exit(1)
	}
	r.logger.Info("Successfully connected to rabbit")
	return conn
}

func (r *rabbit) failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
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
			r.logger.Info(fmt.Sprintf(" [x] %s", d.Body))
			var msg, err = deserialize(d.Body)
			if err != nil {
				r.logger.Error(fmt.Sprintf("Error deserializing message: %v", err))
				continue
			}
			messageId, ok := msg["messageId"].(string)
			if !ok {
				r.logger.Error("messageId is missing or not a string")
				continue
			}
			alreadyProcessed, err := r.InfraHandler.TryMarkMessageAsProcessed(messageId)
			r.failOnError(err, "Error marking message as processed")
			if alreadyProcessed {
				r.logger.Info(fmt.Sprintf("Message %v already processed", msg["messageId"]))
			} else {
				go r.processMessage(msg)
			}
		}
	}()

	r.logger.Info(fmt.Sprintf(" [*] Listening on exchange %s. To exit press CTRL+C", exchange))
	<-forever
}

func (r *rabbit) processMessage(msg Message) {
	r.logger.Info("Processing message")
	messageName, ok := msg["name"].(string)
	if !ok {
		r.logger.Error("message name is missing or not a string")
		return
	}
	correlationId, correlationIdOk := msg["correlationId"].(string)
	if !correlationIdOk {
		r.logger.Error("correlationId is missing or not a string")
		return
	}
	if messageName == "items reserved" {
		itemsReservedEvent, err := inbound.MapToItemsReservedEvent(msg)
		if err != nil {
			r.logger.Error(fmt.Sprintf("Error mapping message to reservation: %v", err))
			return
		}
		orderPayment := models.OrderPayment{
			Amount:   itemsReservedEvent.Price,
			Item:     itemsReservedEvent.Item,
			Quantity: itemsReservedEvent.Quantity,
		}
		paymentSucceded := r.PaymentHandler.ChargeCustomer(orderPayment)
		if paymentSucceded {
			paymentSucceededEvent := outbound.PaymentEvent{
				CorrelationId: correlationId,
				Name:          "payment succeeded",
				MessageId:     uuid.New().String(),
				OrderId:       itemsReservedEvent.OrderId,
			}
			r.publishMessage(paymentSucceededEvent)
		}
	} else {
		r.logger.Info("Unknown message type")
	}
}

func (r *rabbit) publishMessage(paymentEvent outbound.PaymentEvent) {
	ch, err := r.conn.Channel()
	r.failOnError(err, "Failed to open a channel")
	defer ch.Close()
	body, err := json.Marshal(paymentEvent)
	if err != nil {
		r.logger.Error(err.Error())
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
	r.logger.Info(fmt.Sprintf(" [x] Sent %s", body))
}
