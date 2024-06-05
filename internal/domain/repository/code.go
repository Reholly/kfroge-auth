package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var (
	ErrorNotFoundAccountForCode        = errors.New("аккаунт с таким кодом не найден")
	ErrorInvalidResetPasswordCode      = errors.New("невалидный код сброса")
	ErrorInvalidResetPasswordCodeValue = errors.New("невалидное значение")
)

type CodeRepository interface {
	CreatePasswordResetCode(ctx context.Context, account uuid.UUID) (string, error)
	CreateEmailConfirmationCode(ctx context.Context, account uuid.UUID) (string, error)
	GetAccountIDByResetPasswordCode(ctx context.Context, code string) (uuid.UUID, error)
	GetAccountIDByEmailConfirmationCode(ctx context.Context, code string) (uuid.UUID, error)
}
