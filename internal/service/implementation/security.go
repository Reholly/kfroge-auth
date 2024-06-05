package implementation

import (
	"context"
	"fmt"
	"sso-service/config"
	"sso-service/internal/domain/entity"
	"sso-service/internal/domain/service"
	"sso-service/lib/hash"
)

type SecurityService struct {
	logger         service.Logger
	codeService    service.CodeService
	accountService service.AccountService
	mailService    service.MailService
	config         config.AuthConfig
}

func NewSecurityService(
	logger service.Logger,
	codeService service.CodeService,
	accountService service.AccountService,
	mailService service.MailService,
	config config.AuthConfig,
) service.SecurityService {
	return &SecurityService{
		logger:         logger,
		codeService:    codeService,
		accountService: accountService,
		mailService:    mailService,
		config:         config,
	}
}

func (ss *SecurityService) SendResetPassword(ctx context.Context, account entity.Account) error {
	code, err := ss.codeService.CreatePasswordResetCode(ctx, account)
	if err != nil {
		ss.logger.Error(fmt.Sprintf("[SecurityService] ошибка при сбросе пароля: %s", err.Error()))
		return err
	}

	go ss.mailService.SendMail(
		ctx,
		account.Email,
		"Сброс пароля",
		fmt.Sprintf("Ваш код для сброса пароля %s", code))

	return nil
}

func (ss *SecurityService) ResetPassword(ctx context.Context, code, password, repeatedPassword string) error {
	if password != repeatedPassword {
		ss.logger.Error(fmt.Sprintf("[SecurityService] ошибка сброса пароля : пароли не совпадают. Поступивший код: %s", code))
		return service.ErrorPasswordsNotEqual
	}

	account, err := ss.codeService.GetAccountByResetPasswordCode(ctx, code)

	if err != nil {
		return err
	}

	err = ss.accountService.UpdatePassword(ctx, account, hash.GetHashSHA256(password, ss.config.PasswordSalt))
	if err != nil {
		ss.logger.Error(fmt.Sprintf("[SecurityService] ошибка сброса пароля %s", err.Error()))
		return err
	}

	return nil
}

func (ss *SecurityService) UpdateAccountInfo(ctx context.Context, account entity.Account, info entity.AccountInfo) error {
	account, err := ss.accountService.GetAccountById(ctx, account.ID)
	if err != nil {
		ss.logger.Error(fmt.Sprintf("[SecurityService] ошибка обновления профиля пользователя : %+v, текст ошибки: %s", account, err.Error()))
		return err
	}
	if err := entity.ValidateAccountInfo(info); err != nil {
		return err
	}

	return ss.accountService.UpdateMainInfo(ctx, account, info)
}

func (ss *SecurityService) ConfirmAccountEmail(ctx context.Context, code string) error {
	account, err := ss.codeService.GetAccountByEmailConfirmationCode(ctx, code)
	if err != nil {
		ss.logger.Error(fmt.Sprintf("[SecurityService] ошибка получения аккаунт при подтверждении почты, code : %s", code))
		return err
	}

	err = ss.accountService.ConfirmEmail(ctx, account)

	if err != nil {
		ss.logger.Error(fmt.Sprintf("[SecurityService] ошибка подтверждения почты для аккаунта с username: %s, ошибка : %s", account.Username, err.Error()))
		return err
	}

	return nil
}
