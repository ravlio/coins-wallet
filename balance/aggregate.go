package balance

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/aggregatestore/events"
	"github.com/ravlio/wallet/pkg/money"
)

const AggregateType = "balance"

type Aggregate struct {
	*events.AggregateBase
	AccountID uuid.UUID
	Currency  money.Currency
	Balance   money.Money
}

var _ = eh.Aggregate(&Aggregate{})

func NewAggregate(id uuid.UUID) *Aggregate {
	return &Aggregate{
		AggregateBase: events.NewAggregateBase(AggregateType, id),
	}
}

func (a *Aggregate) HandleCommand(ctx context.Context, cmd eh.Command) error {
	var event eh.EventType
	var data eh.EventData

	switch cmd := cmd.(type) {
	case *Credit:
		cd := CommonData{
			AccountID: cmd.AccountID,
			Currency:  cmd.Currency,
			Amount:    cmd.Amount,
		}
		if !money.IsCurrency(cmd.Currency) {
			event = CreditFailedEvent
			data = &CreditFailedData{
				CommonData: cd,
				Error:      fmt.Errorf("invalid currency"),
			}
		} else if cmd.Amount < 0 {
			event = CreditFailedEvent
			data = &CreditFailedData{
				CommonData: cd,
				Error:      fmt.Errorf("negative amount"),
			}
		} else {
			event = CreditedEvent
			data = &CreditedData{
				CommonData: cd,
			}
		}

	case *Debit:
		cd := CommonData{
			AccountID: cmd.AccountID,
			Currency:  cmd.Currency,
			Amount:    cmd.Amount,
		}
		if !money.IsCurrency(cmd.Currency) {
			event = DebitFailedEvent
			data = &DebitFailedData{
				CommonData: cd,
				Error:      fmt.Errorf("invalid currency"),
			}
		} else if cmd.Amount < 0 {
			event = DebitFailedEvent
			data = &DebitFailedData{
				CommonData: cd,
				Error:      fmt.Errorf("negative amount"),
			}
		} else if a.Balance < cmd.Amount {
			event = DebitFailedEvent
			data = &DebitFailedData{
				CommonData: cd,
				Error:      fmt.Errorf("balance %d < amount %d", a.Balance, cmd.Amount),
			}
		} else {
			event = DebitedEvent
			data = &DebitedData{
				CommonData: cd,
			}
		}

	default:
		return errors.New("unknown command")
	}

	a.StoreEvent(event, data, time.Now())

	return nil

}

func (a *Aggregate) ApplyEvent(ctx context.Context, event eh.Event) error {
	switch event.EventType() {
	case CreditedEvent:
		if data, ok := event.Data().(*CreditedData); ok {
			a.AccountID = data.AccountID
			a.Currency = data.Currency
			a.Balance += data.Amount
		} else {
			return fmt.Errorf("invalid event data type: %v", event.Data())
		}

	case DebitedEvent:
		if data, ok := event.Data().(*CreditedData); ok {
			a.AccountID = data.AccountID
			a.Currency = data.Currency
			a.Balance -= data.Amount
		} else {
			return fmt.Errorf("invalid event data type: %v", event.Data())
		}
	}

	return nil
}
