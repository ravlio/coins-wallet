package balance

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/ravlio/wallet/pkg/db"
	"github.com/ravlio/wallet/pkg/errutil"
	"github.com/ravlio/wallet/pkg/money"
)

type Repository interface {
	Credit(ctx context.Context, accountID uint32, currency money.Currency, amount money.Money) error
	CreditTx(ctx context.Context, accountID uint32, currency money.Currency, amount money.Money) (db.Tx, error)
	Debit(ctx context.Context, accountID uint32, currency money.Currency, amount money.Money) error
	DebitTx(ctx context.Context, accountID uint32, currency money.Currency, amount money.Money) (db.Tx, error)
	GetBalance(ctx context.Context, accountID uint32, currency money.Currency) (*Balance, error)
	ListBalances(ctx context.Context, accountID uint32) ([]*Balance, error)
}

type PostgresRepository struct {
	conn *pgx.Conn
}

func NewRepository(conn *pgx.Conn) *PostgresRepository {
	return &PostgresRepository{conn: conn}
}

func (p *PostgresRepository) Credit(ctx context.Context, accountID uint32, currency money.Currency, amount money.Money) error {
	tx, err := p.credit(ctx, accountID, currency, amount)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (p *PostgresRepository) CreditTx(ctx context.Context, accountID uint32, currency money.Currency, amount money.Money) (db.Tx, error) {
	tx, err := p.credit(ctx, accountID, currency, amount)
	if err != nil {
		return nil, err
	}

	return db.NewTx(tx), nil
}

func (p *PostgresRepository) credit(ctx context.Context, accountID uint32, currency money.Currency, amount money.Money) (tx pgx.Tx, err error) {
	tx, err = p.conn.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	_, err = tx.Exec(
		ctx,
		`INSERT INTO balances (account_id, currency, amount, type, created_at) VALUES($1, $2, $3, $4, NOW())`,
		accountID, currency, amount, opTypeCredit,
	)

	if err != nil {
		return tx, fmt.Errorf("scan error: %w", errutil.NewInternalServerError(err))
	}

	return tx, nil
}

func (p *PostgresRepository) Debit(ctx context.Context, accountID uint32, currency money.Currency, amount money.Money) error {
	tx, err := p.debit(ctx, accountID, currency, amount)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (p *PostgresRepository) DebitTx(ctx context.Context, accountID uint32, currency money.Currency, amount money.Money) (db.Tx, error) {
	tx, err := p.debit(ctx, accountID, currency, amount)
	if err != nil {
		return nil, err
	}

	return db.NewTx(tx), nil
}

func (p *PostgresRepository) debit(ctx context.Context, accountID uint32, currency money.Currency, amount money.Money) (tx pgx.Tx, err error) {
	tx, err = p.conn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction error: %w", errutil.NewInternalServerError(err))
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	_, err = tx.Exec(
		ctx,
		`INSERT INTO balances (account_id, currency, amount, type, created_at) VALUES($1, $2, $3, $4, NOW())`,
		accountID, currency, amount, opTypeDebit,
	)

	if err != nil {
		return nil, fmt.Errorf("scan error: %w", errutil.NewInternalServerError(err))
	}

	balance, err := p.GetBalance(ctx, accountID, currency)
	if err != nil {
		return nil, err
	}

	if balance.Balance < 0 {
		return nil, errors.New("balance can't be below 0")
	}

	err = tx.Commit(ctx)

	if err != nil {
		return nil, fmt.Errorf("commit transaction error: %w", errutil.NewInternalServerError(err))
	}

	return tx, nil
}
func (p *PostgresRepository) GetBalance(ctx context.Context, accountID uint32, currency money.Currency) (*Balance, error) {
	ret := &Balance{AccountID: accountID, Currency: currency}

	row := p.conn.QueryRow(
		ctx,
		`SELECT COALESCE(SUM(CASE WHEN type='credit' THEN amount ELSE -amount END),0) FROM balances WHERE account_id=$1 AND currency=$2`,
		accountID,
		currency,
	)

	err := row.Scan(&ret.Balance)
	if err != nil {
		return nil, fmt.Errorf("scan error: %w", errutil.NewInternalServerError(err))
	}

	return ret, nil
}

func (p *PostgresRepository) ListBalances(ctx context.Context, accountID uint32) ([]*Balance, error) {
	rows, err := p.conn.Query(
		ctx,
		`SELECT 
				currency, 
				COALESCE(SUM(CASE WHEN type='credit' THEN amount ELSE -amount END),0) 
			FROM balances WHERE account_id=$1 GROUP BY currency`,
		accountID,
	)

	if err != nil {
		return nil, fmt.Errorf("list balances error: %w", errutil.NewInternalServerError(err))
	}

	ret := make([]*Balance, 0)

	for rows.Next() {
		b := &Balance{AccountID: accountID}
		if err = rows.Scan(&b.Currency, &b.Balance); err != nil {
			return nil, fmt.Errorf("scan error: %w", errutil.NewInternalServerError(err))
		}
		ret = append(ret, b)
	}

	return ret, nil
}
