package graceful_shutdown

import (
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

var shutdownHooks []func()

func init() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signalChan
		log.Info().Msgf("Received signal: %s. Shutting down gracefully...", sig)
		for _, hook := range shutdownHooks {
			hook()
		}
	}()

}

func AddShutdownHook(fn func()) {
	shutdownHooks = append(shutdownHooks, fn)
}
