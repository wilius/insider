package message

import (
	"context"
	"insider/types"
)

type UnsentMessageService interface {
	Fetch(ctx context.Context, count uint) (*[]Dto, error)
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

func (s *dataService) Fetch(ctx context.Context, count uint) (*[]Dto, error) {
	items, err := s.repo.FetchForSending(ctx, count)
	if err != nil {
		return nil, err
	}

	mapped, err := types.MapToDTOList(items, mapToDTO)
	if err != nil {
		return nil, err
	}

	return mapped, err
}
