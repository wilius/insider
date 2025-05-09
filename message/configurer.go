package message

import (
	"github.com/go-chi/chi/v5"
	"insider/database"
)

func Configure(mux *chi.Mux) {
	handler := newHandler(
		newDataServices(
			newRepository(
				database.Instance(),
			),
		),
	)

	mux.Get("/message", handler.List)
	mux.Post("/message", handler.Create)
}
