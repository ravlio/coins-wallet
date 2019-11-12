package account

import (
	eh "github.com/looplab/eventhorizon"
)

const (
	CreatedEvent      eh.EventType = "account:created"
	CreateFailedEvent eh.EventType = "account:createFailed"
)

type CreatedData struct {
	Name  string
	Email string
}

type CreatedFailedData struct {
	Error error
}
