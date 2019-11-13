package account_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/ravlio/wallet/account"
	"github.com/ravlio/wallet/pkg/errutil"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	conn, err := pgx.Connect(context.Background(), "postgres://test:@localhost:5432/wallet")
	require.NoError(t, err)

	_, err = conn.Exec(context.Background(), `TRUNCATE accounts`)
	require.NoError(t, err)

	repo := account.NewRepository(conn)

	svc := account.NewService(repo)

	var id1, id2 uint32
	var list []*account.Account

	t.Run("Non-existing account getting should error", func(t *testing.T) {
		_, err := svc.GetAccount(context.Background(), 1)
		require.Error(t, err)
		require.True(t, errors.Is(err, errutil.ErrNotFound))
	})

	t.Run("Non-existing account updating should error", func(t *testing.T) {
		acc := &account.Account{ID: 1}
		_, err := svc.UpdateAccount(context.Background(), acc)
		require.Error(t, err)
		require.True(t, errors.Is(err, errutil.ErrNotFound))
	})

	t.Run("Non-existing account deleting should error", func(t *testing.T) {
		err := svc.DeleteAccount(context.Background(), 1)
		require.Error(t, err)
		require.True(t, errors.Is(err, errutil.ErrNotFound))
	})

	t.Run("Create account 1", func(t *testing.T) {
		acc := &account.Account{
			ID:    123, // should be ignored
			Name:  "acc1",
			Email: "acc1@gmail.com",
		}
		ret, err := svc.CreateAccount(context.Background(), acc)
		require.NoError(t, err)
		require.True(t, ret.ID > 0)
		require.Equal(t, acc.Name, ret.Name)
		require.Equal(t, acc.Email, ret.Email)

		id1 = ret.ID

		list = append(list, ret)
	})

	t.Run("Create account 2", func(t *testing.T) {
		acc := &account.Account{
			ID:    123, // should be ignored
			Name:  "acc2",
			Email: "acc2@gmail.com",
		}
		ret, err := svc.CreateAccount(context.Background(), acc)
		require.NoError(t, err)
		require.True(t, ret.ID > 0)
		require.Equal(t, acc.Name, ret.Name)
		require.Equal(t, acc.Email, ret.Email)

		id2 = ret.ID
		list = append(list, ret)
	})

	t.Run("List should return 2 accounts", func(t *testing.T) {
		ret, err := svc.ListAccounts(context.Background())
		require.NoError(t, err)
		require.Equal(t, list, ret)
	})

	t.Run("Update account 1", func(t *testing.T) {
		req := &account.Account{
			ID:    id1,
			Name:  "newname",
			Email: "newemail@gmail.com",
		}
		resp, err := svc.UpdateAccount(context.Background(), req)

		require.NoError(t, err)

		require.Equal(t, req, resp)
	})

	t.Run("Delete account 2", func(t *testing.T) {
		err := svc.DeleteAccount(context.Background(), id2)

		require.NoError(t, err)
	})

	t.Run("List should return 1 account", func(t *testing.T) {
		ret, err := svc.ListAccounts(context.Background())

		eq := []*account.Account{{
			ID:    id1,
			Name:  "newname",
			Email: "newemail@gmail.com",
		}}

		require.NoError(t, err)
		require.Equal(t, eq, ret)
	})
}
