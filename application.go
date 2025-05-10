package main

import (
	"insider/database"
	"insider/graceful_shutdown"
	"insider/sender"
	"insider/server"
)

func main() {
	database.Configure()
	server.Configure()
	sender.Start()
	graceful_shutdown.KeepAppUp()
}
