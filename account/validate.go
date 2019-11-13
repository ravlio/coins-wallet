package account

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/ravlio/wallet/pkg/errutil"
)

func (s *Service) validateCreateAccount(ctx context.Context, req *Account) error {
	return s.validateCommon(ctx, req)
}

func (s *Service) validateUpdateAccount(ctx context.Context, req *Account) error {
	return s.validateCommon(ctx, req)
}

func (s *Service) validateCommon(ctx context.Context, req *Account) error {
	if len(req.Name) == 0 {
		return errutil.NewBadRequestFieldError("name", "empty")
	}

	if err := validation.Validate(req.Email, is.Email); err != nil {
		return errutil.NewBadRequestFieldError("name", err.Error())
	}

	return nil
}
