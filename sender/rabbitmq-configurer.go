package sender

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"github.com/wagslane/go-rabbitmq"
	"insider/event_bus"
)

func configureRabbitMQ() *rabbitmq.Publisher {
	declareDelayedMessageQueue()
	configureMessageStatusCheckerConsumer()
	return configurePublisher()

}

func configurePublisher() *rabbitmq.Publisher {
	eventPublisher, err := event_bus.NewPublisher()

	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Unexpected exception while trying to configure publisher")
	}

	return eventPublisher
}

func declareDelayedMessageQueue() {
	err := event_bus.DeclareQueue(
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
	consumer, err := event_bus.NewConsumer(
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
			return rabbitmq.NackRequeue
		}

		return rabbitmq.Ack
	})

	if err != nil {
		log.Fatal().Err(err).Msgf("Unexpected exception while running consumer")
	}
}
