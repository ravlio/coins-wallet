package tranfser

import (
	"context"

	"github.com/ravlio/wallet/pkg/errutil"
	"github.com/ravlio/wallet/pkg/money"
)

func (s *Service) validateTransfer(ctx context.Context, fromID, toID uint32, currency money.Currency, amount money.Money) error {
	if !money.IsCurrency(currency) {
		return errutil.NewBadRequestFieldError("currency", "non-existent currency")
	}

	if amount < 0 {
		return errutil.NewBadRequestFieldError("amount", "negative amount")
	}

	if fromID < 1 {
		return errutil.NewBadRequestFieldError("fromId", "should be more than 1")
	}
	if toID < 1 {
		return errutil.NewBadRequestFieldError("toId", "should be more than 1")
	}

	return nil
}
