package tranfser

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ravlio/wallet/account"
	"github.com/ravlio/wallet/balance"
	eb "github.com/ravlio/wallet/pkg/event_bus"
	"github.com/ravlio/wallet/pkg/money"
	"github.com/ravlio/wallet/types"
)

type Status uint32

const (
	StatusProcessing Status = iota + 1
	StatusSuccess
	StatusFailure
)

type Transfer struct {
	ID        uint32
	FromID    uint32
	ToID      uint32
	Currency  money.Currency
	Amount    money.Money
	Status    Status
	CreatedAt *time.Time
}
type Client interface {
	Transfer(ctx context.Context, from, to uint32, currency money.Currency, amount money.Money) (*Transfer, error)
	ListTransfers(ctx context.Context) ([]*Transfer, error)
}

type Service struct {
	repo      Repository
	eb        eb.Broker
	accountCl account.Client
}

func NewService(repo Repository, eb eb.Broker) *Service {
	return &Service{repo: repo, eb: eb}
}

/*
transfer money with distributed transaction via saga
1. generate transaction id
2. send credit event
3. wait for response
4. send debit event
5. wait for response
6. send success transaction event
*/

func (s *Service) Transfer(ctx context.Context, fromID, toID uint32,
	currency money.Currency, amount money.Money) (ret *Transfer, err error) {
	if err = s.validateTransfer(ctx, fromID, toID, currency, amount); err != nil {
		return nil, err
	}

	// blocking channel for sequences
	ch := make(chan error)

	// timeout mechanism
	timeout := time.NewTimer(time.Second / 2)
	go func() {
		<-timeout.C
		timeout.Stop()
		ch <- errors.New("tx timeout")
	}()

	// create transaction id
	tx := uuid.New()
	md := map[string]interface{}{"tx": tx}

	defer func() {
		// publish failed event if something went wrong
		if err != nil {
			err = s.eb.Publish(types.TransferFailEvent, nil, md)
		}
	}()

	var sub1, sub2, sub3, sub4 eb.Subscription

	// subscribe to all events first

	sub1, err = s.eb.Subscribe(types.BalanceDebitedEvent, func(msg *eb.Message) {
		if msg.Metadata["tx"] == tx { // check if event belongs to current transaction
			ch <- nil // send signal to channel
		}
	})

	if err != nil {
		return nil, err
	}

	defer sub1.Unsubscribe()

	sub2, err = s.eb.Subscribe(types.BalanceDebitFailEvent, func(msg *eb.Message) {
		if msg.Metadata["tx"] == tx {
			// debit failed
			ch <- msg.Payload.(*balance.DebitFailEvent).Error
		}
	})

	if err != nil {
		return nil, err
	}

	defer sub2.Unsubscribe()

	sub3, err = s.eb.Subscribe(types.BalanceCreditedEvent, func(msg *eb.Message) {
		if msg.Metadata["tx"] == tx {
			ch <- nil
		}
	})

	if err != nil {
		return nil, err
	}

	defer sub3.Unsubscribe()

	sub4, err = s.eb.Subscribe(types.BalanceCreditFailEvent, func(msg *eb.Message) {
		if msg.Metadata["tx"] == tx {
			ch <- msg.Payload.(*balance.DebitFailEvent).Error
		}
	})

	if err != nil {
		return nil, err
	}

	defer sub4.Unsubscribe()

	err = s.eb.Publish(types.BalanceDebitEvent, &balance.DebitEvent{Event: balance.Event{
		AccountID: fromID,
		Currency:  currency,
		Amount:    amount,
	},
	}, md)

	if err != nil {
		return nil, err
	}

	// wait for debit
	if err = <-ch; err != nil {
		return nil, err
	}

	err = s.eb.Publish(types.BalanceCreditEvent, &balance.CreditEvent{Event: balance.Event{
		AccountID: toID,
		Currency:  currency,
		Amount:    amount,
	}}, md)

	if err != nil {
		return nil, err
	}

	// wait for credit
	if err = <-ch; err != nil {
		return nil, err
	}

	transfer := &Transfer{
		FromID:    fromID,
		ToID:      toID,
		Currency:  currency,
		Amount:    amount,
		Status:    StatusSuccess,
		CreatedAt: nil,
	}
	ret, err = s.repo.CreateTransfer(ctx, transfer)
	if err != nil {
		return nil, err
	}

	_ = s.eb.Publish(types.TransferSuccessEvent, nil, md)

	return ret, nil

}

func (s *Service) ListTransfers(ctx context.Context) ([]*Transfer, error) {
	return s.repo.ListTransfers(ctx)
}

func (s *Service) Start() error {
	return nil
}

func (s *Service) Stop() error {
	return nil
}
