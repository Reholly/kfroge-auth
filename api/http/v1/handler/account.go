package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sso-service/api/http/v1/dto"
	"sso-service/internal/domain/entity"
)

const (
	ConfirmationRedirectUrl = "http://kforge.ru/#/sign_in"
)

// SendResetPasswordCode godoc
// @Tags account
// @Summary Запрос на сброс пароля
// @Description Этот эндпоинт нужен для запроса сброса пароля на почту. Пользователь вводит Email,
// @Description по этой почте ищется пользователь и генерируется одноразовый код, а затем отправляется на указанную почту.
// @Accept json
// @Produce json
// @Param input body dto.Email true "Email"
// @Success 200
// @Failure 400 {object} dto.ApiError
// @Failure 403 {object} dto.ApiError
// @Failure 500 {object} dto.ApiError
// @Failure 502 {object} dto.ApiError
// @Router /api/auth/account/sendresetcode [put]
func (api *Api) SendResetPasswordCode(c *gin.Context) {
	var email dto.Email
	if err := c.ShouldBindJSON(&email); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, dto.NewApiError(ErrorNotValidBody))
		return
	}

	account, err := api.service.AccountService.GetAccountByUsernameOrEmail(c.Request.Context(), email.Value)
	if err != nil {
		c.AbortWithStatusJSON(api.ChooseStatusCode(err), dto.NewApiError(err))
		return
	}

	err = api.service.SecurityService.SendResetPassword(c.Request.Context(), account)
	if err != nil {
		c.AbortWithStatusJSON(api.ChooseStatusCode(err), dto.NewApiError(err))
		return
	}

	c.Status(http.StatusOK)
}

// ConfirmResetPassword godoc
// @Summary Сброс пароля
// @Description Этот эндпоинт нужен для сброса пароля. Нужен одноразовый код, хранящийся 20 минут, пароль,
// @Description В теле запроса пользователь должен ввести код, пришедший на почту, новый пароль и повторенный новый пароль.
// @Tags account
// @Accept json
// @Produce json
// @Param input body dto.PasswordResetConfirmation true "Reset"
// @Success 200
// @Failure 400 {object} dto.ApiError
// @Failure 403 {object} dto.ApiError
// @Failure 500 {object} dto.ApiError
// @Failure 502 {object} dto.ApiError
// @Router /api/auth/account/confirmreset [put]
func (api *Api) ConfirmResetPassword(c *gin.Context) {
	var passwordReset dto.PasswordResetConfirmation
	if err := c.ShouldBindJSON(&passwordReset); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewApiError(err))
		return
	}

	if err := passwordReset.Validate(); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewApiError(err))
		return
	}

	err := api.service.ResetPassword(c.Request.Context(), passwordReset.Code, passwordReset.Password, passwordReset.RepeatedPassword)

	if err != nil {
		c.AbortWithStatusJSON(api.ChooseStatusCode(err), dto.NewApiError(err))
		return
	}

	c.Status(http.StatusOK)
}

// ConfirmEmail godoc
// @Tags account
// @Summary Подтверждение почты
// @Description Эндпоинт для подвтерждения почты по ссылке. Через query параметры получаются
// @Description code, а затем подтверждается эта почта уже в сервисе.
// @Accept json
// @Produce json
// @Param input query dto.Email true "Confirm"
// @Success 200
// @Failure 400 {object} dto.ApiError
// @Failure 403 {object} dto.ApiError
// @Failure 500 {object} dto.ApiError
// @Failure 502 {object} dto.ApiError
// @Router /api/auth/account/confirm [get]
func (api *Api) ConfirmEmail(c *gin.Context) {
	code := c.Query("code")

	if code == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewApiError(ErrorNotValidBody))
		return
	}

	err := api.service.ConfirmAccountEmail(c.Request.Context(), code)

	if err != nil {
		c.AbortWithStatusJSON(api.ChooseStatusCode(err), dto.NewApiError(err))
		return
	}

	c.Redirect(http.StatusFound, ConfirmationRedirectUrl)
}

// ChangeMainInfo godoc
// @Tags account
// @Summary Редактирование главной информации
// @Description Эндпоинт для редактирование главной информации аккаунта
// @Accept json
// @Produce json
// @Security Bearer
// @Param input body dto.MainInfo true "Username"
// @Success 200
// @Failure 400 {object} dto.ApiError
// @Failure 403 {object} dto.ApiError
// @Failure 500 {object} dto.ApiError
// @Failure 502 {object} dto.ApiError
// @Router /api/auth/account/info [put]
func (api *Api) ChangeMainInfo(c *gin.Context) {
	var mainInfo dto.MainInfo

	if err := c.ShouldBindJSON(&mainInfo); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewApiError(ErrorNotValidBody))
		return
	}

	if err := mainInfo.Validate(); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewApiError(ErrorNotValidBody))
		return
	}

	account, err := api.ExtractAccount(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, dto.NewApiError(err))
		return
	}

	err = api.service.AccountService.UpdateMainInfo(c.Request.Context(), account, entity.AccountInfo{
		Name:       mainInfo.Name,
		Lastname:   mainInfo.Lastname,
		Patronymic: mainInfo.Patronymic,
	})

	if err != nil {
		c.AbortWithStatusJSON(api.ChooseStatusCode(err), dto.NewApiError(err))
		return
	}

	c.Status(http.StatusOK)
}

// GetProfile godoc
// @Tags account
// @Summary Получение профиля
// @Description Эндпоинт для получения главной информации аккаунта
// @Accept json
// @Security Bearer
// @Produce json
// @Success 200 {object} dto.Profile
// @Failure 400 {object} dto.ApiError
// @Failure 403 {object} dto.ApiError
// @Failure 500 {object} dto.ApiError
// @Failure 502 {object} dto.ApiError
// @Router /api/auth/account/info [get]
func (api *Api) GetProfile(c *gin.Context) {
	account, err := api.ExtractAccount(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, dto.NewApiError(err))
		return
	}

	c.JSON(http.StatusOK, dto.MapToDto(account))
}
