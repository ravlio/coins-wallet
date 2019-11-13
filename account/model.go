package account

import (
	"github.com/jackc/pgtype"
)

type column string

const (
	columnID    column = "id"
	columnName  column = "name"
	columnEmail column = "email"
)

var allColumns = []column{columnID, columnName, columnEmail}

type model struct {
	ID    pgtype.Int4
	Name  pgtype.Text
	Email pgtype.Text
}

func (m *model) FromEntity(e *Account) error {
	var err error

	if err = m.ID.Set(e.ID); err != nil {
		return err
	}

	if err = m.Name.Set(e.Name); err != nil {
		return err
	}

	if err = m.Email.Set(e.Email); err != nil {
		return err
	}

	return nil
}

func (m *model) ToEntity() *Account {
	acc := &Account{}
	if m.ID.Status == pgtype.Present {
		acc.ID = uint32(m.ID.Int)
	}

	if m.Name.Status == pgtype.Present {
		acc.Name = m.Name.String
	}

	if m.Email.Status == pgtype.Present {
		acc.Email = m.Email.String
	}
	return acc
}

func (m *model) Values(c ...column) []interface{} {
	if len(c) == 0 {
		c = allColumns
	}

	ret := make([]interface{}, 0, len(c))
	for _, v := range c {
		switch v {
		case columnID:
			ret = append(ret, &m.ID)
		case columnName:
			ret = append(ret, &m.Name)
		case columnEmail:
			ret = append(ret, &m.Email)
		}
	}

	return ret
}
