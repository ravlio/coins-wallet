package balance

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/ravlio/wallet/pkg/db"
)

func NewMigrations(conn *pgx.Conn, name string) (*db.Migrations, error) {
	m := db.NewMigrations(conn, name)

	err := m.AddUp(1, "create enums", func(ctx context.Context, conn *pgx.Conn) error {
		_, err := conn.Exec(ctx, `CREATE TYPE balance_type AS ENUM ('debit', 'credit');`)

		return err
	})

	if err != nil {
		return nil, err
	}

	err = m.AddUp(2, "create table", func(ctx context.Context, conn *pgx.Conn) error {
		_, err := conn.Exec(ctx, `create table balances
(
    id         serial PRIMARY KEY,
    account_id int4,
    currency   int2,
    amount     int8,
    type       balance_type,
    created_at timestamp
);`)

		return err
	})

	if err != nil {
		return nil, err
	}

	err = m.AddUp(3, "create index", func(ctx context.Context, conn *pgx.Conn) error {
		_, err := conn.Exec(ctx, `CREATE INDEX idx_balance_account_id_currency
    ON balances (account_id, currency);`)

		return err
	})

	if err != nil {
		return nil, err
	}

	err = m.AddDown(1, "drop table", func(ctx context.Context, conn *pgx.Conn) error {
		_, err := conn.Exec(ctx, `drop table balance`)

		return err
	})

	if err != nil {
		return nil, err
	}

	err = m.AddDown(2, "drop enums", func(ctx context.Context, conn *pgx.Conn) error {
		_, err := conn.Exec(ctx, `DROP TYPE balance_type;`)

		return err
	})

	if err != nil {
		return nil, err
	}

	err = m.AddDown(3, "drop index", func(ctx context.Context, conn *pgx.Conn) error {
		_, err := conn.Exec(ctx, `drop index idx_balance_account_id_currency;`)

		return err
	})

	if err != nil {
		return nil, err
	}

	return m, nil

}
