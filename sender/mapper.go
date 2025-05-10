package sender

import (
	"insider/message"
	"insider/provider"
)

func mapTo(dto *message.DTO) *provider.SendMessageInput {
	return &provider.SendMessageInput{
		PhoneNumber: dto.PhoneNumber,
		Message:     dto.Message,
	}
}
