package main

func main() {
	logger := NewLogger()
	db := NewDBAccess(logger)
	paymentHandler := NewPayment(db, logger)
	infraHandler := NewInfraHandler(db, logger)
	worker := NewRabbitWorker(paymentHandler, infraHandler, logger)
	worker.StartListen()
}
