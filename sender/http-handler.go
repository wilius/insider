package sender

import (
	"github.com/go-chi/render"
	"insider/scheduler"
	"net/http"
)

type handler struct {
	manager *scheduler.Manager
}

func newHandler(*scheduler.Manager) *handler {
	return &handler{
		manager: manager,
	}
}

func (h *handler) Start(_ http.ResponseWriter, r *http.Request) {
	manager.Start()
	render.Status(r, http.StatusAccepted)
}

func (h *handler) Stop(_ http.ResponseWriter, r *http.Request) {
	manager.Stop()
	render.Status(r, http.StatusAccepted)
}
