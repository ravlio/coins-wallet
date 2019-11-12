package balance

import (
	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
	"github.com/ravlio/wallet/pkg/money"
)

const (
	CreditedEvent     eh.EventType = "balance:credited"
	DebitedEvent      eh.EventType = "balance:debited"
	CreditFailedEvent eh.EventType = "balance:creditFailed"
	DebitFailedEvent  eh.EventType = "balance:debitFailed"
)

type CommonData struct {
	AccountID uuid.UUID
	Currency  money.Currency
	Amount    money.Money
}
type CreditedData struct {
	CommonData
}

type DebitedData struct {
	CommonData
}

type CreditFailedData struct {
	CommonData
	Error error
}

type DebitFailedData struct {
	CommonData
	Error error
}
