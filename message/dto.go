package message

import (
	"insider/constants"
	"insider/types"
	"insider/util"
	"net/http"
	"time"
)

var dtoStatus = struct {
	Created constants.MessageStatus
	Sending constants.MessageStatus
	Sent    constants.MessageStatus
	Failed  constants.MessageStatus
}{
	Created: constants.Created,
	Sending: constants.Sending,
	Sent:    constants.Sent,
	Failed:  constants.Failed,
}

type createRequest struct {
	PhoneNumber string `json:"phoneNumber" validate:"required,e164"`
	Message     string `json:"message" validate:"required,min=1,max=120"`
}

func (s *createRequest) Bind(_ *http.Request) error {
	return util.Validate(s)
}

type HttpDTO struct {
	ID          types.EntityId          `json:"id"`
	PhoneNumber string                  `json:"phoneNumber"`
	Message     string                  `json:"message"`
	Status      constants.MessageStatus `json:"status"`
	CreateDate  time.Time               `json:"createDate"`
	UpdateDate  *time.Time              `json:"updateDate"`
}

type DTO struct {
	ID          int64
	PhoneNumber string
	Message     string
	Status      constants.MessageStatus
	CreateDate  time.Time
	UpdateDate  *time.Time
}

type Filter struct {
	types.PagedFilter
	Query  *string
	Status *string
}
