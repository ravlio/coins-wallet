package balance

import (
	"context"
	"errors"
	"fmt"

	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/eventhandler/projector"
)

type Projector struct{}

func NewProjector() *Projector {
	return &Projector{}
}

func (p *Projector) ProjectorType() projector.Type {
	return "BalanceProjector"
}

func (p *Projector) Project(ctx context.Context, event eh.Event, entity eh.Entity) (eh.Entity, error) {
	i, ok := entity.(*Balance)
	if !ok {
		return nil, errors.New("model is of incorrect type")
	}

	switch event.EventType() {
	case CreditedEvent:
		data, ok := event.Data().(*CreditedData)
		if !ok {
			return nil, fmt.Errorf("projector: invalid event data type: %v", event.Data())
		}

		i.ID = event.AggregateID()
		i.AccountID = data.AccountID
		i.Currency = data.Currency
		i.Balance += data.Amount

	case DebitedEvent:
		data, ok := event.Data().(*DebitedData)
		if !ok {
			return nil, fmt.Errorf("projector: invalid event data type: %v", event.Data())
		}

		i.ID = event.AggregateID()
		i.AccountID = data.AccountID
		i.Currency = data.Currency
		i.Balance -= data.Amount

	default:
		return nil, errors.New("could not handle event: " + event.String())
	}

	i.Version++
	return i, nil
}
