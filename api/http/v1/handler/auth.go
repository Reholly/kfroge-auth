package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sso-service/api/http/v1/dto"
	"sso-service/internal/domain/entity"
)

// SignUp godoc
// @Tags auth
// @Summary Регистрация аккаунта
// @Description Этот эндпоинт служит для выполнения регистрации новых аккаунтов на KForge.
// @Accept json
// @Produce json
// @Param input body dto.Registration true "Register"
// @Success 200
// @Failure 400 {object} dto.ApiError
// @Failure 403 {object} dto.ApiError
// @Failure 500 {object} dto.ApiError
// @Failure 502 {object} dto.ApiError
// @Router /api/auth/register [post]
func (api *Api) SignUp(c *gin.Context) {
	var registration dto.Registration
	if err := c.ShouldBindJSON(&registration); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewApiError(ErrorNotValidBody))
		return
	}

	if err := registration.Validate(); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewApiError(err))
		return
	}

	mainInfo := entity.AccountInfo{
		Name:       registration.Name,
		Lastname:   registration.Lastname,
		Patronymic: registration.Patronymic,
	}

	err := api.service.AuthService.SignUp(
		c.Request.Context(),
		registration.Username,
		registration.Email,
		registration.Password,
		registration.RepeatedPassword,
		mainInfo)

	if err != nil {
		c.AbortWithStatusJSON(api.ChooseStatusCode(err), dto.NewApiError(err))
	}

	c.Status(http.StatusOK)
}

// SignIn godoc
// @Tags auth
// @Summary Вход
// @Description Эндпоинт для входа на KForge. На выходе эндпоинт отдает Access JWT токен (живущий 5 минут) с набором следующих claim:
// @Description 1) role : student / admin / moderator
// @Description 2) exp : время, когда токен перестанет действовать.
// @Description 3) username : имя пользователя. Вообще, нужно для связи с остальными микросервисами.
// @Description А также refresh токен, живущий 14 дней
// @Accept json
// @Produce json
// @Param input body dto.LogIn true "LogIn"
// @Success 200 {object} dto.TokenPair
// @Failure 400 {object} dto.ApiError
// @Failure 403 {object} dto.ApiError
// @Failure 500 {object} dto.ApiError
// @Router /api/auth/login [post]
func (api *Api) SignIn(c *gin.Context) {
	var loginDto dto.LogIn
	if err := c.ShouldBindJSON(&loginDto); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewApiError(ErrorNotValidBody))
		return
	}

	if err := loginDto.Validate(); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewApiError(err))
		return
	}

	tokenPair, err := api.service.SignIn(c.Request.Context(), loginDto.UsernameOrEmail, loginDto.Password)

	if err != nil {
		code := api.ChooseStatusCode(err)
		c.AbortWithStatusJSON(code, dto.NewApiError(err))
		return
	}

	c.JSON(http.StatusOK, dto.TokenPair{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	})
}

// Refresh godoc
// @Tags auth
// @Summary Обновление токена
// @Description Эндпоинт для обновления токенов по refresh токену. Результат состоит из пары новых refresh + access токенов
// @Accept json
// @Produce json
// @Param input body dto.Token true "Refresh"
// @Success 200 {object} dto.TokenPair
// @Failure 400 {object} dto.ApiError
// @Failure 403 {object} dto.ApiError
// @Failure 500 {object} dto.ApiError
// @Router /api/auth/refresh [post]
func (api *Api) Refresh(c *gin.Context) {
	var token dto.Token

	if err := c.ShouldBindJSON(&token); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewApiError(ErrorNotValidBody))
		return
	}

	tokenPair, err := api.service.AuthService.Refresh(c.Request.Context(), token.Value)

	if err != nil {
		c.AbortWithStatusJSON(api.ChooseStatusCode(err), dto.NewApiError(err))
		return
	}

	c.JSON(http.StatusOK, dto.TokenPair{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	})
}
