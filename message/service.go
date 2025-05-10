package message

import (
	"context"
	"insider/constants"
	"insider/types"
)

type UnsentMessageService interface {
	Fetch(count uint) (*[]DTO, error)
	MarkAsSent(id int64, providerMessageId string, provider constants.ProviderType) error
	MarkAsFailed(id int64) error
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

func (s *dataService) Fetch(count uint) (*[]DTO, error) {
	items, err := s.repo.FetchForSending(count)
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
