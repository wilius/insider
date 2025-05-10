package sender

import (
	"github.com/rs/zerolog/log"
	"insider/configs"
	"insider/database"
	"insider/message"
	"insider/provider"
	"time"
)

func newScheduler() *scheduler {
	messageService = message.GetUnsentMessageService()

	schedulerConfig := configs.Instance().
		GetScheduler()

	schedulerInstance = &scheduler{
		interval: schedulerConfig.GetInterval(),
	}

	schedulerInstance.Start()
	return schedulerInstance
}

type scheduler struct {
	interval time.Duration
	stop     chan struct{}
}

func (s scheduler) Start() {
	s.stop = make(chan struct{})
	go s.doStart()
}

func (s scheduler) doStart() {
	s.doRun() // Run the task immediately
	ticker := time.NewTicker(schedulerInstance.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ticker.Stop()
			s.doRun()
			ticker.Reset(schedulerInstance.interval)
		case <-schedulerInstance.stop:
			log.Info().
				Msg("Scheduler stopped")
			return
		}
	}
}

func (s scheduler) doRun() {
	log.Info().
		Msg("Running scheduler")

	itemsToSend, err := s.fetchMessages()
	if err != nil {
		return
	}

	for _, item := range *itemsToSend {
		s.doSendMessage(&item)
	}
}

func (s scheduler) fetchMessages() (*[]message.DTO, error) {
	tx := database.Instance().Begin()
	count := configs.Instance().
		GetScheduler().
		GetItemCountPerCycle()

	items, err := messageService.Fetch(tx, count)
	if err != nil {
		log.Err(err).
			Msg("Error while fetching items")
		return nil, err
	}

	itemsToSend := make([]message.DTO, 0)
	for _, item := range *items {
		err = publisher.PublishEvent(
			&eventDto{Id: item.ID},
			"delayed_message",
		)

		if err != nil {
			log.Err(err).
				Msg("Error while publishing event")
			err = messageService.MarkAsCreated(item.ID)
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

func (s scheduler) doSendMessage(item *message.DTO) {
	log.Info().
		Msgf("Sending message to %s with content %s with id %d", item.PhoneNumber, item.Message, item.ID)

	inputItem := mapTo(item)

	messageProvider := provider.Instance()

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

func (s scheduler) Stop() {
	close(schedulerInstance.stop)
}
