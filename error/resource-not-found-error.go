package error

import "net/http"

// ResourceNotFound a custom error type if the desired resource doesn't exist
type ResourceNotFound struct {
	message string
}

func (e ResourceNotFound) Error() string {
	return e.message
}

func (e ResourceNotFound) Code() string {
	return "RESOURCE_NOT_FOUND"
}

func (e ResourceNotFound) Title() string {
	return "Resource Not Found"
}

func (e ResourceNotFound) HttpStatus() int {
	return http.StatusNotFound
}

func NewResourceNotFound(message string) ResourceNotFound {
	return ResourceNotFound{
		message: message,
	}
}
