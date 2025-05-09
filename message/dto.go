package message

import (
	"insider/types"
	"insider/util"
	"net/http"
	"time"
)

var dtoStatus = struct {
	Created messageStatus
	Sending messageStatus
	Sent    messageStatus
	Failed  messageStatus
}{
	Created: Created,
	Sending: Sending,
	Sent:    Sent,
	Failed:  Failed,
}

type createRequest struct {
	PhoneNumber string `json:"name" validate:"required,e164"`
	Message     string `json:"title" validate:"required,min=1,max=120"`
}

func (s *createRequest) Bind(_ *http.Request) error {
	return util.Validate(s)
}

type Dto struct {
	ID          *types.EntityId `json:"id"`
	PhoneNumber string          `json:"phoneNumber"`
	Message     string          `json:"message"`
	Status      messageStatus   `json:"status"`
	CreateDate  time.Time       `json:"createDate"`
	UpdateDate  time.Time       `json:"updateDate"`
}
