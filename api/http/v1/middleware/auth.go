package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"sso-service/api/http/v1/dto"
	"sso-service/api/http/v1/handler"
	"sso-service/internal/domain/service"
	"strings"
)

const (
	AuthorizationHeader = "Authorization"
)

var (
	ErrorNoToken      = errors.New("отсутствует токен")
	ErrorAccountInBan = errors.New("аккаунт заблокирован")
)

func AuthMiddleware(
	tokenService service.TokenService,
	permissionService service.PermissionService,
	accountService service.AccountService,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerHeader := c.Request.Header.Get(AuthorizationHeader)
		if bearerHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.NewApiError(ErrorNoToken))
			return
		}

		token := strings.Split(bearerHeader, " ")[1]
		claims, err := tokenService.ValidateAndParseClaims(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.NewApiError(err))
			return
		}
		usernameClaim, err := permissionService.GetUsernameClaim(c.Request.Context(), claims)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.NewApiError(err))
			return
		}
		username := usernameClaim.Value.(string)
		account, err := accountService.GetAccountByUsernameOrEmail(c.Request.Context(), username)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.NewApiError(err))
			return
		}

		if account.IsBanned {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.NewApiError(ErrorAccountInBan))
			return
		}

		c.Set(handler.ClaimsContext, claims)
		c.Set(handler.AccountContext, account)

		c.Next()
	}
}
