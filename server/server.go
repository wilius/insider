package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	"insider/configs"
	"insider/graceful_shutdown"
	"insider/message"
	"insider/sender"
	"net/http"
	"time"
)

var httpServer *http.Server

func Configure() {
	serverConfig := configs.
		Instance().
		GetServer()

	port := fmt.Sprintf(":%d", serverConfig.GetPort())

	log.Info().
		Msgf("starting server at localhost:%s", port)

	httpServer = &http.Server{
		Addr:    port,
		Handler: buildRouters(),
	}

	graceful_shutdown.AddShutdownHook(func() {
		log.Info().Msg("Closing http listener")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := httpServer.Shutdown(ctx)
		if err != nil {
			log.Err(err).
				Msgf("http server shutdown failed %v", err)
		}
	})

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Info().
					Msg("Listener shutdown has been gracefully completed")
			} else {
				log.Fatal().
					Msgf("server start failed: %v", err)
			}
		}
	}()
}

func buildRouters() *chi.Mux {
	r := chi.NewRouter()
	r.Use(chiMiddleware.Heartbeat("/health"))

	message.Configure(r)
	sender.Configure(r)

	return r
}
