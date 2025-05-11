package sender

import (
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	rabbitInternal "insider/rabbitmq"
)

func createDefaultExchangePublisher() *rabbitInternal.Publisher {
	publisher, err := rabbitInternal.NewPublisher("")

	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Unexpected exception while trying to configure publisher")
	}

	return publisher
}

func declareDelayedMessageQueue() {
	err := rabbitInternal.DeclareQueue(
		"delayed_message",
		true, false, false, false,
		amqp.Table{
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": "message.checker",
			"x-message-ttl":             60000,
		},
	)

	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Unexpected exception while trying to declare delayed_message queue")
	}
}
