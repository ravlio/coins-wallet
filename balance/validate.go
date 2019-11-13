package balance

import (
	"context"

	"github.com/ravlio/wallet/pkg/errutil"
	"github.com/ravlio/wallet/pkg/money"
)

func (s *Service) validateCommon(ctx context.Context, accountID uint32, currency money.Currency) error {
	if _, err := s.accountCl.GetAccount(ctx, accountID); err != nil {
		return errutil.NewBadRequestFieldError("accountId", err.Error())
	}

	if !money.IsCurrency(currency) {
		return errutil.NewBadRequestFieldError("currency", "non-existent currency")
	}

	return nil
}

func (s *Service) validateCredit(ctx context.Context, accountID uint32, currency money.Currency, amount money.Money) error {
	if err := s.validateCommon(ctx, accountID, currency); err != nil {
		return err
	}

	if amount < 0 {
		return errutil.NewBadRequestFieldError("amount", "negative amount")
	}

	return nil
}

func (s *Service) validateDebit(ctx context.Context, accountID uint32, currency money.Currency, amount money.Money) error {
	if err := s.validateCommon(ctx, accountID, currency); err != nil {
		return err
	}

	if amount < 0 {
		return errutil.NewBadRequestFieldError("amount", "negative amount")
	}

	return nil
}
