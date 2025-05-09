package error

// Response a DTO object to use in case of to respond error to any endpoint call
type Response struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}
