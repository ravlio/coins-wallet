package account

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/ravlio/wallet/pkg/errutil"
)

type Repository interface {
	GetAccountByID(ctx context.Context, id uint32) (*Account, error)
	CreateAccount(ctx context.Context, acc *Account) (*Account, error)
	UpdateAccount(ctx context.Context, acc *Account) (*Account, error)
	DeleteAccount(ctx context.Context, id uint32) error
	ListAccounts(ctx context.Context) ([]*Account, error)
}

type PostgresRepository struct {
	conn *pgx.Conn
}

func NewRepository(conn *pgx.Conn) *PostgresRepository {
	return &PostgresRepository{conn: conn}
}

func (p *PostgresRepository) GetAccountByID(ctx context.Context, id uint32) (*Account, error) {
	m := &model{}

	row := p.conn.QueryRow(
		ctx, `SELECT id, name, email FROM accounts WHERE id=$1`, id)
	err := row.Scan(m.Values(columnID, columnName, columnEmail)...)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errutil.ErrNotFound
		}

		return nil, fmt.Errorf("scan error: %w", errutil.NewInternalServerError(err))

	}

	return m.ToEntity(), nil

}

func (p *PostgresRepository) CreateAccount(ctx context.Context, acc *Account) (*Account, error) {
	writeModel := &model{}

	if err := writeModel.FromEntity(acc); err != nil {
		return nil, err
	}

	row := p.conn.QueryRow(
		ctx,
		`INSERT INTO accounts (name, email) VALUES($1, $2) RETURNING id, name, email`,
		writeModel.Values(columnName, columnEmail)...,
	)

	readModel := &model{}

	err := row.Scan(readModel.Values(columnID, columnName, columnEmail)...)

	if err != nil {
		return nil, fmt.Errorf("scan error: %w", errutil.NewInternalServerError(err))
	}

	return readModel.ToEntity(), nil
}

func (p *PostgresRepository) UpdateAccount(ctx context.Context, acc *Account) (*Account, error) {
	writeModel := &model{}

	if err := writeModel.FromEntity(acc); err != nil {
		return nil, err
	}

	row := p.conn.QueryRow(
		ctx,
		`UPDATE accounts SET name=$1, email=$2 WHERE id=$3 RETURNING id, name, email`,
		writeModel.Values(columnName, columnEmail, columnID)...,
	)

	readModel := &model{}

	err := row.Scan(readModel.Values(columnID, columnName, columnEmail)...)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errutil.ErrNotFound
		}

		return nil, fmt.Errorf("scan error: %w", errutil.NewInternalServerError(err))
	}

	return readModel.ToEntity(), nil
}

func (p *PostgresRepository) DeleteAccount(ctx context.Context, id uint32) error {
	res, err := p.conn.Exec(ctx, `DELETE FROM accounts WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("delete error: %w", errutil.NewInternalServerError(err))
	}

	if res.RowsAffected() == 0 {
		return errutil.ErrNotFound
	}

	return nil
}

func (p *PostgresRepository) ListAccounts(ctx context.Context) ([]*Account, error) {
	rows, err := p.conn.Query(ctx, `SELECT id, name, email FROM accounts`)
	if err != nil {
		return nil, fmt.Errorf("list accounts error: %w", errutil.NewInternalServerError(err))
	}

	ret := make([]*Account, 0)

	for rows.Next() {
		readModel := &model{}
		if err = rows.Scan(readModel.Values(columnID, columnName, columnEmail)...); err != nil {
			return nil, fmt.Errorf("scan error: %w", errutil.NewInternalServerError(err))
		}

		ret = append(ret, readModel.ToEntity())
	}

	return ret, nil
}
