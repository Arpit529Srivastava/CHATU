package main

import (
	"context"
	"log"
	"net/http"
)

func main() {
	setupAPI()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func setupAPI() {
	ctx := context.Background()
	manager := NewManager(ctx)
	http.HandleFunc("/ws", manager.serveWS) // the manager will take the request and add it to the connection of the websocket
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/login", manager.loginHandler)
}
