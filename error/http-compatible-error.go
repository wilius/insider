package error

// HttpCompatibleError an interface used by error classes which supposed to be compatible with error response.
// added here in order to get rid of one-by-one error mapping for error responses. instead, the interface provides
// the desired values in order to respond an error for a api client interaction
type HttpCompatibleError interface {
	// Code corresponds what to show as "code" property of an error response
	Code() string

	// Title corresponds what to show as "title" property of an error response
	Title() string

	// HttpStatus indicates which http status code have to return back to a client interaction
	HttpStatus() int
}
