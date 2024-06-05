package handler

import (
	"context"
	"github.com/Reholly/kforge-proto/src/gen/auth"
	"github.com/pkg/errors"
	"sso-service/internal/service"
)

type AuthHandler struct {
	auth.UnimplementedAuthServiceServer
	service service.ServiceManager
}

func (handler *AuthHandler) ValidateToken(ctx context.Context, request *auth.ValidateTokenRequest) (*auth.ValidateTokenResponse, error) {
	claims, err := handler.service.TokenService.ValidateAndParseClaims(request.Title)

	if err != nil {
		return &auth.ValidateTokenResponse{
			Role:     "",
			Username: "",
			Email:    "",
			Id:       "",
			IsValid:  false,
		}, err
	}

	usernameClaim, err := handler.service.PermissionService.GetUsernameClaim(ctx, claims)

	if err != nil {
		return &auth.ValidateTokenResponse{
			Role:     "",
			Username: "",
			Email:    "",
			Id:       "",
			IsValid:  false,
		}, err
	}

	usernameString, ok := usernameClaim.Value.(string)
	if !ok {
		return &auth.ValidateTokenResponse{
			Role:     "",
			Username: "",
			Email:    "",
			Id:       "",
			IsValid:  false,
		}, errors.New("ошибка преобразования логина к строке")
	}

	account, err := handler.service.AccountService.GetAccountByUsernameOrEmail(ctx, usernameString)

	if err != nil {
		return &auth.ValidateTokenResponse{
			Role:     "",
			Username: "",
			Email:    "",
			Id:       "",
			IsValid:  false,
		}, err
	}

	return &auth.ValidateTokenResponse{
		Role:     account.Role,
		Username: account.Username,
		Email:    account.Email,
		Id:       account.ID.String(),
		IsValid:  true,
	}, nil
}

func NewAuthGrpcServer(service service.ServiceManager) auth.AuthServiceServer {
	return &AuthHandler{
		service: service,
	}
}
