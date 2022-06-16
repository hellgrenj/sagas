package ws

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/hellgrenj/sagas/monitoring/models"
)

type Hub struct {
	connections map[*websocket.Conn]bool
	register    chan *websocket.Conn
	unregister  chan *websocket.Conn
	events      chan models.Event
}

func newHub() *Hub {
	return &Hub{
		events:      make(chan models.Event),
		register:    make(chan *websocket.Conn),
		unregister:  make(chan *websocket.Conn),
		connections: make(map[*websocket.Conn]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case connection := <-h.register:
			h.connections[connection] = true
			log.Printf("number of connections %v", len(h.connections))
		case connection := <-h.unregister:
			delete(h.connections, connection)
			connection.Close()
			log.Println("removed connection")
			log.Printf("number of connections %v", len(h.connections))
		case event := <-h.events:
			log.Printf("Publishing event: %s with correlationId %v and messageId %v", event.Name, event.CorrelationId, event.MessageId)
			for c := range h.connections {
				log.Printf("Sending event to connection %v", c.RemoteAddr())
				err := c.WriteJSON(event)
				if err != nil {
					log.Printf("Error writing message to websocket: %v", err)
				}
			}
		}
	}
}
