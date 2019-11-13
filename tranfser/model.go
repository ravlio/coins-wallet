package tranfser

import (
	"github.com/jackc/pgtype"
	"github.com/ravlio/wallet/pkg/money"
)

type column string

const (
	columnID        column = "id"
	columnFromID    column = "from_id"
	columnToID      column = "to_id"
	columnCurrency  column = "currency"
	columnAmount    column = "amount"
	columnStatus    column = "status"
	columnCreatedAt column = "created_at"
)

var allColumns = []column{
	columnID, columnFromID, columnToID, columnCurrency,
	columnAmount, columnStatus, columnCreatedAt,
}

type modelStatus string

const (
	modelStatusProcessing modelStatus = "processing"
	modelStatusSuccess    modelStatus = "success"
	modelStatusFailure    modelStatus = "failure"
)

type model struct {
	ID        pgtype.Int4
	FromID    pgtype.Int4
	ToID      pgtype.Int4
	Currency  pgtype.Int2
	Amount    pgtype.Int8
	Status    pgtype.Text
	CreatedAt pgtype.Timestamp
}

func (m *model) FromEntity(e *Transfer) error {
	var err error

	if err = m.ID.Set(e.ID); err != nil {
		return err
	}

	if err = m.FromID.Set(e.FromID); err != nil {
		return err
	}

	if err = m.ToID.Set(e.ToID); err != nil {
		return err
	}

	if err = m.Currency.Set(e.Currency); err != nil {
		return err
	}

	if err = m.Amount.Set(e.Amount); err != nil {
		return err
	}

	switch e.Status {
	case StatusProcessing:
		if err = m.Status.Set(modelStatusProcessing); err != nil {
			return err
		}
	case StatusSuccess:
		if err = m.Status.Set(modelStatusSuccess); err != nil {
			return err
		}
	case StatusFailure:
		if err = m.Status.Set(modelStatusFailure); err != nil {
			return err
		}
	}

	if err = m.CreatedAt.Set(e.CreatedAt); err != nil {
		return err
	}

	return nil
}

func (m *model) ToEntity() *Transfer {
	acc := &Transfer{}
	if m.ID.Status == pgtype.Present {
		acc.ID = uint32(m.ID.Int)
	}

	if m.FromID.Status == pgtype.Present {
		acc.FromID = uint32(m.FromID.Int)
	}

	if m.ToID.Status == pgtype.Present {
		acc.ToID = uint32(m.ToID.Int)
	}

	if m.Currency.Status == pgtype.Present {
		acc.Currency = money.Currency(m.Currency.Int)
	}

	if m.Amount.Status == pgtype.Present {
		acc.Amount = money.Money(m.Amount.Int)
	}

	if m.Status.Status == pgtype.Present {
		switch modelStatus(m.Status.String) {
		case modelStatusProcessing:
			acc.Status = StatusProcessing
		case modelStatusSuccess:
			acc.Status = StatusSuccess
		case modelStatusFailure:
			acc.Status = StatusFailure
		}
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
		case columnFromID:
			ret = append(ret, &m.FromID)
		case columnToID:
			ret = append(ret, &m.ToID)
		case columnCurrency:
			ret = append(ret, &m.Currency)
		case columnAmount:
			ret = append(ret, &m.Amount)
		case columnStatus:
			ret = append(ret, &m.Status)
		case columnCreatedAt:
			ret = append(ret, &m.CreatedAt)
		}
	}

	return ret
}
