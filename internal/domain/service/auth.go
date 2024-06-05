package service

import (
	"context"
	"github.com/pkg/errors"
	"sso-service/internal/domain/entity"
)

var (
	ErrorWrongPassword          = errors.New("неправильный пароль")
	ErrorAccountIsInBan         = errors.New("аккаунт заблокирован")
	ErrorPasswordMustMatch      = errors.New("пароли должны совпадать")
	ErrorEmailNotConfirmed      = errors.New("почта не подтверждена")
	ErrorCouldNotExtractIdClaim = errors.New("ошибка получения id")
)

type AuthService interface {
	SignIn(ctx context.Context, usernameOrEmail, password string) (entity.TokenPair, error)
	SignUp(ctx context.Context, username, email, password, repeatedPassword string, info entity.AccountInfo) error
	Refresh(ctx context.Context, refreshToken string) (entity.TokenPair, error)
}
