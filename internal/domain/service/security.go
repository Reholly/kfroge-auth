package service

import (
	"context"
	"github.com/pkg/errors"
	"sso-service/internal/domain/entity"
)

var (
	ErrorPasswordsNotEqual = errors.New("пароли не совпадают")
)

type SecurityService interface {
	SendResetPassword(ctx context.Context, account entity.Account) error
	ResetPassword(ctx context.Context, code, password, repeatedPassword string) error
	UpdateAccountInfo(ctx context.Context, account entity.Account, info entity.AccountInfo) error
	ConfirmAccountEmail(ctx context.Context, code string) error
}
