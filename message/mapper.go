package message

import (
	"fmt"
	customError "insider/error"
	"insider/types"
	"net/url"
)

func mapToDTO(item *entity) (*Dto, error) {
	status, err := mapStatus(&item.Status)
	if err != nil {
		return nil, err
	}

	d := &Dto{
		ID: types.EntityId{
			Id: item.ID,
		},
		PhoneNumber: item.PhoneNumber,
		Message:     item.Message,
		Status:      *status,
		CreateDate:  item.CreateDate,
		UpdateDate:  item.UpdateDate,
	}

	return d, nil
}

func mapStatus(status *messageStatus) (*messageStatus, error) {
	switch *status {
	case Created:
		return &dtoStatus.Created, nil
	case Sending:
		return &dtoStatus.Sending, nil
	case Sent:
		return &dtoStatus.Sent, nil
	case Failed:
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
		Status:      Created,
	}

	return item
}

func mapQueryToFilter(valuesPtr *url.Values) *Filter {
	return &Filter{
		PagedFilter: *types.ParseQueryForPageFilter(valuesPtr),
		Query:       getFirst(valuesPtr, "query"),
	}
}

func getFirst(valuesPtr *url.Values, key string) *string {
	if query, exists := (*valuesPtr)[key]; exists && len(query) > 0 {
		return &query[0]
	}

	return nil
}
