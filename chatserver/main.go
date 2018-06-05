package main

import (
	"cryptolessons/chatserver/cmd"
	"log"
)

func main() {
	log.Println("Starting server")
	cmd.StartServer()
}
