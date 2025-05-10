package main

import (
	"insider/graceful_shutdown"
	"insider/sender"
	"insider/server"
)

func main() {
	server.Configure()
	sender.Start()
	graceful_shutdown.KeepAppUp()
}
