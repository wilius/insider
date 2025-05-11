package sender

import (
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"github.com/wagslane/go-rabbitmq"
	error2 "insider/error"
	"insider/message"
	rabbitInternal "insider/rabbitmq"
)

type delayedCheckConsumer struct {
	messageService message.UnsentMessageService
}

func newDelayedCheckConsumer(messageService message.UnsentMessageService) *delayedCheckConsumer {
	rabbitmqConsumer, err := rabbitInternal.NewConsumer(
		"message.checker",
		rabbitmq.WithConsumerOptionsConcurrency(1),
		rabbitmq.WithConsumerOptionsQueueDurable,
	)

	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Unexpected exception while trying to configure consumer for router-entity")
	}

	consumer := &delayedCheckConsumer{
		messageService: messageService,
	}

	consumer.Start(rabbitmqConsumer)

	return consumer
}

func (c *delayedCheckConsumer) Start(consumer *rabbitmq.Consumer) {
	go c.doStart(consumer)
}

func (c *delayedCheckConsumer) doStart(consumer *rabbitmq.Consumer) {
	err := consumer.Run(func(d rabbitmq.Delivery) rabbitmq.Action {
		event := eventDto{}
		if err := json.Unmarshal(d.Body, &event); err != nil {
			return rabbitmq.NackRequeue
		}

		err := c.doHandleCheckingMessage(&event)

		if err != nil {
			var casted error2.ResourceNotFound
			if !errors.As(err, &casted) {
				return rabbitmq.NackRequeue
			}
		}

		return rabbitmq.Ack
	})

	if err != nil {
		log.Fatal().Err(err).Msgf("Unexpected exception while running consumer")
	}
}

func (c *delayedCheckConsumer) doHandleCheckingMessage(event *eventDto) error {
	log.Info().
		Msgf("Received message with ID: %d", event.Id)

	return c.messageService.MarkAsCreated(event.Id)
}
