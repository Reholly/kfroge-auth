package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"sso-service/internal/domain/entity"
)

// Возможные ошибки, которые может выдать AccountRepository
var (
	ErrorAccountNotFound      = errors.New("аккаунт не найден")
	ErrorAccountAlreadyExists = errors.New("аккаунт с таким логином или почтой уже существует")
)

type AccountRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (entity.Account, error)
	FindByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (entity.Account, error)
	FindByUsername(ctx context.Context, username string) (entity.Account, error)
	FindByEmail(ctx context.Context, email string) (entity.Account, error)
	UpdatePasswordByID(ctx context.Context, id uuid.UUID, hashedPassword string) error
	UpdateInfoById(ctx context.Context, id uuid.UUID, info entity.AccountInfo) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
	Create(ctx context.Context, account entity.Account) error
	ChangeBanStatusByID(ctx context.Context, id uuid.UUID, status bool) error
	ChangeRoleByID(ctx context.Context, id uuid.UUID, role string) error
	ConfirmEmailByID(ctx context.Context, id uuid.UUID) error
}
