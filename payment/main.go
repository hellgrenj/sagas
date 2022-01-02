package main

import (
	"github.com/hellgrenj/sagas/payment/db"
	"github.com/hellgrenj/sagas/payment/infra"
	"github.com/hellgrenj/sagas/payment/logic"
)

func main() {
	db := db.NewDBAccess(infra.NewLogger("db"))
	paymentHandler := logic.NewPaymentHandler(db, infra.NewLogger("payment"))
	infraHandler := infra.NewInfraHandler(db, infra.NewLogger("infra"))
	worker := infra.NewRabbitWorker(paymentHandler, infraHandler, infra.NewLogger("rabbit"))
	worker.StartListen()
}
