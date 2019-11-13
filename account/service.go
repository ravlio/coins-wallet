package account

import (
	"context"
)

type Account struct {
	ID    uint32
	Name  string
	Email string
}
type Client interface {
	CreateAccount(ctx context.Context, req *Account) (*Account, error)
	GetAccount(ctx context.Context, id uint32) (*Account, error)
	DeleteAccount(ctx context.Context, id uint32) error
	UpdateAccount(ctx context.Context, req *Account) (*Account, error)
	ListAccounts(ctx context.Context) ([]*Account, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateAccount(ctx context.Context, req *Account) (*Account, error) {
	if err := s.validateCreateAccount(ctx, req); err != nil {
		return nil, err
	}

	return s.repo.CreateAccount(ctx, req)
}

func (s *Service) GetAccount(ctx context.Context, id uint32) (*Account, error) {
	return s.repo.GetAccountByID(ctx, id)
}

func (s *Service) DeleteAccount(ctx context.Context, id uint32) error {
	return s.repo.DeleteAccount(ctx, id)
}

func (s *Service) UpdateAccount(ctx context.Context, req *Account) (*Account, error) {
	if err := s.validateUpdateAccount(ctx, req); err != nil {
		return nil, err
	}

	return s.repo.UpdateAccount(ctx, req)
}

func (s *Service) ListAccounts(ctx context.Context) ([]*Account, error) {
	return s.repo.ListAccounts(ctx)
}
