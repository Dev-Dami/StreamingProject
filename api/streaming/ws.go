package streaming

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Upgrades HTTP connections to WebSocket connections.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// connected clients.
var clients = make(map[*websocket.Conn]bool)

// broadcast channel.
var broadcast = make(chan []byte, 30)

func init() {
	go broadcaster()
}

// Handle WebSocket connections.
func ServeWS(w http.ResponseWriter, r *http.Request) {
	go broadcaster()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("WebSocket upgrade:", err)
		return
	}
	defer conn.Close()

	clients[conn] = true
	defer delete(clients, conn)

	// Keep connection alive
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Broadcast frame to all clients.
func Broadcast(frame []byte) {
	select {
	case broadcast <- frame:
	default:
	}
}

// broadcaster.
func broadcaster() {
	for frame := range broadcast {
		for conn := range clients {
			if err := conn.WriteMessage(websocket.BinaryMessage, frame); err != nil {
				delete(clients, conn)
				conn.Close()
			}
		}
	}
}
