package main

func main() {

	db := NewDBAccess(NewLogger("db"))
	paymentHandler := NewPaymentHandler(db, NewLogger("payment"))
	infraHandler := NewInfraHandler(db, NewLogger("infra"))
	worker := NewRabbitWorker(paymentHandler, infraHandler, NewLogger("rabbit"))
	worker.StartListen()
}
