package ws

import (
	"log"
	"net/http"
	"time"

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

const (
	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var connections = map[*websocket.Conn]bool{}

func connect(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	c.SetReadLimit(512)
	c.SetReadDeadline(time.Now().Add(pongWait))
	c.SetPongHandler(func(string) error { c.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	connections[c] = true
	go func() { // write pings to client..
		pingTicker := time.NewTicker(pingPeriod)
		for range pingTicker.C {
			log.Println("pinging client..")
			c.SetWriteDeadline(time.Now().Add(30 * time.Second))
			if err := c.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}()
	go func() { // read pongs from client
		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				log.Printf("removing connection %v. Pong failed with error %v", c.RemoteAddr(), err)
				delete(connections, c)
				c.Close()
				break
			}
		}
	}()

}
func publishEventToWsConnections(eventChan chan models.Event) {
	for {
		event := <-eventChan
		log.Printf("Publishing event: %s with correlationId %v and messageId %v", event.Name, event.CorrelationId, event.MessageId)
		for c := range connections {
			log.Printf("Sending event to connection %v", c.RemoteAddr())
			err := c.WriteJSON(event)
			if err != nil {
				log.Printf("Error writing message to websocket: %v", err)
			}
		}
	}
}
