package sender

import (
	"context"
	"github.com/rs/zerolog/log"
	"insider/configs"
	"insider/provider"
	"insider/types"
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
	items, err := messageService.Fetch(context.Background(), count)
	if err != nil {
		log.Err(err).
			Msg("Error while fetching items")
		return
	}

	inputItems, err := types.MapToDTOList(items, mapTo)
	if err != nil {
		log.Err(err).
			Msg("Error while fetching items")
		return
	}

	for _, inputItem := range *inputItems {
		log.Info().
			Msgf("Sending message to %s with content %s", inputItem.PhoneNumber, inputItem.Message)
		response, err := messageProvider.Send(&inputItem)
		if err != nil {
			log.Err(err).
				Msgf("Error while sending message to %s with content %s", inputItem.PhoneNumber, inputItem.Message)
		} else {
			log.Info().
				Msgf("Message sent to %s with content %s tracked by %s", inputItem.PhoneNumber, inputItem.Message, response.MessageId)
		}

	}
}
