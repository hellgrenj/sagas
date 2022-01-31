package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func StartListen(eventChan chan string) {
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
	// // defer c.Close()
	// for {
	// 	mt, message, err := c.ReadMessage()
	// 	if err != nil {
	// 		log.Println("read:", err)
	// 		break
	// 	}
	// 	log.Printf("recv: %s", message)
	// 	err = c.WriteMessage(mt, message)
	// 	if err != nil {
	// 		log.Println("write:", err)
	// 		break
	// 	}
	// }
}
func publishEventToWsConnections(eventChan chan string) {
	for {
		event := <-eventChan
		log.Printf("Publishing event: %s", event)
		for index, c := range connections {
			log.Printf("Sending event to connection %d", index)
			err := c.WriteMessage(websocket.TextMessage, []byte(event))
			if err != nil {
				log.Printf("Error writing message to websocket: %v", err)
				// remove connection from list
				connections = append(connections[:index], connections[index+1:]...)
			}
		}
	}
}
