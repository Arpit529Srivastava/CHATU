package main

import (
	//"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true }, // Allow all origins (adjust for security)
	}
)

type Manager struct {
	clients  ClientList
	handlers map[string]EventHandler
	sync.RWMutex
}

// Creates a new Manager instance
func NewManager() *Manager {
	m := &Manager{
		clients:  make(ClientList),
		handlers: make(map[string]EventHandler),
	}
	m.setupEventHandlers()
	return m
}

// Registers event handlers
func (m *Manager) setupEventHandlers() {
	m.handlers[EventSendMessage] = SendMessage
}

// Example handler for "send_message" events
func SendMessage(event Event, c *Client) error {
	fmt.Println(event)
	return nil
}

// Routes an event to the appropriate handler
func (m *Manager) routeEvent(event Event, c *Client) error {
	if handler, exists := m.handlers[event.Type]; exists {
		return handler(event, c)
	}
	log.Printf("No handler found for event type: %s", event.Type)
	return errors.New("unknown event type")
}

// Serves WebSocket connections
func (m *Manager) serveWS(w http.ResponseWriter, r *http.Request) {
	log.Println("New WebSocket connection established")
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	client := NewClient(conn, m)
	m.addClient(client)
	go client.readMessage()
	go client.writeMessage()
}
func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	m.clients[client] = true
	log.Printf("Client added. Total clients: %d", len(m.clients))
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	if _, exists := m.clients[client]; exists {
		client.connection.Close()
		delete(m.clients, client)
		log.Printf("Client removed. Total clients: %d", len(m.clients))
	}
}

// Broadcasts an event to all connected clients
// func (m *Manager) broadcast(event Event) {
// 	m.RLock()
// 	defer m.RUnlock()
// 	for client := range m.clients {
// 		select {
// 		case client.egress <- event:
// 		default:
// 			log.Println("Client egress channel full, dropping event")
// 		}
// 	}
// }
