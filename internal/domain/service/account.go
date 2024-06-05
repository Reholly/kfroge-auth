package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"sso-service/internal/domain/entity"
)

var (
	ErrorAccountService = errors.New("внутрення ошибка сервиса")
)

type AccountService interface {
	ChangeBanStatus(ctx context.Context, account entity.Account, status bool) error
	ChangeRole(ctx context.Context, account entity.Account, role string) error
	CreateAccount(ctx context.Context, username, email, hashedPassword string, info entity.AccountInfo) error
	GetAccountByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (entity.Account, error)
	GetAccountById(ctx context.Context, id uuid.UUID) (entity.Account, error)
	UpdatePassword(ctx context.Context, account entity.Account, newPassword string) error
	UpdateMainInfo(ctx context.Context, account entity.Account, info entity.AccountInfo) error
	ConfirmEmail(ctx context.Context, account entity.Account) error
}
