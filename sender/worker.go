package sender

import (
	"github.com/rs/zerolog/log"
	"insider/configs"
	"insider/database"
	"insider/message"
	"insider/provider"
	"insider/rabbitmq"
)

type sendMessageWorker struct {
	publisher      *rabbitmq.Publisher
	messageService message.UnsentMessageService
	consumer       *delayedCheckConsumer
}

func newSendMessageWorker(
	publisher *rabbitmq.Publisher,
	messageService message.UnsentMessageService,
) *sendMessageWorker {

	consumer := newDelayedCheckConsumer(messageService)

	return &sendMessageWorker{
		publisher:      publisher,
		messageService: messageService,
		consumer:       consumer,
	}
}

func (w sendMessageWorker) doRun() {
	log.Info().
		Msg("Running scheduler")

	itemsToSend, err := w.fetchMessages()
	if err != nil {
		return
	}

	for _, item := range *itemsToSend {
		w.doSendMessage(&item)
	}
}

func (w sendMessageWorker) fetchMessages() (*[]message.DTO, error) {
	tx := database.Instance().Begin()
	count := configs.Instance().
		GetScheduler().
		GetItemCountPerCycle()

	items, err := w.messageService.Fetch(tx, count)
	if err != nil {
		log.Err(err).
			Msg("Error while fetching items")
		return nil, err
	}

	itemsToSend := make([]message.DTO, 0)
	for _, item := range *items {
		err = w.publisher.PublishEvent(
			&eventDto{Id: item.ID},
			"delayed_message",
		)

		if err != nil {
			log.Err(err).
				Msg("Error while publishing event")
			err = w.messageService.MarkAsCreated(item.ID)
			if err != nil {
				log.Err(err).
					Msg("Error while marking as created")
			}
		} else {
			itemsToSend = append(itemsToSend, item)
		}
	}
	tx.Commit()
	return &itemsToSend, nil
}

func (w sendMessageWorker) doSendMessage(item *message.DTO) {
	log.Info().
		Msgf("Sending message to %s with content %s with id %d", item.PhoneNumber, item.Message, item.ID)

	inputItem := mapTo(item)

	messageProvider := provider.Instance()

	response, err := messageProvider.Send(inputItem)
	if err != nil {
		log.Err(err).
			Msgf("Error while sending message to %s with content %s with id %d", inputItem.PhoneNumber, inputItem.Message, item.ID)
		err = w.messageService.MarkAsFailed(item.ID)
	} else {
		log.Info().
			Msgf("Message sent to %s with content %s with id %d tracked by %s", inputItem.PhoneNumber, inputItem.Message, item.ID, response.MessageId)
		err = w.messageService.MarkAsSent(item.ID, response.MessageId, messageProvider.Type())
	}

	if err != nil {
		log.Err(err).
			Msgf("Error while marking message with id %d", item.ID)
	}
}
