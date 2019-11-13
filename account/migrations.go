package account

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/ravlio/wallet/pkg/db"
)

func NewMigrations(conn *pgx.Conn, name string) (*db.Migrations, error) {
	m := db.NewMigrations(conn, name)

	err := m.AddUp(1, "create table", func(ctx context.Context, conn *pgx.Conn) error {
		_, err := conn.Exec(ctx, `create table accounts (
   							 id    serial primary key,
							name  text,
    						email text
						)`)

		return err
	})

	if err != nil {
		return nil, err
	}

	err = m.AddDown(1, "drop table", func(ctx context.Context, conn *pgx.Conn) error {
		_, err := conn.Exec(ctx, `drop table accounts`)

		return err
	})

	return m, nil

}
