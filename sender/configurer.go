package sender

import (
	"github.com/go-chi/chi/v5"
	"insider/graceful_shutdown"
	"insider/message"
	"insider/rabbitmq"
	"sync"
)

var (
	schedulerInstance *scheduler
	connectionOnce    sync.Once
	messageService    message.UnsentMessageService
	publisher         *rabbitmq.Publisher
)

func Configure(mux *chi.Mux) {
	connectionOnce.Do(func() {
		schedulerInstance = newScheduler()
		configureRabbitMQ()
		graceful_shutdown.AddShutdownHook(func() {
			schedulerInstance.Stop()
		})
	})

	handlerInstance := newHandler()

	mux.Get("/messages", handlerInstance.Start)
	mux.Post("/messages", handlerInstance.Stop)
}
