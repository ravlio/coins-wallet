package account

import (
	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
)

type Account struct {
	ID      uuid.UUID
	Version int
	Name    string
	Email   string
}

var _ = eh.Entity(&Account{})
var _ = eh.Versionable(&Account{})

func (i *Account) EntityID() uuid.UUID {
	return i.ID
}

func (i *Account) AggregateVersion() int {
	return i.Version
}
