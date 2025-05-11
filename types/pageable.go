package types

import "errors"

type Pageable struct {
	First            bool        `json:"first"`
	Last             bool        `json:"last"`
	Size             int         `json:"size"`
	Content          interface{} `json:"content"`
	Page             int         `json:"page"`
	NumberOfElements int         `json:"numberOfElements"`
}

func MapToPageDTO[T any, K any](
	pagedResult *Pageable,
	mapper func(*T) (*K, error),
) (*Pageable, error) {
	if items, ok := pagedResult.Content.(*[]T); ok {
		mapped, err := MapToDTOList(items, mapper)
		if err != nil {
			return nil, err
		}

		pagedResult.Content = mapped
		return pagedResult, nil
	}

	return nil, errors.New("couldn't cast content")
}

func MapToDTOList[T any, K any](
	items *[]T,
	mapper func(*T) (*K, error)) (*[]K, error) {
	dtoList := make([]K, len(*items))

	// Convert each entity to a DTO
	for i, item := range *items {
		mapped, err := mapper(&item)
		if err != nil {
			return nil, err
		}
		dtoList[i] = *mapped
	}

	return &dtoList, nil
}

func MapToPageable[T any](items *[]T, filter *PagedFilter) *Pageable {
	size := filter.Size
	page := filter.Page
	dereference := *items
	itemLen := len(dereference)
	hasNextPage := itemLen > size

	if itemLen > size {
		itemLen = size
	}

	totalPages := page
	totalElements := (page-1)*size + itemLen
	if hasNextPage {
		dereference = dereference[:size]
		totalPages = totalPages + 1
		totalElements += 1
	}

	dereference = dereference[0:itemLen]
	pageable := Pageable{
		First:            page == 1,
		Last:             !hasNextPage,
		Size:             size,
		Content:          &dereference,
		Page:             page,
		NumberOfElements: itemLen,
	}

	return &pageable
}
