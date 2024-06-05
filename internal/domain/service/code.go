package service

import (
	"context"
	"github.com/pkg/errors"
	"sso-service/internal/domain/entity"
)

var (
	ErrorCodeService = errors.New("внутрення ошибка сервиса")
)

type CodeService interface {
	CreatePasswordResetCode(ctx context.Context, account entity.Account) (string, error)
	CreateEmailConfirmationCode(ctx context.Context, account entity.Account) (string, error)
	GetAccountByResetPasswordCode(ctx context.Context, code string) (entity.Account, error)
	GetAccountByEmailConfirmationCode(ctx context.Context, code string) (entity.Account, error)
}
