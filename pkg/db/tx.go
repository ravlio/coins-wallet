package db

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type Tx interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type PgxTx struct {
	tx pgx.Tx
}

func NewTx(tx pgx.Tx) Tx {
	return &PgxTx{tx: tx}
}

func (p *PgxTx) Commit(ctx context.Context) error {
	return p.tx.Commit(ctx)
}

func (p *PgxTx) Rollback(ctx context.Context) error {
	return p.tx.Rollback(ctx)
}
