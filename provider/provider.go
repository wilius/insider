package provider

type Provider interface {
	Send(input *SendMessageInput) (*SendMessageOutput, error)
}

type SendMessageInput struct {
	PhoneNumber string
	Message     string
}

type SendMessageOutput struct {
	MessageId string
}
