package message

import (
	"github.com/go-chi/render"
	customError "insider/error"
	"insider/types"
	"net/http"
)

type handler struct {
	service *dataService
}

func newHandler(service *dataService) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	data := &createRequest{}
	if err := render.Bind(r, data); err != nil {
		customError.ResponseError(&w, r, customError.NewInvalidParamsError(err.Error()))
		return
	}

	create := mapCreate(data)
	output, err := h.service.Create(
		r.Context(),
		create,
	)

	if err != nil {
		customError.ResponseError(&w, r, customError.NewProcessingError(err.Error()))
		return
	}

	render.Status(r, http.StatusCreated)
	dto, err := mapToHttpDTO(output)
	if err != nil {
		customError.ResponseError(&w, r, customError.NewProcessingError(err.Error()))
		return
	}

	render.JSON(w, r, dto)
}

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	filter := mapQueryToFilter(&query)

	pagedResult, err := h.service.List(r.Context(), filter)
	if err != nil {
		customError.ResponseError(&w, r, err)
		return
	}

	casted, err := types.MapToPageDTO(pagedResult, mapToHttpDTO)
	if err != nil {
		customError.ResponseError(&w, r, err)
		return
	}
	render.JSON(w, r, casted)
}
