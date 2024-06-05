package service

import (
	"context"
	"github.com/pkg/errors"
	"sso-service/internal/domain/entity"
)

var (
	ErrorAdminBan                = errors.New("невозможно забанить главного администратора")
	ErrorAdminUnban              = errors.New("админ и так не забанен")
	ErrorAdminChangeRole         = errors.New("админ не может сменить роль")
	ErrorCouldNotBan             = errors.New("ошибка при бане пользователя")
	ErrorCouldNotCreateModerator = errors.New("ошибка при создании модератора")
	ErrorCouldNotDeleteModerator = errors.New("ошибка при удалении модератора")
)

type AdminService interface {
	BanUser(ctx context.Context, account entity.Account, reason string) error
	UnbanUser(ctx context.Context, account entity.Account) error
	CreateModerator(ctx context.Context, account entity.Account) error
	DeleteModerator(ctx context.Context, account entity.Account) error
}
