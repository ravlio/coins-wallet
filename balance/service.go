package balance

import (
	"context"
	"errors"
	"time"

	"github.com/ravlio/wallet/account"
	"github.com/ravlio/wallet/pkg/db"
	eb "github.com/ravlio/wallet/pkg/event_bus"
	"github.com/ravlio/wallet/pkg/money"
	"github.com/ravlio/wallet/types"
)

type Balance struct {
	AccountID uint32
	Currency  money.Currency
	Balance   money.Money
	CreatedAt *time.Time
}

type Client interface {
	Credit(ctx context.Context, accountID uint32, currency money.Currency, amount money.Money) error
	Debit(ctx context.Context, accountID uint32, currency money.Currency, amount money.Money) error
	GetBalance(ctx context.Context, accountID uint32, currency money.Currency) (*Balance, error)
	ListBalances(ctx context.Context, accountID uint32) ([]*Balance, error)
}

type Service struct {
	repo      Repository
	eb        eb.Broker
	accountCl account.Client
	subs      []eb.Subscription
}

func NewService(repo Repository, account account.Client, eb eb.Broker) (*Service, error) {
	svc := &Service{repo: repo, eb: eb, accountCl: account}
	err := svc.subscribe()
	if err != nil {
		return nil, err
	}

	return svc, nil
}

func (s *Service) Credit(ctx context.Context, accountID uint32, currency money.Currency, amount money.Money) error {
	if err := s.validateCredit(ctx, accountID, currency, amount); err != nil {
		return err
	}

	return s.repo.Credit(ctx, accountID, currency, amount)
}

func (s *Service) Debit(ctx context.Context, accountID uint32, currency money.Currency, amount money.Money) error {
	if err := s.validateDebit(ctx, accountID, currency, amount); err != nil {
		return err
	}

	return s.repo.Debit(ctx, accountID, currency, amount)
}

func (s *Service) GetBalance(ctx context.Context, accountID uint32, currency money.Currency) (*Balance, error) {
	if err := s.validateCommon(ctx, accountID, currency); err != nil {
		return nil, err
	}

	return s.repo.GetBalance(ctx, accountID, currency)
}

func (s *Service) ListBalances(ctx context.Context, accountID uint32) ([]*Balance, error) {
	return s.repo.ListBalances(ctx, accountID)
}

func (s *Service) subscribe() error {
	var sub1, sub2 eb.Subscription
	var err error

	sub1, err = s.eb.Subscribe(types.BalanceCreditEvent, s.handleBalanceCreditEvent)

	if err != nil {
		return err
	}

	s.subs = append(s.subs, sub1)

	sub2, err = s.eb.Subscribe(types.BalanceDebitEvent, s.handleBalanceDebitEvent)

	if err != nil {
		return err
	}

	s.subs = append(s.subs, sub2)

	return nil
}

func (s *Service) handleBalanceCreditEvent(msg *eb.Message) {
	ctx := context.Background()

	ev := msg.Payload.(*CreditEvent)

	if err := s.validateCredit(ctx, ev.AccountID, ev.Currency, ev.Amount); err != nil {
		_ = s.eb.Publish(types.BalanceCreditFailEvent, nil, msg.Metadata)
		return
	}

	tx, err := s.repo.CreditTx(ctx, ev.AccountID, ev.Currency, ev.Amount)
	if err != nil {
		_ = s.eb.Publish(types.BalanceCreditFailEvent, nil, msg.Metadata)
		return
	}

	err = s.handleSaga(ctx, tx, msg)
	if err != nil {
		// TODO: ugly
		if err.Error() != "transfer fail" {
			_ = s.eb.Publish(types.BalanceCreditFailEvent, nil, msg.Metadata)
		}
	}
}

func (s *Service) handleBalanceDebitEvent(msg *eb.Message) {
	ctx := context.Background()

	ev := msg.Payload.(*DebitEvent)

	if err := s.validateDebit(ctx, ev.AccountID, ev.Currency, ev.Amount); err != nil {
		_ = s.eb.Publish(types.BalanceDebitFailEvent, nil, msg.Metadata)
		return
	}

	tx, err := s.repo.CreditTx(ctx, ev.AccountID, ev.Currency, ev.Amount)
	if err != nil {
		_ = s.eb.Publish(types.BalanceDebitFailEvent, nil, msg.Metadata)
		return
	}

	err = s.handleSaga(ctx, tx, msg)

	if err != nil {
		// TODO: ugly
		if err.Error() != "transfer fail" {
			_ = s.eb.Publish(types.BalanceDebitFailEvent, nil, msg.Metadata)
		}
	}
}

func (s *Service) handleSaga(ctx context.Context, tx db.Tx, msg *eb.Message) error {
	var err error

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	// blocking channel for sequences
	ch := make(chan error)

	// timeout mechanism
	timeout := time.NewTimer(time.Second / 2)
	go func() {
		<-timeout.C
		timeout.Stop()
		ch <- errors.New("tx timeout")
	}()

	sub1, _ := s.eb.Subscribe(types.TransferSuccessEvent, func(m *eb.Message) {
		if m.Metadata["tx"] == msg.Metadata["tx"] {
			ch <- nil
		}
	})

	defer sub1.Unsubscribe()

	sub2, _ := s.eb.Subscribe(types.TransferFailEvent, func(m *eb.Message) {
		if m.Metadata["tx"] == msg.Metadata["tx"] {
			ch <- errors.New("transfer fail")
		}
	})

	defer sub2.Unsubscribe()

	err = <-ch

	if err != nil {
		return err
	}

	err = tx.Commit(ctx)

	if err != nil {
		return err
	}

	return nil
}
