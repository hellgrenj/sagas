package main

import (
	"github.com/hellgrenj/sagas/monitoring/models"
	"github.com/hellgrenj/sagas/monitoring/rabbit"
	"github.com/hellgrenj/sagas/monitoring/ws"
)

func main() {
	// TODO this is a spike - go over error handling and logging and TODO's etc...
	// TODO fix connection management.. getting duplicates now.. ping-pong setup?
	eventChan := make(chan models.Event)
	go rabbit.StartListen(eventChan)
	ws.StartListen(eventChan)
}
