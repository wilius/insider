package error

import "net/http"

// ProcessingError custom error used when encountered with an unexpected error
type ProcessingError struct {
	message string
}

func (e ProcessingError) Error() string {
	return e.message
}

func (e ProcessingError) Code() string {
	return "PROCESSING_ERROR"
}

func (e ProcessingError) Title() string {
	return "Processing error"
}

func (e ProcessingError) HttpStatus() int {
	return http.StatusUnprocessableEntity
}

func NewProcessingError(message string) ProcessingError {
	return ProcessingError{
		message: message,
	}
}
