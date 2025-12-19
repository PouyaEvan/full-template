package http

import (
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

// WebSocketHandler handles websocket connections
type WebSocketHandler struct {
	clients    map[*websocket.Conn]bool
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	broadcast  chan []byte
	mu         sync.Mutex
}

// NewWebSocketHandler creates a new WebSocketHandler
func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{
		clients:    make(map[*websocket.Conn]bool),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		broadcast:  make(chan []byte),
	}
}

// Run starts the hub loop
func (h *WebSocketHandler) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Println("Client connected")

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
			}
			h.mu.Unlock()
			log.Println("Client disconnected")

		case message := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
					log.Println("write error:", err)
					client.Close()
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

// HandleWebSocket handles the websocket connection
func (h *WebSocketHandler) HandleWebSocket(c *websocket.Conn) {
	defer func() {
		h.unregister <- c
		c.Close()
	}()

	h.register <- c

	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("read error:", err)
			}
			return
		}
		// For now, we just ignore incoming messages or echo them if needed.
		// But the requirement is mostly for notifications (server -> client).
	}
}

// BroadcastMessage sends a message to all connected clients
func (h *WebSocketHandler) BroadcastMessage(message string) {
	h.broadcast <- []byte(message)
}

// WebSocketMiddleware to upgrade connection
func WebSocketMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return c.Status(fiber.StatusUpgradeRequired).SendString("Upgrade Required")
	}
}
