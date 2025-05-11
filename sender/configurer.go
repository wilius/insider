package sender

import (
	"github.com/go-chi/chi/v5"
	"insider/graceful_shutdown"
	"insider/message"
	"insider/scheduler"
)

var (
	manager *scheduler.Manager
	worker  *sendMessageWorker
)

func Configure(mux *chi.Mux) {
	declareDelayedMessageQueue()

	publisher := createDefaultExchangePublisher()

	worker = newSendMessageWorker(
		publisher,
		message.GetUnsentMessageService(),
	)

	manager = scheduler.NewManager(worker.doRun)

	graceful_shutdown.AddShutdownHook(func() {
		manager.Stop()
	})

	handlerInstance := newHandler(manager)

	mux.Put("/sender/start", handlerInstance.Start)
	mux.Put("/sender/stop", handlerInstance.Stop)
}
