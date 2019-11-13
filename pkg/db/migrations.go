package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

type Migrations struct {
	up      []Migration
	down    []Migration
	conn    *pgx.Conn
	svcName string // svcName of service for migrations collection
}

type Migration struct {
	id   uint32
	name string
	conn *pgx.Conn
	hnd  MigrationHnd
}
type MigrationHnd func(ctx context.Context, conn *pgx.Conn) error

func NewMigrations(conn *pgx.Conn, svcName string) *Migrations {
	return &Migrations{conn: conn, svcName: svcName}
}

func (m *Migrations) checkSequence(l []Migration, id uint32) error {
	if len(l) == 0 {
		if id != 1 {
			return fmt.Errorf("id sequence error upon adding migration. Expected 1, got %d", id)
		}
		return nil
	}

	last := l[len(l)-1]

	if last.id != id-1 {
		return fmt.Errorf("id sequence error upon adding migration. Expected %d, got %d", last.id+1, id)
	}

	return nil
}

func (m *Migrations) createTable(ctx context.Context) error {
	_, err := m.conn.Exec(ctx,
		`create table if not exists migrations (
    			service VARCHAR UNIQUE NOT NULL,
    			id integer NOT NULL DEFAULT 1
				);`)

	if err != nil {
		return fmt.Errorf("create migation table error: %w", err)
	}

	return nil
}

func (m *Migrations) AddUp(id uint32, name string, hnd MigrationHnd) error {
	err := m.checkSequence(m.up, id)
	if err != nil {
		return err
	}

	m.up = append(m.up, Migration{id: id, name: name, hnd: hnd})
	return nil
}

func (m *Migrations) AddDown(id uint32, name string, hnd MigrationHnd) error {
	err := m.checkSequence(m.down, id)
	if err != nil {
		return err
	}

	m.down = append(m.down, Migration{id: id, name: name, hnd: hnd})
	return nil
}

func (m *Migrations) Up() (count int, err error) {
	ctx := context.Background()
	err = m.createTable(ctx)

	if err != nil {
		return 0, err
	}

	for _, v := range m.up {
		var curid uint32

		row := m.conn.QueryRow(ctx, `SELECT id FROM migrations WHERE service=$1`, m.svcName)
		err = row.Scan(&curid)
		if err != nil && err != pgx.ErrNoRows {
			return 0, fmt.Errorf("get migration id error: %w", err)
		}

		if curid != v.id-1 {
			continue
		}

		err = v.hnd(ctx, v.conn)
		if err != nil {
			return 0, fmt.Errorf("error while migration %q (id: %d): %w", v.name, v.id, err)
		}

		row = m.conn.QueryRow(ctx,
			`INSERT INTO migrations (service,id) VALUES($1,1)
				ON CONFLICT(service) DO UPDATE SET id=migrations.id+1 RETURNING id`,
			m.svcName)

		var newid uint32

		err = row.Scan(&newid)

		if err != nil {
			return 0, fmt.Errorf("id increment error: %w", err)
		}

		if newid != v.id {
			return 0, fmt.Errorf("id increment error. Expected %d, got %d", v.id, newid)
		}

		count++
	}

	return count, nil
}

func (m *Migrations) Down() error {
	ctx := context.Background()
	err := m.createTable(ctx)

	if err != nil {
		return err
	}

	for i := len(m.down) - 1; i >= 0; i-- {
		v := m.down[i]
		var curid uint32

		row := m.conn.QueryRow(ctx, `SELECT id FROM migrations WHERE service=$1`, m.svcName)
		err := row.Scan(&curid)
		if err != nil {
			if err == pgx.ErrNoRows {
				continue
			}
			return fmt.Errorf("get migration id error: %w", err)
		}

		if curid != v.id {
			return fmt.Errorf("last executed migration id sequence error. Expected %d, got %d", curid, v.id)
		}

		err = v.hnd(ctx, v.conn)
		if err != nil {
			return fmt.Errorf("error while migration %q (id: %d): %w", v.name, v.id, err)
		}

		row = m.conn.QueryRow(ctx,
			`UPDATE migrations SET id=id-1 WHERE service=$1 RETURNING id`,
			m.svcName)

		var newid uint32

		err = row.Scan(&newid)

		if err != nil {
			return fmt.Errorf("id decrement error: %w", err)
		}

		if newid != v.id-1 {
			return fmt.Errorf("id increment error. Expected %d, got %d", v.id-1, newid)
		}
	}

	return nil
}
