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

	mux.Get("/messages", handler.List)
	mux.Post("/messages", handler.Create)
}
