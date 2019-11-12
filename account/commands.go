package account

import (
	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
)

const (
	CreateCommand eh.CommandType = "Create"
)

type Create struct {
	ID    uuid.UUID
	Name  string
	Email string
}

var _ = eh.Command(&Create{})

func (c *Create) AggregateType() eh.AggregateType { return AggregateType }
func (c *Create) AggregateID() uuid.UUID          { return c.ID }
func (c *Create) CommandType() eh.CommandType     { return CreateCommand }
