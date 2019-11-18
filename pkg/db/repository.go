package db

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/ravlio/wallet/pkg/errutil"
)

type SelectOptions struct {
	Cols []string
}

type InsertReturningOptions struct {
	WriteCols   []string
	WriteValues map[string]interface{}
	ReadCols    []string
}

type UpdateReturningOptions struct {
	WriteCols   []string
	WriteValues map[string]interface{}
	ReadCols    []string
}

type Repository interface {
	Select(ctx context.Context, model Model, key string, value interface{}, opts *SelectOptions) error
	InsertReturning(ctx context.Context, writeModel, readModel Model, opts *InsertReturningOptions) error
	UpdateReturning(ctx context.Context, writeModel, readModel Model, opts *UpdateReturningOptions) error
	Delete(ctx context.Context, key string, value interface{}) error
}

type PGRepository struct {
	conn  *pgx.Conn
	table string
}

func NewRepository(conn *pgx.Conn, table string) *PGRepository {
	return &PGRepository{conn: conn, table: table}
}

func (r *PGRepository) Select(ctx context.Context, model Model, key string, value interface{}, opts *SelectOptions) error {
	var query string

	var sel string
	if opts.Cols == nil {
		sel = "*"
	} else {
		sel = strings.Join(opts.Cols, ", ")
	}

	query = fmt.Sprintf(`SELECT %s FROM %s WHERE %s=$1`, sel, r.table, key)

	row := r.conn.QueryRow(
		ctx, query, value)

	err := row.Scan(model.Values(opts.Cols...)...)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errutil.ErrNotFound
		}

		return fmt.Errorf("scan error: %w", errutil.NewInternalServerError(err))

	}

	return nil
}

func (r *PGRepository) InsertReturning(ctx context.Context, writeModel, readModel Model, opts *InsertReturningOptions) error {
	var finalWriteCols = make([]string, 0, len(opts.WriteCols)+len(opts.WriteValues))

	finalWriteCols = append(finalWriteCols, opts.WriteCols...)

	var writeVals = make([]interface{}, 0, len(opts.WriteValues))

	for k, v := range opts.WriteValues {
		finalWriteCols = append(finalWriteCols, k)
		writeVals = append(writeVals, v)
	}

	var returning string

	if opts.ReadCols == nil {
		returning = "*"
	} else {
		returning = strings.Join(opts.ReadCols, ", ")
	}

	placeholders := make([]string, 0, len(finalWriteCols))
	for k := range finalWriteCols {
		placeholders = append(placeholders, "$"+strconv.Itoa(k+1))
	}

	query := fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES(%s) RETURNING %s`,
		r.table,
		strings.Join(placeholders, ", "),
		strings.Join(finalWriteCols, ", "),
		returning,
	)

	row := r.conn.QueryRow(
		ctx,
		query,
		append(readModel.Values(opts.ReadCols...), writeVals...)...,
	)

	err := row.Scan(readModel.Values(opts.ReadCols...)...)

	if err != nil {
		return fmt.Errorf("scan error: %w", errutil.NewInternalServerError(err))
	}

	return nil
}

func (r *PGRepository) UpdateReturning(ctx context.Context, writeModel, readModel Model, opts *UpdateReturningOptions) error {
	var finalWriteCols = make([]string, 0, len(opts.WriteCols)+len(opts.WriteValues))

	finalWriteCols = append(finalWriteCols, opts.WriteCols...)

	var writeVals = make([]interface{}, 0, len(opts.WriteValues))

	for k, v := range opts.WriteValues {
		finalWriteCols = append(finalWriteCols, k)
		writeVals = append(writeVals, v)
	}

	var returning string

	if opts.ReadCols == nil {
		returning = "*"
	} else {
		returning = strings.Join(opts.ReadCols, ", ")
	}

	sets := make([]string, 0, len(finalWriteCols))
	for k, v := range finalWriteCols {
		sets = append(sets, v+"=$"+strconv.Itoa(k+1))
	}

	query := fmt.Sprintf(
		`UPDATE %s SET (%s) RETURNING %s`,
		r.table,
		strings.Join(finalWriteCols, ", "),
		returning,
	)

	row := r.conn.QueryRow(
		ctx,
		query,
		append(readModel.Values(opts.ReadCols...), writeVals...)...,
	)

	err := row.Scan(readModel.Values(opts.ReadCols...)...)

	if err != nil {
		return fmt.Errorf("scan error: %w", errutil.NewInternalServerError(err))
	}

	return nil
}

func (r *PGRepository) Delete(ctx context.Context, key string, value interface{}) error {
	query := fmt.Sprintf("SELECT count(*) FROM (DELETE FROM %s WHERE %s=$1 RETURNING %s)", r.table, key, key)
	row := r.conn.QueryRow(ctx, query, value)

	var count int

	err := row.Scan(&count)

	if err != nil {
		return fmt.Errorf("scan error: %w", errutil.NewInternalServerError(err))
	}

	if count == 0 {
		return errutil.ErrNotFound
	}

	return nil
}
