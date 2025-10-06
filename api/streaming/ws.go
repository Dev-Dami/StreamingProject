package streaming

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Client represents a WebSocket client
type Client struct {
	conn  *websocket.Conn
	mutex sync.Mutex
	id    int
}

// WriteMessage safely writes a message to the WebSocket
func (c *Client) WriteMessage(messageType int, data []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.conn.WriteMessage(messageType, data)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var (
	clients      = make(map[*Client]bool)
	clientsMutex sync.RWMutex
	broadcast    = make(chan []byte, 200) // Larger buffer
	clientIDGen  = 0
)

// Start broadcaster when package loads
func init() {
	go startBroadcaster()
}

// ServeWS handles WebSocket connections
func ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	clientsMutex.Lock()
	clientIDGen++
	client := &Client{
		conn: conn,
		id:   clientIDGen,
	}
	clients[client] = true
	log.Printf("Client %d connected. Total clients: %d", client.id, len(clients))
	clientsMutex.Unlock()

	// Handle client cleanup
	defer func() {
		conn.Close()
		clientsMutex.Lock()
		delete(clients, client)
		log.Printf("Client %d disconnected. Total clients: %d", client.id, len(clients))
		clientsMutex.Unlock()
	}()

	// Keep connection alive with ping/pong
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Read messages to detect client disconnection
	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return // This will trigger the defer cleanup
			}
		}
	}()

	// Send periodic pings
	for {
		select {
		case <-ticker.C:
			if err := client.WriteMessage(websocket.PingMessage, nil); err != nil {
				return // Connection failed
			}
		}
	}
}

// Broadcast sends frame to all connected clients
func Broadcast(frame []byte) {
	select {
	case broadcast <- frame:
		// Frame queued successfully
	default:
		// Channel full, drop frame (this is okay for video streaming)
	}
}

// startBroadcaster runs the broadcasting loop
func startBroadcaster() {
	log.Println("Broadcaster started")
	frameCount := 0

	for frame := range broadcast {
		frameCount++
		
		clientsMutex.RLock()
		if len(clients) == 0 {
			clientsMutex.RUnlock()
			continue // No clients, skip frame
		}
		
		// Send frame to all clients
		var failedClients []*Client
		for client := range clients {
			if err := client.WriteMessage(websocket.BinaryMessage, frame); err != nil {
				failedClients = append(failedClients, client)
			}
		}
		clientCount := len(clients)
		clientsMutex.RUnlock()

		// Remove failed clients
		if len(failedClients) > 0 {
			clientsMutex.Lock()
			for _, client := range failedClients {
				delete(clients, client)
				client.conn.Close()
				log.Printf("Removed failed client %d", client.id)
			}
			clientsMutex.Unlock()
		}

		// Log progress
		if frameCount%100 == 0 {
			log.Printf("Broadcast frame %d to %d clients (size: %d bytes)", frameCount, clientCount, len(frame))
		}
	}
}
