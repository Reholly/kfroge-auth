package service

import (
	"github.com/pkg/errors"
	"sso-service/internal/domain/entity"
)

var (
	ErrorRefreshTokenCreate = errors.New("невозможно создать refresh токен для пользователя")
	ErrorAccessTokenCreate  = errors.New("невозможно создать access токен для пользователя")
	ErrorInvalidToken       = errors.New("токен невалидный")
	ErrorParseToken         = errors.New("ошибка распаковки токена")
)

type TokenService interface {
	CreateTokenPair(account entity.Account) (entity.TokenPair, error)
	ValidateAndParseClaims(token string) ([]entity.Claim, error)
}
