package rabbitmq

import (
	"github.com/wagslane/go-rabbitmq"
	"insider/graceful_shutdown"
)

func NewConsumer(
	queue string,
	optionFunctions ...func(*rabbitmq.ConsumerOptions),
) (*rabbitmq.Consumer, error) {
	configure()

	optionFunctions = append(
		optionFunctions,
		rabbitmq.WithConsumerOptionsConsumerAutoAck(false),
	)

	consumer, err := rabbitmq.NewConsumer(
		processingConnection,
		queue,
		optionFunctions...,
	)
	if err != nil {
		if consumer != nil {
			consumer.Close()
		}

		return nil, err
	}

	graceful_shutdown.AddShutdownHook(consumer.Close)
	return consumer, nil
}
