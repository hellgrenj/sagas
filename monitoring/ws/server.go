package ws

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hellgrenj/sagas/monitoring/models"
)

func StartListen(eventChan chan models.Event) {
	hub := newHub()
	go hub.run()
	go publishEventsToWsConnections(eventChan, hub)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		connect(w, r, hub)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	return origin == "http://localhost:1337" || origin == "http://localhost"
}}

const (
	// Time allowed to read the next pong message from the client.
	pongWait = 10 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

func connect(w http.ResponseWriter, r *http.Request, hub *Hub) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	c.SetReadLimit(512)
	c.SetReadDeadline(time.Now().Add(pongWait))
	c.SetPongHandler(func(string) error { c.SetReadDeadline(time.Now().Add(pongWait)); return nil })
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
				hub.unregister <- c
				break
			}
		}
	}()
	hub.register <- c
}
func publishEventsToWsConnections(eventChan chan models.Event, hub *Hub) {
	for {
		event := <-eventChan
		hub.events <- event
	}
}
