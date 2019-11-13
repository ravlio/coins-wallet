package balance_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/ravlio/wallet/account"
	"github.com/ravlio/wallet/balance"
	"github.com/ravlio/wallet/pkg/errutil"
	"github.com/ravlio/wallet/pkg/event_bus"
	"github.com/ravlio/wallet/pkg/money"
	"github.com/stretchr/testify/require"
)

type accountMock struct {
	CreateAccountHnd func(ctx context.Context, req *account.Account) (*account.Account, error)
	GetAccountHnd    func(ctx context.Context, id uint32) (*account.Account, error)
	DeleteAccountHnd func(ctx context.Context, id uint32) error
	UpdateAccountHnd func(ctx context.Context, req *account.Account) (*account.Account, error)
	ListAccountsHnd  func(ctx context.Context) ([]*account.Account, error)
}

func (a *accountMock) CreateAccount(ctx context.Context, req *account.Account) (*account.Account, error) {
	return a.CreateAccountHnd(ctx, req)
}

func (a *accountMock) GetAccount(ctx context.Context, id uint32) (*account.Account, error) {
	return a.GetAccountHnd(ctx, id)
}

func (a *accountMock) DeleteAccount(ctx context.Context, id uint32) error {
	return a.DeleteAccountHnd(ctx, id)
}

func (a *accountMock) UpdateAccount(ctx context.Context, req *account.Account) (*account.Account, error) {
	return a.UpdateAccountHnd(ctx, req)
}

func (a *accountMock) ListAccounts(ctx context.Context) ([]*account.Account, error) {
	return a.ListAccountsHnd(ctx)
}

var _ account.Client = &accountMock{}

func TestService(t *testing.T) {
	conn, err := pgx.Connect(context.Background(), "postgres://test:@localhost:5432/wallet")
	require.NoError(t, err)

	_, err = conn.Exec(context.Background(), `TRUNCATE balances`)
	require.NoError(t, err)

	repo := balance.NewRepository(conn)

	accs := &accountMock{
		GetAccountHnd: func(ctx context.Context, id uint32) (i *account.Account, e error) {
			if id == 3 {
				return nil, errutil.ErrNotFound
			}

			return &account.Account{ID: id}, nil
		},
	}

	eb := event_bus.NewLocalBroker()
	svc, err := balance.NewService(repo, accs, eb)

	require.NoError(t, err)

	t.Run("non-existing account", func(t *testing.T) {
		_, err := svc.GetBalance(context.Background(), 3, 0)
		require.Equal(t, "accountId", err.(*errutil.BadRequestError).Field)
	})

	t.Run("non-existing currency", func(t *testing.T) {
		_, err = svc.GetBalance(context.Background(), 1, 0)
		require.Equal(t, "currency", err.(*errutil.BadRequestError).Field)
	})

	t.Run("account 1 get zero balance", func(t *testing.T) {
		b, err := svc.GetBalance(context.Background(), 1, money.CurrencyUSD)
		require.NoError(t, err)

		require.Equal(t, uint32(1), b.AccountID)
		require.Equal(t, money.CurrencyUSD, b.Currency)
		require.Equal(t, money.Money(0), b.Balance)
	})

	t.Run("account 1 credit negative amount should fail", func(t *testing.T) {
		err := svc.Credit(context.Background(), 1, money.CurrencyUSD, money.Float64(-10))
		require.Error(t, err)
	})

	t.Run("account 1 debit negative balance should fail", func(t *testing.T) {
		err := svc.Debit(context.Background(), 1, money.CurrencyUSD, money.Float64(10))
		require.Error(t, err)
	})

	t.Run("account 1 increase balance to 0.1 usd", func(t *testing.T) {
		err := svc.Credit(context.Background(), 1, money.CurrencyUSD, money.Float64(0.1))
		require.NoError(t, err)
	})

	t.Run("account 1 increase balance to 0.5 usd", func(t *testing.T) {
		err := svc.Credit(context.Background(), 1, money.CurrencyUSD, money.Float64(0.5))
		require.NoError(t, err)
	})

	t.Run("account 1 increase balance to 2.3 rub", func(t *testing.T) {
		err := svc.Credit(context.Background(), 1, money.CurrencyRUB, money.Float64(2.3))
		require.NoError(t, err)
	})

	t.Run("account 1 decrease balance to 0.65 rub", func(t *testing.T) {
		err := svc.Debit(context.Background(), 1, money.CurrencyRUB, money.Float64(0.65))
		require.NoError(t, err)
	})

	t.Run("account 2 increase balance to 121.12345 rub", func(t *testing.T) {
		err := svc.Credit(context.Background(), 2, money.CurrencyRUB, money.Float64(121.12345))
		require.NoError(t, err)
	})

	t.Run("account 1 list balances", func(t *testing.T) {
		l, err := svc.ListBalances(context.Background(), 1)
		require.NoError(t, err)

		exp := []*balance.Balance{
			{
				AccountID: 1,
				Currency:  money.CurrencyRUB,
				Balance:   money.Float64(1.65),
			},
			{
				AccountID: 1,
				Currency:  money.CurrencyUSD,
				Balance:   money.Float64(0.6),
			},
		}

		require.Equal(t, exp, l)
	})

	t.Run("account 2 list balances", func(t *testing.T) {
		l, err := svc.ListBalances(context.Background(), 2)
		require.NoError(t, err)

		exp := []*balance.Balance{
			{
				AccountID: 2,
				Currency:  money.CurrencyRUB,
				Balance:   money.Float64(121.12345),
			},
		}

		require.Equal(t, exp, l)
	})

}
