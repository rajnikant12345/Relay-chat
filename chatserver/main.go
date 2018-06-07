package main

import (
	"Relay-chat/chatserver/cmd"
	"log"
)

func main() {

	log.Println("Starting server")
	cmd.StartServer()
}
