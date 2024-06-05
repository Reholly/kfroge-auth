package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"sso-service/internal/domain/entity"
	"sso-service/internal/domain/repository"
	"sso-service/internal/domain/service"
	serviceImpl "sso-service/internal/service"
)

const (
	ClaimsContext  = "claims"
	AccountContext = "account"
)

var (
	ErrorNoClaims       = errors.New("отсутствуют claims в контексте запроса")
	ErrorInvalidClaims  = errors.New("невалидные claims")
	ErrorNoAccount      = errors.New("отсутствует аккаунт в контексте запроса")
	ErrorInvalidAccount = errors.New("невалидный аккаунт")
	ErrorNotValidBody   = errors.New("невалидное тело запроса")
)

type Api struct {
	service serviceImpl.ServiceManager
}

func NewApi(service serviceImpl.ServiceManager) Api {
	return Api{
		service: service,
	}
}

func (api *Api) ExtractClaims(c *gin.Context) ([]entity.Claim, error) {
	fromRequest, exist := c.Get(ClaimsContext)
	if !exist {
		return nil, ErrorNoClaims
	}

	claims, ok := fromRequest.([]entity.Claim)

	if !ok {
		return nil, ErrorInvalidClaims
	}

	return claims, nil
}

func (api *Api) ExtractAccount(c *gin.Context) (entity.Account, error) {
	fromRequest, exist := c.Get(AccountContext)
	if !exist {
		return entity.Account{}, ErrorNoAccount
	}

	account, ok := fromRequest.(entity.Account)

	if !ok {
		return entity.Account{}, ErrorInvalidAccount
	}

	return account, nil
}

func (api *Api) CheckAdministrationPermissions(c *gin.Context) bool {
	account, err := api.ExtractAccount(c)
	if err != nil {
		api.service.Logger.Error(fmt.Sprintf("[Api] ошибка получения аккаунта из контекста %s", err.Error()))
		return false
	}

	return api.service.PermissionService.IsInAdministration(c.Request.Context(), account)
}

func (api *Api) ChooseStatusCode(err error) int {
	switch {
	case errors.Is(err, repository.ErrorAccountNotFound),
		errors.Is(err, repository.ErrorNotFoundAccountForCode):
		return http.StatusNotFound

	case errors.Is(err, service.ErrorAccountIsInBan),
		errors.Is(err, service.ErrorAdminBan),
		errors.Is(err, service.ErrorAdminUnban),
		errors.Is(err, service.ErrorAdminChangeRole):
		return http.StatusForbidden

	case errors.Is(err, service.ErrorEmailNotConfirmed),
		errors.Is(err, ErrorNoAccount),
		errors.Is(err, ErrorInvalidAccount),
		errors.Is(err, ErrorNoClaims),
		errors.Is(err, ErrorInvalidClaims):
		return http.StatusUnauthorized

	case errors.Is(err, service.ErrorMailSending):
		return http.StatusBadGateway

	case errors.Is(err, repository.ErrorInvalidResetPasswordCodeValue),
		errors.Is(err, repository.ErrorInvalidResetPasswordCode),
		errors.Is(err, repository.ErrorAccountAlreadyExists),
		errors.Is(err, service.ErrorPasswordMustMatch),
		errors.Is(err, service.ErrorWrongPassword),
		errors.Is(err, service.ErrorNoUsernameClaim),
		errors.Is(err, service.ErrorNoIdClaim),
		errors.Is(err, service.ErrorCouldNotExtractIdClaim),
		errors.Is(err, service.ErrorPasswordsNotEqual),
		errors.Is(err, service.ErrorAccessTokenCreate),
		errors.Is(err, service.ErrorInvalidToken),
		errors.Is(err, service.ErrorParseToken),
		errors.Is(err, service.ErrorRefreshTokenCreate):
		return http.StatusBadRequest

	case errors.Is(err, service.ErrorAccountService),
		errors.Is(err, service.ErrorCodeService),
		errors.Is(err, service.ErrorCouldNotBan),
		errors.Is(err, service.ErrorCouldNotCreateModerator),
		errors.Is(err, service.ErrorCouldNotDeleteModerator):
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
