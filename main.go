package main

import (
	"github.com/shevchik87/chat-goland/socket"
)

func main() {
	hub := socket.NewHub()
	go hub.Run()

	a := App{}
	a.Initialize("postgres", "postgres", "chat")

	a.Run(":8080")
}
