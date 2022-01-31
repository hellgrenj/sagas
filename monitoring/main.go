package main

import (
	"github.com/hellgrenj/sagas/monitoring/rabbit"
	"github.com/hellgrenj/sagas/monitoring/ws"
)

func main() {
	eventChan := make(chan string)
	go rabbit.StartListen(eventChan)
	ws.StartListen(eventChan)
}
