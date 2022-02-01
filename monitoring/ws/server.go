package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/hellgrenj/sagas/monitoring/models"
)

func StartListen(eventChan chan models.Event) {
	go publishEventToWsConnections(eventChan)
	http.HandleFunc("/ws", connect)

	log.Fatal(http.ListenAndServe(":8080", nil)) // keeps process alive
}

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	return origin == "http://localhost:1337"
}}

var connections = []*websocket.Conn{}

func connect(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	connections = append(connections, c)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
}
func publishEventToWsConnections(eventChan chan models.Event) {
	for {
		event := <-eventChan
		log.Printf("Publishing event: %s with correlationId %v and messageId %v", event.Name, event.CorrelationId, event.MessageId)
		for index, c := range connections {
			log.Printf("Sending event to connection %d", index)
			err := c.WriteJSON(event)
			if err != nil {
				log.Printf("Error writing message to websocket: %v", err)
				// remove connection from list
				connections = append(connections[:index], connections[index+1:]...)
			}
		}
	}
}
