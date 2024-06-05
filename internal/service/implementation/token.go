package implementation

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"sso-service/config"
	"sso-service/internal/domain/entity"
	"sso-service/internal/domain/service"
	"time"
)

type TokenService struct {
	config            config.JwtConfig
	permissionService service.PermissionService
	logger            service.Logger
}

func NewTokenService(
	config config.JwtConfig,
	permissionService service.PermissionService,
	logger service.Logger,
) service.TokenService {
	return &TokenService{
		config:            config,
		permissionService: permissionService,
		logger:            logger,
	}
}

func (ts *TokenService) CreateTokenPair(account entity.Account) (entity.TokenPair, error) {
	refreshTokenBase := jwt.New(jwt.SigningMethodHS256)
	refreshTokenClaims := refreshTokenBase.Claims.(jwt.MapClaims)
	refreshTokenClaims[entity.IdClaim] = account.ID
	refreshTokenClaims[entity.ExpirationClaim] = time.Now().
		Add(time.Hour * time.Duration(ts.config.RefreshTokenTimeToLiveInHours)).
		Unix()
	refreshToken, err := refreshTokenBase.SignedString([]byte(ts.config.JwtSecret))

	if err != nil {
		ts.logger.Error(fmt.Sprintf("[TokenService] ошибка создания refresh токена: %s", err.Error()))
		return entity.TokenPair{}, service.ErrorRefreshTokenCreate
	}

	accessTokenBase := jwt.New(jwt.SigningMethodHS256)
	accessTokenClaims := accessTokenBase.Claims.(jwt.MapClaims)
	accessTokenClaims[entity.IdClaim] = account.ID
	accessTokenClaims[entity.ExpirationClaim] = time.Now().
		Add(time.Second * time.Duration(ts.config.AccessTokenTimeToLiveInSeconds)).
		Unix()
	accessTokenClaims[entity.UsernameClaim] = account.Username
	accessTokenClaims[entity.RoleClaim] = account.Role

	accessToken, err := accessTokenBase.SignedString([]byte(ts.config.JwtSecret))

	if err != nil {
		ts.logger.Error(fmt.Sprintf("[TokenService] ошибка создания access токена %s", err.Error()))
		return entity.TokenPair{}, service.ErrorAccessTokenCreate
	}

	return entity.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (ts *TokenService) ValidateAndParseClaims(token string) ([]entity.Claim, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			ts.logger.Error(fmt.Sprintf("[TokenService] алгоритм подписи токена невалиден: %v", token.Header["alg"]))
			return nil, service.ErrorInvalidToken
		}

		return []byte(ts.config.JwtSecret), nil
	})

	if err != nil {
		ts.logger.Error(fmt.Sprintf("[TokenService] ошибка распаковки токена: %s", err.Error()))
		return nil, service.ErrorInvalidToken
	}

	if !parsedToken.Valid {
		ts.logger.Error("[TokenService] ошибка, невалидный токен")
		return nil, service.ErrorInvalidToken
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		domainClaims := make([]entity.Claim, 0, len(claims))

		for title, value := range claims {
			if title == "" || value == "" {
				continue
			}
			domainClaims = append(domainClaims, entity.Claim{
				Title: title,
				Value: value,
			})
		}

		return domainClaims, nil
	}

	ts.logger.Error(fmt.Sprintf("[TokenService] ошибка приведения к MapClaims, распаковка токена не удалась. Токен: %+v", parsedToken))
	return nil, service.ErrorParseToken
}
