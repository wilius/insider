package types

import (
	"encoding/json"
	"errors"
	"github.com/lib/pq"
	"regexp"
	"strconv"
)

type EntityId struct {
	Id int64
}

func MapEntityId(id *int64) *EntityId {
	if id != nil {
		return &EntityId{
			Id: *id,
		}
	}

	return nil
}

func UnwrapEntityId(id *EntityId) *int64 {
	if id != nil {
		return &id.Id
	}

	return nil
}

func UnwrapEntityIds(ids *[]EntityId) *[]int64 {
	if ids == nil {
		return nil
	}

	converted := make([]int64, len(*ids))
	for i, id := range *ids {
		converted[i] = *UnwrapEntityId(&id)
	}

	return &converted
}

func UnwrapEntityIdsBy[T any](items *[]T, fn func(*T) *EntityId) *[]int64 {
	if items == nil {
		return nil
	}

	converted := make([]int64, len(*items))
	for i, item := range *items {
		id := fn(&item)
		converted[i] = *UnwrapEntityId(id)
	}

	return &converted
}

func MapToEntityIdsToInt64(categories *[]EntityId) *pq.Int64Array {
	if categories == nil {
		return nil
	}

	dtoList := make(pq.Int64Array, len(*categories))
	for i, item := range *categories {
		dtoList[i] = item.Id
	}

	return &dtoList
}

func MapInt64ToEntityId(categories *pq.Int64Array) *[]EntityId {
	if categories == nil {
		return nil
	}

	dtoList := make([]EntityId, len(*categories))
	for i, item := range *categories {
		dtoList[i] = EntityId{
			Id: item,
		}
	}

	return &dtoList
}

func UnmarshalEntityId(str *string) (*EntityId, error) {
	item := EntityId{}

	if err := item.Unmarshal(str); err != nil {
		return nil, err
	}

	return &item, nil
}

func (d *EntityId) UnmarshalJSON(b []byte) error {
	text := ""
	err := json.Unmarshal(b, &text)
	if err != nil {
		return err
	}

	return d.Unmarshal(&text)
}

func (d *EntityId) Unmarshal(text *string) error {
	re := regexp.MustCompile(`^([^-]+-)*([^-]+)$`)

	// Find the fragment
	match := re.FindStringSubmatch(*text)
	l := len(match)
	if l != 3 {
		return errors.New("invalid id")
	}

	decrypt, err := Decrypt(&match[2])
	if err != nil {
		return err
	}

	data, err := strconv.ParseInt(*decrypt, 10, 64)
	if err != nil {
		return err
	}

	d.Id = data

	return nil
}

func (d *EntityId) MarshalJSON() ([]byte, error) {
	asString := strconv.FormatUint(uint64(d.Id), 10)
	encrypt, err := Encrypt(&asString)
	if err != nil {
		return nil, err
	}

	return json.Marshal(*encrypt)
}

func (d *EntityId) Marshal() (*string, error) {
	asString := strconv.FormatUint(uint64(d.Id), 10)
	encrypt, err := Encrypt(&asString)
	if err != nil {
		return nil, err
	}

	return encrypt, nil
}

func DecodeEntityId(str *string) (*int64, error) {
	if str == nil {
		return nil, nil
	}

	id, err := UnmarshalEntityId(str)
	if err != nil {
		return nil, err
	}

	return &id.Id, err
}

func DecodeEntityIds(strs *[]string) (*[]int64, error) {
	if strs == nil {
		return nil, nil
	}

	result := make([]int64, 0)
	for _, str := range *strs {
		id, err := UnmarshalEntityId(&str)
		if err != nil {
			return nil, err
		}

		result = append(result, id.Id)
	}

	return &result, nil
}
