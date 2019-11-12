package balance

import (
	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
	"github.com/ravlio/wallet/pkg/money"
)

type Balance struct {
	ID        uuid.UUID
	AccountID uuid.UUID
	Currency  money.Currency
	Balance   money.Money
	Version   int
}

var _ = eh.Entity(&Balance{})
var _ = eh.Versionable(&Balance{})

func (i *Balance) EntityID() uuid.UUID {
	return i.ID
}

func (i *Balance) AggregateVersion() int {
	return i.Version
}
