package main

import (
	"insider/graceful_shutdown"
	"insider/server"
)

func main() {
	server.Configure()
	graceful_shutdown.KeepAppUp()
}
