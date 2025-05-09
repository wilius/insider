package error

import "net/http"

// InvalidParamsError an error created when the given input for any kind of api interaction of client contains undesired value for any input property
type InvalidParamsError struct {
	message string
}

func (e InvalidParamsError) Error() string {
	return e.message
}

func (e InvalidParamsError) Code() string {
	return "INVALID_PARAMS"
}

func (e InvalidParamsError) Title() string {
	return "Invalid params"
}

func (e InvalidParamsError) HttpStatus() int {
	return http.StatusBadRequest
}

func NewInvalidParamsError(message string) InvalidParamsError {
	return InvalidParamsError{
		message: message,
	}
}
