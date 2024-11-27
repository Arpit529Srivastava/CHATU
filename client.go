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
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
	}
}
// reading and writing the messages
// websockets only allow one at a time not simultaneously
func (c *Client) readMessage(){
	defer func(){
		//cleaning up the connection
		c.manager.removeClient(c)
	}()
	for {
		messageType, payload, err :=c.connection.ReadMessage()

		if err != nil{
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure){
				log.Printf("error while reading message : %v", err)
			}
			break
		}
		log.Println(messageType)
		log.Println(string(payload))
	}
}
// starting with concurrency