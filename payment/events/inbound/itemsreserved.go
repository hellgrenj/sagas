package inbound

import (
	"errors"
)

type ItemsReservedEvent struct {
	Price    float64
	Item     string
	Quantity float64
	OrderId  float64
}

func MapToItemsReservedEvent(msg map[string]interface{}) (ItemsReservedEvent, error) {
	var r ItemsReservedEvent
	reservation, ok := msg["reservation"].(map[string]interface{})
	if !ok {
		return r, errors.New("reservation is missing or not a json object")

	}
	price, priceOk := reservation["price"].(float64)
	if !priceOk {
		return r, errors.New("reservation price is missing or not a float")
	}
	item, itemOk := reservation["item"].(string)
	if !itemOk {
		return r, errors.New("reservation item s missing or not a string")
	}
	quantity, quantityOk := reservation["quantity"].(float64)
	if !quantityOk {
		return r, errors.New("reservation quantity is missing or not a float")
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
