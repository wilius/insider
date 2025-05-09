package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	"insider/configs"
	"insider/database"
	"insider/graceful_shutdown"
	"net/http"
)

var httpServer *http.Server

func StartServer() {
	serverConfig := configs.Instance().GetServer()

	port := fmt.Sprintf(":%d", serverConfig.GetPort())

	log.Info().
		Msgf("starting server at localhost:%s", port)

	httpServer = &http.Server{
		Addr:    port,
		Handler: buildRouter(),
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

func buildRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(chiMiddleware.Heartbeat("/health"))

	database.Configure()
	return r
}
