package types

import (
	"encoding/json"
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

func (d *EntityId) MarshalJSON() ([]byte, error) {
	asString := strconv.FormatUint(uint64(d.Id), 10)
	encrypt, err := Encrypt(&asString)
	if err != nil {
		return nil, err
	}

	return json.Marshal(*encrypt)
}
