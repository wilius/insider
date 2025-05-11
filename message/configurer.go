package message

import (
	"github.com/go-chi/chi/v5"
	"insider/database"
	"sync"
)

var (
	repo    *repository
	service *dataService
	once    sync.Once
)

func GetUnsentMessageService() UnsentMessageService {
	configureService()

	return service
}

func Configure(mux *chi.Mux) {
	configureService()

	handler := newHandler(service)

	mux.Get("/messages", handler.List)
	mux.Post("/messages", handler.Create)
}

func configureService() {
	once.Do(func() {
		repo = newRepository(
			database.Instance(),
		)

		service = newDataServices(repo)
	})
}
