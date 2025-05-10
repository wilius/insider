package message_provider

type ProviderType string

const (
	WebhookSite ProviderType = "webhook.site"
)

type Provider interface {
	Send(input *SendMessageInput) (*SendMessageOutput, error)
}

type SendMessageInput struct {
	phoneNumber string
	message     string
}

type SendMessageOutput struct {
	messageId string
}
