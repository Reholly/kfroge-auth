package implementation

import (
	"context"
	"fmt"
	"sso-service/internal/domain/entity"
	"sso-service/internal/domain/service"
)

type PermissionService struct {
	logger service.Logger
}

func NewPermissionService(logger service.Logger) service.PermissionService {
	return &PermissionService{
		logger: logger,
	}
}

func (ps *PermissionService) GetUsernameClaim(ctx context.Context, claims []entity.Claim) (entity.Claim, error) {
	for _, value := range claims {
		if value.Title == entity.UsernameClaim {
			return value, nil
		}
	}

	ps.logger.Error(fmt.Sprintf("[PermissionService] ошибка получения claim: %s ", service.ErrorNoUsernameClaim.Error()))

	return entity.Claim{}, service.ErrorNoUsernameClaim
}

func (ps *PermissionService) GetIdClaim(ctx context.Context, claims []entity.Claim) (entity.Claim, error) {
	for _, value := range claims {
		if value.Title == entity.IdClaim {
			return value, nil
		}
	}

	ps.logger.Error(fmt.Sprintf("[PermissionService] ошибка получения claim: %s ", service.ErrorNoUsernameClaim.Error()))

	return entity.Claim{}, service.ErrorNoIdClaim
}

func (ps *PermissionService) IsInAdministration(ctx context.Context, account entity.Account) bool {
	return account.Role == entity.AdminRole || account.Role == entity.ModeratorRole
}

func (ps *PermissionService) IsAdmin(ctx context.Context, account entity.Account) bool {
	return account.Role == entity.AdminRole
}

func (ps *PermissionService) IsBanned(ctx context.Context, account entity.Account) bool {
	return account.IsBanned
}
