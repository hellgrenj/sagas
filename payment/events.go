package main

import "errors"

// outbound events
type PaymentEvent struct {
	CorrelationId string  `json:"correlationId"`
	Name          string  `json:"name"`
	MessageId     string  `json:"messageId"`
	OrderId       float64 `json:"orderId"`
}

// inbound events
type ItemsReservedEvent struct {
	Price    float64
	Item     string
	Quantity float64
	OrderId  float64
}

func MapToItemsReservedEvent(msg Message) (ItemsReservedEvent, error) {
	var r ItemsReservedEvent
	reservation, ok := msg["reservation"].(map[string]interface{})
	if !ok {
		return r, errors.New("reservation is missing or not a json object")

	}
	price, priceOk := reservation["price"].(float64)
	if !priceOk {
		return r, errors.New("reservation price is not a float")
	}
	item, itemOk := reservation["item"].(string)
	if !itemOk {
		return r, errors.New("reservation item is not a string")
	}
	quantity, quantityOk := reservation["quantity"].(float64)
	if !quantityOk {
		return r, errors.New("reservation quantity is not a float")
	}
	orderId, orderIdOk := reservation["orderId"].(float64)
	if !orderIdOk {
		return r, errors.New("orderId is missing or not a float")
	}
	r = ItemsReservedEvent{
		Price:    price,
		Item:     item,
		Quantity: quantity,
		OrderId:  orderId,
	}
	return r, nil
}
