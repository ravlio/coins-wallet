package account

import (
	"github.com/jackc/pgtype"
	"github.com/ravlio/wallet/pkg/db"
)

var allColumns = []string{"id", "name", "email", "created_at", "updated_at"}
var selectColumns = []string{"id", "name", "email", "created_at", "updated_at"}
var insertColumns = []string{"name", "email"}
var updateColumns = []string{"name", "email"}

type model struct {
	ID        pgtype.Int4
	Name      pgtype.Text
	Email     pgtype.Text
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

func (m *model) FromEntity(e *Account) error {
	setter := db.NewModelSetter()
	setter.Set(&m.ID, e.ID)
	setter.Set(&m.Name, e.Name)
	setter.Set(&m.Email, e.Email)
	setter.Set(&m.CreatedAt, e.CreatedAt)
	setter.Set(&m.UpdatedAt, e.UpdatedAt)

	return setter.Apply()
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

	if m.CreatedAt.Status == pgtype.Present {
		*acc.CreatedAt = m.CreatedAt.Time
	}

	if m.UpdatedAt.Status == pgtype.Present {
		*acc.UpdatedAt = m.UpdatedAt.Time
	}
	return acc
}

func (m *model) Values(c ...string) []interface{} {
	if len(c) == 0 {
		c = allColumns
	}

	ret := make([]interface{}, 0, len(c))
	for _, v := range c {
		switch v {
		case "id":
			ret = append(ret, &m.ID)
		case "name":
			ret = append(ret, &m.Name)
		case "email":
			ret = append(ret, &m.Email)
		}
	}

	return ret
}
