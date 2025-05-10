package constants

type MessageStatus string

const (
	Created MessageStatus = "CREATED"
	Sending               = "SENDING"
	Sent                  = "SENT"
	Failed                = "FAILED"
)
