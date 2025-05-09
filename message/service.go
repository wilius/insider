package message

import (
	"context"
	"insider/types"
)

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
