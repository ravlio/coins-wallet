package db

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/require"
)

func TestMigrations(t *testing.T) {
	conn, err := pgx.Connect(context.Background(), "postgres://test:@localhost:5432/wallet")
	require.NoError(t, err)

	_, err = conn.Exec(context.Background(), `TRUNCATE migrations`)

	m1 := NewMigrations(conn, "svc1")

	db := map[string]int{}

	err = m1.AddUp(2, "migration 2", func(ctx context.Context, conn *pgx.Conn) error {
		db["tk"] = 1
		return nil
	})

	require.EqualError(t, err, "id sequence error upon adding migration. Expected 1, got 2")

	err = m1.AddUp(1, "migration 1", func(ctx context.Context, conn *pgx.Conn) error {
		db["m1"] = 1
		return nil
	})

	require.NoError(t, err)

	err = m1.AddUp(3, "migration 3", func(ctx context.Context, conn *pgx.Conn) error {
		db["m2"] = 1
		return nil
	})

	require.EqualError(t, err, "id sequence error upon adding migration. Expected 2, got 3")

	err = m1.AddUp(2, "migration 2", func(ctx context.Context, conn *pgx.Conn) error {
		db["m2"] = 2
		return nil
	})

	require.NoError(t, err)

	c, err := m1.Up()

	require.NoError(t, err)
	require.EqualValues(t, 2, c)

	err = m1.AddUp(3, "migration 3", func(ctx context.Context, conn *pgx.Conn) error {
		db["m3"] = 3
		return nil
	})

	require.NoError(t, err)

	c, err = m1.Up()

	require.NoError(t, err)
	require.EqualValues(t, 1, c)

	err = m1.AddUp(4, "migration 4 failing", func(ctx context.Context, conn *pgx.Conn) error {
		return errors.New("fail")
	})

	require.NoError(t, err)

	c, err = m1.Up()

	require.EqualError(t, err, "error while migration \"migration 4 failing\" (id: 4): fail")

	require.EqualValues(t, 1, db["m1"])
	require.EqualValues(t, 2, db["m2"])
	require.EqualValues(t, 3, db["m3"])

	//////////////

	err = m1.AddDown(2, "migration 2 down", func(ctx context.Context, conn *pgx.Conn) error {
		delete(db, "m1")
		return nil
	})

	require.EqualError(t, err, "id sequence error upon adding migration. Expected 1, got 2")

	err = m1.AddDown(1, "migration 1 down", func(ctx context.Context, conn *pgx.Conn) error {
		delete(db, "m1")
		return nil
	})

	require.NoError(t, err)

	err = m1.AddDown(2, "migration 2 down", func(ctx context.Context, conn *pgx.Conn) error {
		delete(db, "m2")
		return nil
	})

	require.NoError(t, err)

	err = m1.AddDown(3, "migration 3 down", func(ctx context.Context, conn *pgx.Conn) error {
		delete(db, "m3")
		return nil
	})

	require.NoError(t, err)

	err = m1.Down()
	require.NoError(t, err)

	_, ok := db["m1"]
	require.False(t, ok)
	_, ok = db["m2"]
	require.False(t, ok)
	_, ok = db["m3"]
	require.False(t, ok)

}
