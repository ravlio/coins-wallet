package account

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/aggregatestore/events"
	"github.com/looplab/eventhorizon/commandhandler/aggregate"
	"github.com/looplab/eventhorizon/commandhandler/bus"
	"github.com/looplab/eventhorizon/eventhandler/projector"
	"github.com/ravlio/wallet/pkg/middleware"
)

type Client interface {
	CreateAccount(ctx context.Context, req *CreateAccountRequest) error
	GetAccount(ctx context.Context, id string) (*Account, error)
}

type CreateAccountRequest struct {
	ID    uuid.UUID
	Name  string
	Email string
}

type Service struct {
	commandHandler eh.CommandHandler
	repo           eh.ReadRepo
}

func NewService(eventStore eh.EventStore,
	eventBus eh.EventBus,
	commandBus *bus.CommandHandler,
	repo eh.ReadWriteRepo) (*Service, error) {

	svc := &Service{repo: repo}

	eventBus.AddObserver(eh.MatchAny(), &middleware.Logger{})

	aggregateStore, err := events.NewAggregateStore(eventStore, eventBus)
	if err != nil {
		return nil, fmt.Errorf("could not create aggregate store: %w", err)
	}

	handler, err := aggregate.NewCommandHandler(AggregateType, aggregateStore)
	if err != nil {
		return nil, fmt.Errorf("could not create command commandHandler: %w", err)
	}

	svc.commandHandler = eh.UseCommandHandlerMiddleware(handler, middleware.LoggingMiddleware)
	err = commandBus.SetHandler(svc.commandHandler, CreateCommand)
	if err != nil {
		return nil, fmt.Errorf("could not set command commandHandler for command %s: %w", CreateCommand, err)
	}

	proj := projector.NewEventHandler(
		NewProjector(), repo)
	proj.SetEntityFactory(func() eh.Entity { return &Account{} })
	eventBus.AddHandler(eh.MatchAnyEventOf(
		CreatedEvent,
		CreateFailedEvent,
	), proj)

	eh.RegisterCommand(func() eh.Command { return &Create{} })

	eh.RegisterAggregate(func(id uuid.UUID) eh.Aggregate {
		return NewAggregate(id)
	})

	eh.RegisterEventData(CreatedEvent, func() eh.EventData {
		return &CreatedData{}
	})

	return svc, nil
}

func (s *Service) CreateAccount(ctx context.Context, req *CreateAccountRequest) error {
	cmd, err := eh.CreateCommand(CreateCommand)
	if err != nil {
		return fmt.Errorf("create command error: %w", err)
	}

	ccmd, ok := cmd.(*Create)
	if !ok {
		return fmt.Errorf("error command typecast")
	}

	ccmd.ID = req.ID
	ccmd.Name = req.Name
	ccmd.Email = req.Email

	err = s.commandHandler.HandleCommand(ctx, ccmd)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetAccount(ctx context.Context, id uuid.UUID) (*Account, error) {
	var (
		data interface{}
		err  error
	)

	if data, err = s.repo.Find(ctx, id); err != nil {
		if rrErr, ok := err.(eh.RepoError); ok && rrErr.Err == eh.ErrEntityNotFound {
			return nil, eh.ErrEntityNotFound
		}

		return nil, errors.New("unexpected error")
	}

	return data.(*Account), nil
}
