package message

import (
	"fmt"
	"insider/constants"
	customError "insider/error"
	"insider/types"
	"net/url"
)

func mapToHttpDTO(item *entity) (*HttpDTO, error) {
	status, err := mapStatus(&item.Status)
	if err != nil {
		return nil, err
	}

	d := &HttpDTO{
		ID:          types.MapEntityId(&item.ID),
		PhoneNumber: item.PhoneNumber,
		Message:     item.Message,
		Status:      *status,
		CreateDate:  item.CreateDate,
		UpdateDate:  item.UpdateDate,
	}

	return d, nil
}

func mapToDTO(item *entity) (*DTO, error) {
	status, err := mapStatus(&item.Status)
	if err != nil {
		return nil, err
	}

	d := &DTO{
		ID:          item.ID,
		PhoneNumber: item.PhoneNumber,
		Message:     item.Message,
		Status:      *status,
		CreateDate:  item.CreateDate,
		UpdateDate:  item.UpdateDate,
	}

	return d, nil
}

func mapStatus(status *constants.MessageStatus) (*constants.MessageStatus, error) {
	switch *status {
	case constants.Created:
		return &dtoStatus.Created, nil
	case constants.Sending:
		return &dtoStatus.Sending, nil
	case constants.Sent:
		return &dtoStatus.Sent, nil
	case constants.Failed:
		return &dtoStatus.Failed, nil
	default:
		return nil, customError.NewProcessingError(
			fmt.Sprintf("Could not determine message status from %s", *status),
		)
	}
}

func mapCreate(data *createRequest) *entity {
	item := &entity{
		PhoneNumber: data.PhoneNumber,
		Message:     data.Message,
		Status:      constants.Created,
	}

	return item
}

func mapQueryToFilter(valuesPtr *url.Values) *Filter {
	status := getFirst(valuesPtr, "status")

	if status == nil {
		tmp := constants.Sent
		status = &tmp
	}

	return &Filter{
		PagedFilter: *types.ParseQueryForPageFilter(valuesPtr),
		Query:       getFirst(valuesPtr, "query"),
		Status:      status,
	}
}

func getFirst(valuesPtr *url.Values, key string) *string {
	if query, exists := (*valuesPtr)[key]; exists && len(query) > 0 {
		return &query[0]
	}

	return nil
}
