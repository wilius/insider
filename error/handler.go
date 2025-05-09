package error

import (
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
	"net/http"
	"reflect"
)

func ResponseError(w *http.ResponseWriter, r *http.Request, err interface{}) {
	if errCasted, ok := err.(error); ok {
		if casted, ok := err.(HttpCompatibleError); ok {
			RespondError(w, r, casted.Code(), casted.Title(), casted.HttpStatus(), errCasted)
		} else {
			RespondError(w, r, "UNEXPECTED_EXCEPTION", "Unexpected Exception", http.StatusInternalServerError, errCasted)
		}
	} else {
		RespondError(w, r, "UNEXPECTED_EXCEPTION", "Unexpected Exception", http.StatusInternalServerError, errors.New(fmt.Sprintf("%v", reflect.ValueOf(err))))
	}

}

func RespondError(w *http.ResponseWriter, r *http.Request, code string, title string, status int, err error) {
	render.Status(r, status)
	log.Error().Msgf("Unexpected error %v", err)
	render.JSON(*w, r, Response{
		Code:    code,
		Message: title,
		Detail:  err.Error(),
	})
}
