package main

import "encoding/json"

// Event represents a generic event with a type and payload
type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"` // RawMessage allows flexible handling of payload data
}

// EventHandler defines a function type for handling events
type EventHandler func(event Event, c *Client) error

// Event types
const (
	EventSendMessage = "send_message"
)

// SendMessageEvent represents the structure for a "send message" event
type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
}
