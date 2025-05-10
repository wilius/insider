package sender

import (
	"github.com/rs/zerolog/log"
	"github.com/wagslane/go-rabbitmq"
	"insider/configs"
	"insider/event_bus"
	"insider/message"
	"insider/provider"
	"time"
)

func runner(publisher *rabbitmq.Publisher) {
	doRun(publisher) // Run the task immediately
	ticker := time.NewTicker(instance.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ticker.Stop()
			doRun(publisher)
			ticker.Reset(instance.interval)
		case <-instance.stop:
			log.Info().
				Msg("Scheduler stopped")
			return
		}
	}
}

func doRun(publisher *rabbitmq.Publisher) {
	log.Info().
		Msg("Running scheduler")

	count := configs.Instance().
		GetScheduler().
		GetItemCountPerCycle()

	messageProvider := provider.Instance()

	// begin here and don't forget to add skip locked keyword into the downstream sql query
	items, err := messageService.Fetch(count)
	if err != nil {
		log.Err(err).
			Msg("Error while fetching items")
		return
	}

	for _, item := range *items {
		event_bus.PublishEvent(publisher, &eventDto{Id: item.ID}, "", "delayed_message")
	}
	// commit here

	for _, item := range *items {
		doSendMessage(&item, messageProvider)
	}
}

func doSendMessage(item *message.DTO, messageProvider provider.Provider) {
	log.Info().
		Msgf("Sending message to %s with content %s with id %d", item.PhoneNumber, item.Message, item.ID)

	inputItem := mapTo(item)

	response, err := messageProvider.Send(inputItem)
	if err != nil {
		log.Err(err).
			Msgf("Error while sending message to %s with content %s with id %d", inputItem.PhoneNumber, inputItem.Message, item.ID)
		err = messageService.MarkAsFailed(item.ID)
	} else {
		log.Info().
			Msgf("Message sent to %s with content %s with id %d tracked by %s", inputItem.PhoneNumber, inputItem.Message, item.ID, response.MessageId)
		err = messageService.MarkAsSent(item.ID, response.MessageId, messageProvider.Type())
	}

	if err != nil {
		log.Err(err).
			Msgf("Error while marking message with id %d", item.ID)
	}
}

func doHandleCheckingMessage(event *eventDto) error {
	log.Info().
		Msgf("Received message with ID: %d", event.Id)

	return messageService.MarkAsCreated(event.Id)
}
