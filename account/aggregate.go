package account

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/aggregatestore/events"
)

const AggregateType = "account"

type Aggregate struct {
	*events.AggregateBase

	name  string
	email string
}

var _ = eh.Aggregate(&Aggregate{})

func NewAggregate(id uuid.UUID) *Aggregate {
	return &Aggregate{
		AggregateBase: events.NewAggregateBase(AggregateType, id),
	}
}

func (a *Aggregate) HandleCommand(ctx context.Context, cmd eh.Command) error {
	switch cmd := cmd.(type) {
	case *Create:
		a.StoreEvent(CreatedEvent, &CreatedData{
			Name:  cmd.Name,
			Email: cmd.Email,
		}, time.Now())

		return nil
	}

	return errors.New("couldn't handle command")
}

func (a *Aggregate) ApplyEvent(ctx context.Context, event eh.Event) error {
	switch event.EventType() {
	case CreatedEvent:
		if data, ok := event.Data().(*CreatedData); ok {
			a.name = data.Name
			a.email = data.Email
		} else {
			return fmt.Errorf("invalid event data type: %v", event.Data())
		}
	}

	return nil
}
