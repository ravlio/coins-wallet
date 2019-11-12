package tranfser

import (
	"context"

	"github.com/ravlio/wallet/pkg/money"
)

type Service interface {
	Transfer(ctx context.Context, from, to string, amount money.Money) error
}
