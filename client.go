package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// ping interval should be lower than the pong wait
// we are providing a 90% time frame to wait
var (
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	manager    *Manager
	egress     chan Event
}

// Constructor for creating a new Client instance
func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan Event),
	}
}

// Gracefully shuts down the client
func (c *Client) shutdown() {
	log.Println("Shutting down client")
	c.manager.removeClient(c)

	// Close the egress channel to signal write goroutine to stop
	close(c.egress)

	// Close the WebSocket connection
	if err := c.connection.Close(); err != nil {
		log.Printf("Error closing WebSocket connection: %v", err)
	}
}

// Handles incoming messages from the client
func (c *Client) readMessage() {
	defer c.shutdown()

	if err := c.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println("Error setting read deadline:", err)
		return
	}

	c.connection.SetPongHandler(func(appData string) error {
		log.Println("Received pong")
		return c.connection.SetReadDeadline(time.Now().Add(pongWait))
	})
	c.connection.SetReadLimit(512) // setting the limit for sending the message

	for {
		_, payload, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected close error: %v", err)
			} else {
				log.Println("Connection closed by client")
			}
			break
		}

		var request Event
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("Error unmarshalling event: %v", err)
			continue
		}

		if err := c.manager.routeEvent(request, c); err != nil {
			log.Printf("Error handling event: %v", err)
		}
	}
}

// Handles outgoing messages to the client
func (c *Client) writeMessage() {
	defer c.shutdown()

	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				// Channel closed, send close message
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Printf("Error sending close message: %v", err)
				}
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Printf("Error marshalling message: %v", err)
				continue
			}

			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("Error sending message: %v", err)
				return
			}
			log.Println("Message sent successfully")

		case <-ticker.C:
			log.Println("Sending ping")
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Printf("Error sending ping: %v", err)
				return
			}
		}
	}
}
