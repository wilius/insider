package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	"insider/configs"
	"insider/database"
	"insider/graceful_shutdown"
	"insider/message"
	"insider/message_provider"
	"net/http"
)

var httpServer *http.Server

func StartServer() {
	database.Configure()

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
		_ = httpServer.Close()
	})

	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatal().Msgf("server start failed: %v", err)
	} else {
		log.Info().Msg("Received shutdown")
	}
}

func buildRouters() *chi.Mux {
	r := chi.NewRouter()
	r.Use(chiMiddleware.Heartbeat("/health"))

	// TODO remove
	message_provider.Instance()
	message.Configure(r)

	return r
}
