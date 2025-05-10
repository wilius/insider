package sender

import (
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"github.com/wagslane/go-rabbitmq"
	error2 "insider/error"
	rabbitInternal "insider/rabbitmq"
)

func configureRabbitMQ() {
	declareDelayedMessageQueue()
	configureMessageStatusCheckerConsumer()
	configurePublisher()
}

func configurePublisher() {
	var err error
	publisher, err = rabbitInternal.NewPublisher("")

	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Unexpected exception while trying to configure publisher")
	}
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

func configureMessageStatusCheckerConsumer() {
	consumer, err := rabbitInternal.NewConsumer(
		"message.checker",
		rabbitmq.WithConsumerOptionsConcurrency(1),
		rabbitmq.WithConsumerOptionsQueueDurable,
	)

	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Unexpected exception while trying to configure consumer for router-entity")
	}

	go listenChanges(consumer)
}

func listenChanges(consumer *rabbitmq.Consumer) {
	err := consumer.Run(func(d rabbitmq.Delivery) rabbitmq.Action {
		event := eventDto{}
		if err := json.Unmarshal(d.Body, &event); err != nil {
			return rabbitmq.NackRequeue
		}

		err := doHandleCheckingMessage(&event)

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

func doHandleCheckingMessage(event *eventDto) error {
	log.Info().
		Msgf("Received message with ID: %d", event.Id)

	return messageService.MarkAsCreated(event.Id)
}
