package main

import (
	"insider/database"
	"insider/sender"
	"insider/server"
)

func main() {
	database.Configure()
	server.Configure()
	sender.Start()
	select {}
}
