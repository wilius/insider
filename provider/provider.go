package provider

import "insider/constants"

type Provider interface {
	Type() constants.ProviderType
	Send(input *SendMessageInput) (*SendMessageOutput, error)
}

type SendMessageInput struct {
	PhoneNumber string
	Message     string
}

type SendMessageOutput struct {
	MessageId string
}
