package sender

import (
	"github.com/go-chi/render"
	"net/http"
)

type handler struct {
}

func newHandler() *handler {
	return &handler{}
}

func (h *handler) Start(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusAccepted)
}

func (h *handler) Stop(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusAccepted)
}
