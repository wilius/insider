package graceful_shutdown

import (
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

var shutdownHooks []func()
var appUpChannel = make(chan int, 1)

func init() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signalChan
		log.Info().Msgf("Received signal: %s. Shutting down gracefully...", sig)
		for _, hook := range shutdownHooks {
			hook()
		}

		appUpChannel <- 1
	}()

}

func AddShutdownHook(fn func()) {
	shutdownHooks = append(shutdownHooks, fn)
}

func KeepAppUp() {
	select {
	case <-appUpChannel:
	}
}
