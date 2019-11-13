package tranfser

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/ravlio/wallet/pkg/errutil"
)

type Repository interface {
	CreateTransfer(ctx context.Context, req *Transfer) (*Transfer, error)
	ListTransfers(ctx context.Context) ([]*Transfer, error)
}

type PostgresRepository struct {
	conn *pgx.Conn
}

func NewRepository(conn *pgx.Conn) *PostgresRepository {
	return &PostgresRepository{conn: conn}
}

func (p *PostgresRepository) CreateTransfer(ctx context.Context, req *Transfer) (*Transfer, error) {
	writeModel := &model{}

	if err := writeModel.FromEntity(req); err != nil {
		return nil, err
	}

	row := p.conn.QueryRow(
		ctx,
		`INSERT INTO transfers 
		(from_id, to_id,currency,amount,status,created_at) 
		VALUES($1, $2, $3, $4, $5) 
		RETURNING id, from_id,to_id,currency,amount,status,created_at`,
		writeModel.Values(columnFromID, columnToID, columnCurrency, columnAmount, columnStatus)...,
	)

	readModel := &model{}

	err := row.Scan(readModel.Values(columnID, columnFromID, columnToID,
		columnCurrency, columnAmount, columnStatus, columnCreatedAt)...)

	if err != nil {
		return nil, fmt.Errorf("scan error: %w", errutil.NewInternalServerError(err))
	}

	return readModel.ToEntity(), nil

}

func (p *PostgresRepository) ListTransfers(ctx context.Context) ([]*Transfer, error) {
	rows, err := p.conn.Query(ctx, `SELECT id, from_id,to_id,currency,amount,status,created_a FROM transfers`)
	if err != nil {
		return nil, fmt.Errorf("list transfers error: %w", errutil.NewInternalServerError(err))
	}

	ret := make([]*Transfer, 0)

	for rows.Next() {
		readModel := &model{}
		if err = rows.Scan(readModel.Values(columnID, columnFromID, columnToID,
			columnCurrency, columnAmount, columnStatus, columnCreatedAt)...); err != nil {
			return nil, fmt.Errorf("scan error: %w", errutil.NewInternalServerError(err))
		}

		ret = append(ret, readModel.ToEntity())
	}

	return ret, nil
}
