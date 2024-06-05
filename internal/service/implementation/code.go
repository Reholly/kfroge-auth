package implementation

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"sso-service/internal/domain/entity"
	"sso-service/internal/domain/repository"
	"sso-service/internal/domain/service"
	repositoryImpl "sso-service/internal/repository"
)

type CodeService struct {
	logger  service.Logger
	account service.AccountService
	repo    repositoryImpl.RepositoryManager
}

func NewCodeService(logger service.Logger, accountService service.AccountService, manager repositoryImpl.RepositoryManager) service.CodeService {
	return &CodeService{
		logger:  logger,
		repo:    manager,
		account: accountService,
	}
}

func (cs *CodeService) CreatePasswordResetCode(ctx context.Context, account entity.Account) (string, error) {
	code, err := cs.repo.Code.CreatePasswordResetCode(ctx, account.ID)

	return code, cs.handleError(err)
}

func (cs *CodeService) CreateEmailConfirmationCode(ctx context.Context, account entity.Account) (string, error) {
	code, err := cs.repo.Code.CreateEmailConfirmationCode(ctx, account.ID)

	return code, cs.handleError(err)
}

func (cs *CodeService) GetAccountByResetPasswordCode(ctx context.Context, code string) (entity.Account, error) {
	id, err := cs.repo.Code.GetAccountIDByResetPasswordCode(ctx, code)

	if err != nil {
		return entity.Account{}, cs.handleError(err)
	}

	account, err := cs.account.GetAccountById(ctx, id)

	return account, cs.handleError(err)
}

func (cs *CodeService) GetAccountByEmailConfirmationCode(ctx context.Context, code string) (entity.Account, error) {
	id, err := cs.repo.Code.GetAccountIDByEmailConfirmationCode(ctx, code)

	if err != nil {
		return entity.Account{}, cs.handleError(err)
	}

	account, err := cs.account.GetAccountById(ctx, id)

	return account, cs.handleError(err)
}

func (cs *CodeService) handleError(err error) error {
	if err == nil {
		return nil
	}

	cs.logger.Error(fmt.Sprintf("[CodeService] ошибка: %s", err.Error()))

	causeErr := errors.Cause(err)

	switch {
	case errors.Is(err, repository.ErrorAccountNotFound):
		return causeErr
	case errors.Is(causeErr, repository.ErrorInvalidResetPasswordCode):
		return causeErr
	case errors.Is(causeErr, repository.ErrorInvalidResetPasswordCodeValue):
		return causeErr
	case errors.Is(causeErr, repository.ErrorNotFoundAccountForCode):
		return causeErr
	case err != nil:
		return service.ErrorCodeService
	default:
		return nil
	}
}
