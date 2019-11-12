package balance

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
	"github.com/ravlio/wallet/pkg/money"
)

type Client interface {
	Credit(ctx context.Context, req *Credit) error
	Debit(ctx context.Context, req *Debit) error
	GetBalance(ctx context.Context, account string, currency money.Currency) (*Balance, error)
	ListBalances(ctx context.Context, account string) ([]*Balance, error)
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
	err = commandBus.SetHandler(svc.commandHandler, CreditCommand)
	if err != nil {
		return nil, fmt.Errorf("could not set command commandHandler for command %s: %w", CreditCommand, err)
	}
	err = commandBus.SetHandler(svc.commandHandler, DebitCommand)
	if err != nil {
		return nil, fmt.Errorf("could not set command commandHandler for command %s: %w", DebitCommand, err)
	}

	proj := projector.NewEventHandler(
		NewProjector(), repo)
	proj.SetEntityFactory(func() eh.Entity { return &Balance{} })
	eventBus.AddHandler(eh.MatchAnyEventOf(
		CreditedEvent,
		DebitedEvent,
	), proj)

	eh.RegisterCommand(func() eh.Command { return &Credit{} })
	eh.RegisterCommand(func() eh.Command { return &Debit{} })

	eh.RegisterAggregate(func(id uuid.UUID) eh.Aggregate {
		return NewAggregate(id)
	})

	eh.RegisterEventData(CreditedEvent, func() eh.EventData {
		return &CreditedData{}
	})

	eh.RegisterEventData(CreditFailedEvent, func() eh.EventData {
		return &CreditFailedData{}
	})

	eh.RegisterEventData(DebitedEvent, func() eh.EventData {
		return &DebitedData{}
	})

	eh.RegisterEventData(DebitFailedEvent, func() eh.EventData {
		return &DebitFailedData{}
	})

	return svc, nil
}

func (s *Service) Credit(ctx context.Context, req *Credit) error {
	cmd, err := eh.CreateCommand(CreditCommand)
	if err != nil {
		return fmt.Errorf("create command error: %w", err)
	}

	ccmd, ok := cmd.(*Credit)
	if !ok {
		return fmt.Errorf("error command typecast")
	}

	ccmd.ID = req.ID
	ccmd.AccountID = req.AccountID
	ccmd.Currency = req.Currency
	ccmd.Amount = req.Amount

	err = s.commandHandler.HandleCommand(ctx, ccmd)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Debit(ctx context.Context, req *Credit) error {
	cmd, err := eh.CreateCommand(DebitCommand)
	if err != nil {
		return fmt.Errorf("create command error: %w", err)
	}

	ccmd, ok := cmd.(*Debit)
	if !ok {
		return fmt.Errorf("error command typecast")
	}

	ccmd.ID = req.ID
	ccmd.AccountID = req.AccountID
	ccmd.Currency = req.Currency
	ccmd.Amount = req.Amount

	err = s.commandHandler.HandleCommand(ctx, ccmd)
	if err != nil {
		return err
	}

	return nil
}


func (s *Service) GetBalance(ctx context.Context, account string, currency money.Currency) (*Balance, error) {
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

	return data.(*Balance), nil
}
