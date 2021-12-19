package main

import (
	"context"
	"fmt"
)

type payment struct {
	db     *DBAccess
	logger Logger
}

type OrderPayment struct {
	Amount   float64 `bson:"amount"`
	Item     string  `bson:"item"`
	Quantity float64 `bson:"quantity"`
}

func NewPayment(db *DBAccess, logger Logger) *payment {
	return &payment{db: db, logger: logger}
}
func (p *payment) ChargeCustomer(orderPayment OrderPayment) bool {
	p.logger.Info("charging customer")
	_, err := p.db.conn.Database("payment").Collection("payments").InsertOne(context.TODO(), orderPayment)
	if err != nil {
		p.logger.Error(fmt.Sprintf("failed to create payment %s", err))
		return false
	} else {
		p.logger.Info("payment created")
		return true
	}
}
