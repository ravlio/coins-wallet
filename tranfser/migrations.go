package tranfser

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/ravlio/wallet/pkg/db"
)

func NewMigrations(conn *pgx.Conn, name string) (*db.Migrations, error) {
	m := db.NewMigrations(conn, name)

	err := m.AddUp(1, "create all", func(ctx context.Context, conn *pgx.Conn) error {
		tx, err := conn.Begin(ctx)
		if err != nil {
			return err
		}

		defer tx.Rollback(ctx)

		_, err = tx.Exec(ctx, `CREATE TYPE transfer_status AS ENUM ('processing','success', 'failure');`)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx,
			`
create table transfers
		(
			id         serial PRIMARY KEY,
			from_id    int4,
			to_id      int4,
			currency   int2,
			amount     int8,
			status     transfer_status,
			created_at timestamp
		)`)

		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, `CREATE INDEX idx_transfers_from_id_to_id
		ON transfers (from_id, to_id);`)

		if err != nil {
			return err
		}

		return tx.Commit(ctx)
	})

	if err != nil {
		return nil, err
	}

	err = m.AddDown(1, "drop all", func(ctx context.Context, conn *pgx.Conn) error {
		_, err := conn.Exec(ctx,
			`drop  index idx_transfers_from_id_to_id; 
					drop table transfers;
				drop type transfer_status;`)
		return err
	})

	if err != nil {
		return nil, err
	}

	return m, nil

}
