package sender

import (
	"github.com/rs/zerolog/log"
	"insider/configs"
	"insider/provider"
	"time"
)

func runner() {
	doRun() // Run the task immediately
	ticker := time.NewTicker(instance.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ticker.Stop()
			doRun()
			ticker.Reset(instance.interval)
		case <-instance.stop:
			log.Info().
				Msg("Scheduler stopped")
			return
		}
	}
}

func doRun() {
	log.Info().
		Msg("Running scheduler")
	count := configs.Instance().
		GetScheduler().
		GetItemCountPerCycle()

	messageProvider := provider.Instance()
	items, err := messageService.Fetch(count)
	if err != nil {
		log.Err(err).
			Msg("Error while fetching items")
		return
	}

	for _, item := range *items {
		log.Info().
			Msgf("Sending message to %s with content %s", item.PhoneNumber, item.Message)

		inputItem := mapTo(&item)

		response, err := messageProvider.Send(inputItem)
		if err != nil {
			log.Err(err).
				Msgf("Error while sending message to %s with content %s", inputItem.PhoneNumber, inputItem.Message)
			err = messageService.MarkAsFailed(item.ID)
		} else {
			log.Info().
				Msgf("Message sent to %s with content %s tracked by %s", inputItem.PhoneNumber, inputItem.Message, response.MessageId)
			err = messageService.MarkAsSent(item.ID, response.MessageId, messageProvider.Type())
		}

		if err != nil {
			log.Err(err).
				Msgf("Error while marking message")
		}
	}
}
