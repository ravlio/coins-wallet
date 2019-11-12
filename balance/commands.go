package balance

import (
	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
	"github.com/ravlio/wallet/pkg/money"
)

const (
	CreditCommand eh.CommandType = "Credit"
	DebitCommand  eh.CommandType = "Debit"
)

type Credit struct {
	ID        uuid.UUID
	AccountID uuid.UUID
	Currency  money.Currency
	Amount    money.Money
}

type Debit struct {
	ID       uuid.UUID
	AccountID uuid.UUID
	Currency money.Currency
	Amount   money.Money
}

var _ = eh.Command(&Credit{})

func (c *Credit) AggregateType() eh.AggregateType { return AggregateType }
func (c *Credit) AggregateID() uuid.UUID          { return c.ID }
func (c *Credit) CommandType() eh.CommandType     { return CreditCommand }

var _ = eh.Command(&Debit{})

func (c *Debit) AggregateType() eh.AggregateType { return AggregateType }
func (c *Debit) AggregateID() uuid.UUID          { return c.ID }
func (c *Debit) CommandType() eh.CommandType     { return DebitCommand }
