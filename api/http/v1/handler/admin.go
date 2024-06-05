package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sso-service/api/http/v1/dto"
)

// CreateModerator godoc
// @Tags admin
// @Security Bearer
// @Summary Создание модератора
// @Description Этот эндпоинт нужен для накидывания роли модератора на пользователя по Username, доступен только пользователям с ролью admin.
// @Accept json
// @Produce json
// @Param input body dto.Username true "Username"
// @Success 200
// @Failure 400 {object} dto.ApiError
// @Failure 403 {object} dto.ApiError
// @Failure 500 {object} dto.ApiError
// @Failure 502 {object} dto.ApiError
// @Router /api/auth/admin/createmoder [put]
func (api *Api) CreateModerator(c *gin.Context) {
	if !api.CheckAdministrationPermissions(c) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	var username dto.Username
	if err := c.ShouldBindJSON(&username); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewApiError(ErrorNotValidBody))
		return
	}

	if err := username.Validate(); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewApiError(err))
		return
	}

	user, err := api.service.AccountService.GetAccountByUsernameOrEmail(c.Request.Context(), username.Value)

	if err != nil {
		c.AbortWithStatusJSON(api.ChooseStatusCode(err), dto.NewApiError(err))
		return
	}

	err = api.service.AdminService.CreateModerator(c.Request.Context(), user)

	if err != nil {
		c.AbortWithStatusJSON(api.ChooseStatusCode(err), dto.NewApiError(err))
		return
	}

	c.Status(http.StatusOK)
}

// DeleteModerator godoc
// @Tags admin
// @Security Bearer
// @Summary Удаление модератора
// @Description Этот эндпоинт нужен для снятия роли модератора с пользователя по Username, доступен только пользователям с ролью admin.
// @Accept json
// @Produce json
// @Param input body dto.Username true "Username"
// @Success 200
// @Failure 400 {object} dto.ApiError
// @Failure 403 {object} dto.ApiError
// @Failure 500 {object} dto.ApiError
// @Failure 502 {object} dto.ApiError
// @Router /api/auth/admin/deletemoder [put]
func (api *Api) DeleteModerator(c *gin.Context) {
	if !api.CheckAdministrationPermissions(c) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	var username dto.Username
	if err := c.ShouldBindJSON(&username); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewApiError(ErrorNotValidBody))
		return
	}

	if err := username.Validate(); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewApiError(err))
		return
	}

	user, err := api.service.AccountService.GetAccountByUsernameOrEmail(c.Request.Context(), username.Value)

	if err != nil {
		c.AbortWithStatusJSON(api.ChooseStatusCode(err), dto.NewApiError(err))
		return
	}

	err = api.service.AdminService.DeleteModerator(c.Request.Context(), user)

	if err != nil {
		c.AbortWithStatusJSON(api.ChooseStatusCode(err), dto.NewApiError(err))
		return
	}

	c.Status(http.StatusOK)
}

// BanUser godoc
// @Tags admin
// @Security Bearer
// @Summary Бан аккаунта
// @Description Этот эндпоинт нужен для бана пользователя по Username с указанием причины бана.
// @Description По почте пользователю отправляется письмо с причиной бана.
// @Accept json
// @Produce json
// @Param input body dto.Ban true "Ban"
// @Success 200
// @Failure 400 {object} dto.ApiError
// @Failure 403 {object} dto.ApiError
// @Failure 500 {object} dto.ApiError
// @Failure 502 {object} dto.ApiError
// @Router /api/auth/admin/ban [put]
func (api *Api) BanUser(c *gin.Context) {
	if !api.CheckAdministrationPermissions(c) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	var ban dto.Ban
	if err := c.ShouldBindJSON(&ban); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewApiError(ErrorNotValidBody))
		return
	}

	if err := ban.Validate(); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewApiError(err))
		return
	}

	user, err := api.service.AccountService.GetAccountByUsernameOrEmail(c.Request.Context(), ban.Username)
	if err != nil {
		c.AbortWithStatusJSON(api.ChooseStatusCode(err), dto.NewApiError(err))
		return
	}

	err = api.service.AdminService.BanUser(c.Request.Context(), user, ban.Reason)
	if err != nil {
		c.AbortWithStatusJSON(api.ChooseStatusCode(err), dto.NewApiError(err))
		return
	}

	c.Status(http.StatusOK)
}

// UnbanUser godoc
// @Tags admin
// @Security Bearer
// @Summary Разбан
// @Description Этот эндпоинт нужен для снятия бана с пользователя по Username.
// @Description По почте пользователю отправляется сообщение о разбане.
// @Accept json
// @Produce json
// @Param input body dto.Username true "Ban"
// @Success 200
// @Failure 400 {object} dto.ApiError
// @Failure 403 {object} dto.ApiError
// @Failure 500 {object} dto.ApiError
// @Failure 502 {object} dto.ApiError
// @Router /api/auth/admin/unban [put]
func (api *Api) UnbanUser(c *gin.Context) {
	if !api.CheckAdministrationPermissions(c) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	var username dto.Username
	if err := c.ShouldBindJSON(&username); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewApiError(ErrorNotValidBody))
		return
	}

	if err := username.Validate(); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.NewApiError(err))
		return
	}

	user, err := api.service.AccountService.GetAccountByUsernameOrEmail(c.Request.Context(), username.Value)
	if err != nil {
		c.AbortWithStatusJSON(api.ChooseStatusCode(err), dto.NewApiError(err))
		return
	}

	err = api.service.AdminService.UnbanUser(c.Request.Context(), user)

	if err != nil {
		c.AbortWithStatusJSON(api.ChooseStatusCode(err), dto.NewApiError(err))
		return
	}

	c.Status(http.StatusOK)
}
