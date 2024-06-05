package implementation

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"sso-service/internal/domain/entity"
	"sso-service/internal/domain/repository"
	"sso-service/internal/domain/service"
	repoImpl "sso-service/internal/repository"
)

type AccountService struct {
	repo   repoImpl.RepositoryManager
	logger service.Logger
}

func NewAccountService(repo repoImpl.RepositoryManager, logger service.Logger) service.AccountService {
	return &AccountService{
		repo:   repo,
		logger: logger,
	}
}

func (as *AccountService) ChangeBanStatus(ctx context.Context, account entity.Account, status bool) error {
	err := as.repo.Account.ChangeBanStatusByID(ctx, account.ID, status)

	as.logger.Info(fmt.Sprintf("[AccountService] у пользователя %+v сменился стату блокировки на %t", account, status))

	return as.handleError(err)
}
func (as *AccountService) ChangeRole(ctx context.Context, account entity.Account, role string) error {
	err := as.repo.Account.ChangeRoleByID(ctx, account.ID, role)

	as.logger.Info(fmt.Sprintf("[AccountService] у пользователя %+v сменилась роль на %s", account, role))

	return as.handleError(err)
}
func (as *AccountService) CreateAccount(ctx context.Context, username, email, hashedPassword string, info entity.AccountInfo) error {
	newAccount := entity.Account{
		Username:         username,
		Email:            email,
		Info:             info,
		IsEmailConfirmed: false,
		IsBanned:         false,
		HashedPassword:   hashedPassword,
		Role:             entity.StudentRole,
	}
	return as.handleError(as.repo.Account.Create(ctx, newAccount))
}

func (as *AccountService) UpdatePassword(ctx context.Context, account entity.Account, newPassword string) error {
	err := as.repo.Account.UpdatePasswordByID(ctx, account.ID, newPassword)

	as.logger.Info(fmt.Sprintf("[AccountService] у пользователя %+v сменился пароль", account))

	return as.handleError(err)
}

func (as *AccountService) UpdateMainInfo(ctx context.Context, account entity.Account, info entity.AccountInfo) error {
	err := as.repo.Account.UpdateInfoById(ctx, account.ID, info)

	as.logger.Info(fmt.Sprintf("[AccountService] у пользователя %+v сменилась главная информация на %+v", account, info))

	return as.handleError(err)
}

func (as *AccountService) GetAccountByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (entity.Account, error) {
	account, err := as.repo.Account.FindByUsernameOrEmail(ctx, usernameOrEmail)

	return account, as.handleError(err)
}

func (as *AccountService) GetAccountById(ctx context.Context, id uuid.UUID) (entity.Account, error) {
	account, err := as.repo.Account.FindByID(ctx, id)

	return account, as.handleError(err)
}

func (as *AccountService) ConfirmEmail(ctx context.Context, account entity.Account) error {
	return as.handleError(as.repo.Account.ConfirmEmailByID(ctx, account.ID))
}

func (as *AccountService) handleError(err error) error {
	if err == nil {
		return nil
	}

	as.logger.Error(fmt.Sprintf("[AccountService] ошибка: %s", err.Error()))

	causeErr := errors.Cause(err)

	switch {
	case errors.Is(causeErr, repository.ErrorAccountNotFound):
		return causeErr
	case errors.Is(causeErr, repository.ErrorAccountAlreadyExists):
		return causeErr
	case err != nil:
		return service.ErrorAccountService
	default:
		return nil
	}
}
