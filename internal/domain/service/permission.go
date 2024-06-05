package service

import (
	"context"
	"github.com/pkg/errors"
	"sso-service/internal/domain/entity"
)

var (
	ErrorNoUsernameClaim = errors.New("отсутствует username")
	ErrorNoIdClaim       = errors.New("остутствует id")
)

type PermissionService interface {
	GetUsernameClaim(ctx context.Context, claims []entity.Claim) (entity.Claim, error)
	GetIdClaim(ctx context.Context, claims []entity.Claim) (entity.Claim, error)
	IsInAdministration(ctx context.Context, account entity.Account) bool
	IsAdmin(ctx context.Context, account entity.Account) bool
}
