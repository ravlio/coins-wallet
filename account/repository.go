package account

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/ravlio/wallet/pkg/db"
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
	repo db.Repository
}

func NewRepository(conn *pgx.Conn) *PostgresRepository {
	return &PostgresRepository{conn: conn, repo: db.NewRepository(conn, "accounts")}
}

func (p *PostgresRepository) GetAccountByID(ctx context.Context, id uint32) (*Account, error) {
	m := &model{}

	err := p.repo.Select(ctx, m, "id", id, &db.SelectOptions{Cols: selectColumns})
	if err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (p *PostgresRepository) CreateAccount(ctx context.Context, acc *Account) (*Account, error) {
	writeModel := &model{}

	if err := writeModel.FromEntity(acc); err != nil {
		return nil, err
	}

	readModel := &model{}

	err := p.repo.InsertReturning(
		ctx,
		writeModel,
		readModel,
		&db.InsertReturningOptions{
			WriteCols:   insertColumns,
			WriteValues: map[string]interface{}{"created_at": "NOW()"},
			ReadCols:    selectColumns,
		},
	)

	if err != nil {
		return nil, err
	}

	return readModel.ToEntity(), nil
}

func (p *PostgresRepository) UpdateAccount(ctx context.Context, acc *Account) (*Account, error) {

	writeModel := &model{}

	if err := writeModel.FromEntity(acc); err != nil {
		return nil, err
	}

	readModel := &model{}

	err := p.repo.UpdateReturning(
		ctx,
		writeModel,
		readModel,
		&db.UpdateReturningOptions{
			WriteCols:   updateColumns,
			WriteValues: map[string]interface{}{"updated_at": "NOW()"},
			ReadCols:    selectColumns,
		},
	)

	if err != nil {
		return nil, err
	}

	return readModel.ToEntity(), nil
}

func (p *PostgresRepository) DeleteAccount(ctx context.Context, id uint32) error {
	return p.repo.Delete(ctx, "id", id)
}

func (p *PostgresRepository) ListAccounts(ctx context.Context) ([]*Account, error) {
	rows, err := p.conn.Query(ctx, `SELECT id, name, email FROM accounts`)
	if err != nil {
		return nil, fmt.Errorf("list accounts error: %w", errutil.NewInternalServerError(err))
	}

	ret := make([]*Account, 0)

	for rows.Next() {
		readModel := &model{}
		if err = rows.Scan(readModel.Values("id", "name", "email")...); err != nil {
			return nil, fmt.Errorf("scan error: %w", errutil.NewInternalServerError(err))
		}

		ret = append(ret, readModel.ToEntity())
	}

	return ret, nil
}
