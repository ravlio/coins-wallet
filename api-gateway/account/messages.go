package account

import (
	"time"

	"github.com/ravlio/wallet/account"
)

type Account struct {
	ID        uint32     `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (a *Account) FromEntity(e *account.Account) {
	a.ID = e.ID
	a.Name = e.Name
	a.Email = e.Email
	a.CreatedAt = e.CreatedAt
	a.UpdatedAt = e.UpdatedAt
}

func (a *Account) ToEntity() *account.Account {
	return &account.Account{
		ID:        a.ID,
		Name:      a.Name,
		Email:     a.Email,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

type ListAccountsResponse struct {
	Accounts []*Account `json:"accounts"`
}
