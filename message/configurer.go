package message

import (
	"github.com/go-chi/chi/v5"
	"insider/database"
)

var repo *repository
var service *dataService

func GetUnsentMessageService() UnsentMessageService {
	if service == nil {
		panic("service not initialized yet")
	}

	return service
}

func Configure(mux *chi.Mux) {
	repo = newRepository(
		database.Instance(),
	)

	service = newDataServices(repo)

	handler := newHandler(service)

	mux.Get("/messages", handler.List)
	mux.Post("/messages", handler.Create)
}
