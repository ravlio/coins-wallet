package account

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
	return "AccountProjector"
}

func (p *Projector) Project(ctx context.Context, event eh.Event, entity eh.Entity) (eh.Entity, error) {
	i, ok := entity.(*Account)
	if !ok {
		return nil, errors.New("model is of incorrect type")
	}

	switch event.EventType() {
	case CreatedEvent:
		data, ok := event.Data().(*CreatedData)
		if !ok {
			return nil, fmt.Errorf("projector: invalid event data type: %v", event.Data())
		}

		i.ID = event.AggregateID()
		i.Name = data.Name
		i.Email = data.Email

	default:
		return nil, errors.New("could not handle event: " + event.String())
	}

	i.Version++
	return i, nil
}
