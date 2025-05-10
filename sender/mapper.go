package sender

import (
	"insider/message"
	"insider/provider"
)

func mapTo(dto *message.Dto) (*provider.SendMessageInput, error) {
	return &provider.SendMessageInput{
		PhoneNumber: dto.PhoneNumber,
		Message:     dto.Message,
	}, nil
}
