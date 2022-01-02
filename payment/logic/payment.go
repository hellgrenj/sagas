package logic

import (
	"fmt"

	"github.com/hellgrenj/sagas/payment/models"
)

type Logger interface {
	Info(msg string)
	Error(msg string)
}
type DBAccess interface {
	InsertPayment(orderPayment models.OrderPayment) error
}
type PaymentHandler interface {
	ChargeCustomer(orderPayment models.OrderPayment) bool
}
type payment struct {
	db     DBAccess
	logger Logger
}

func NewPaymentHandler(db DBAccess, logger Logger) *payment {
	return &payment{db: db, logger: logger}
}
func (p *payment) ChargeCustomer(orderPayment models.OrderPayment) bool {
	p.logger.Info("charging customer")
	err := p.db.InsertPayment(orderPayment)
	if err != nil {
		p.logger.Error(fmt.Sprintf("failed to create payment %s", err))
		return false
	} else {
		p.logger.Info("payment created")
		return true
	}
}
