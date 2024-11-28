// it will handle, whenever we are connecting to a new room it will create a client for in the backend
package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type ClientList map[*Client]bool
type Client struct {
	connection *websocket.Conn
	manager    *Manager
	// we are using egress for avoiding concurrencies write on the connection
	egress chan []byte
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress: make(chan []byte),
	}
}

// reading and writing the messages
// websockets only allow one at a time not simultaneously
func (c *Client) readMessage() {
	defer func() {
		//cleaning up the connection
		c.manager.removeClient(c)
	}()
	for {
		messageType, payload, err := c.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error while reading message : %v", err)
			}
			break
		}
		// just a hack ;)
		for wsclient := range c.manager.clients{
			wsclient.egress <- payload
		}
		log.Println(messageType)
		log.Println(string(payload))
	}
}

// starting with concurrency
func (c *Client) writeMessage(){
	defer func(){
		c.manager.removeClient((c))
	}()
	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil{
					log.Println("connection id closed :", err)
				} 
				return
			}
			if err := c.connection.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("failed to send message : %v", err)
			}
		}
	}
}
