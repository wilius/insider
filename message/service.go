package message

import (
	"context"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"insider/constants"
	"insider/types"
	"time"
)

var redisService = &redisDataService{}

type UnsentMessageService interface {
	Fetch(tx *gorm.DB, count uint) (*[]DTO, error)
	MarkAsSent(id int64, providerMessageId string, provider constants.ProviderType) error
	MarkAsFailed(id int64) error
	MarkAsCreated(id int64) error
}

type dataService struct {
	repo *repository
}

func newDataServices(r *repository) *dataService {
	return &dataService{
		repo: r,
	}
}

func (s *dataService) Create(context context.Context, input *entity) (*entity, error) {
	return s.repo.Create(context, input)
}

func (s *dataService) List(ctx context.Context, filter *Filter) (*types.Pageable, error) {
	return s.repo.List(ctx, filter)
}

func (s *dataService) Fetch(tx *gorm.DB, count uint) (*[]DTO, error) {
	items, err := s.repo.FetchForSending(tx, count)
	if err != nil {
		return nil, err
	}

	mapped, err := types.MapToDTOList(items, mapToDTO)
	if err != nil {
		return nil, err
	}

	return mapped, err
}

func (s *dataService) MarkAsSent(id int64, providerMessageId string, provider constants.ProviderType) error {
	var providerType = string(provider)
	err := redisService.Store(context.Background(), id, providerMessageId, time.Now())
	if err != nil {
		log.Warn().
			Err(err).
			Msgf("Failed to store cache for message with id %d. skipping..", id)
	}

	return s.repo.update(
		id,
		constants.Sending,
		&entity{
			Status:            constants.Sent,
			ProviderMessageID: &providerMessageId,
			Provider:          &providerType,
		},
	)
}

func (s *dataService) MarkAsFailed(id int64) error {
	return s.repo.update(
		id,
		constants.Sending,
		&entity{
			Status: constants.Failed,
		},
	)
}

func (s *dataService) MarkAsCreated(id int64) error {
	result, err := redisService.Exists(context.Background(), id)
	if err == nil {
		if *result {
			return nil
		}
		log.Info().
			Msgf("There is no cached provider message result found for message with id %d. trying it with database update query..", id)
	} else {
		log.Warn().
			Err(err).
			Msgf("Failed to check cache for message with id %d. trying it with database update query..", id)
	}

	return s.repo.update(
		id,
		constants.Sending,
		&entity{
			Status: constants.Created,
		},
	)
}
