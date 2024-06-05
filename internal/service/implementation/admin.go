package implementation

import (
	"context"
	"fmt"
	"sso-service/internal/domain/entity"
	"sso-service/internal/domain/service"
)

type AdminService struct {
	logger            service.Logger
	accountService    service.AccountService
	permissionService service.PermissionService
	mailService       service.MailService
}

func NewAdminService(
	logger service.Logger,
	accountService service.AccountService,
	permissionService service.PermissionService,
	mailService service.MailService) service.AdminService {
	return &AdminService{
		logger:            logger,
		accountService:    accountService,
		permissionService: permissionService,
		mailService:       mailService,
	}
}

func (as *AdminService) BanUser(ctx context.Context, account entity.Account, reason string) error {
	if as.permissionService.IsAdmin(ctx, account) {
		return service.ErrorAdminBan
	}

	if err := as.accountService.ChangeBanStatus(ctx, account, true); err != nil {
		as.logger.Error(fmt.Sprintf("[AdminService] ошибка при бане пользователя %+v, ошибка : %s", account, err.Error()))

		return service.ErrorCouldNotBan
	}

	go as.mailService.SendMail(ctx, account.Email, "Бан",
		fmt.Sprintf("Вас заблокировали на платформе kforge.ru по причине %s", reason))

	return nil
}

func (as *AdminService) UnbanUser(ctx context.Context, account entity.Account) error {
	if as.permissionService.IsAdmin(ctx, account) {
		return service.ErrorAdminUnban
	}

	if err := as.accountService.ChangeBanStatus(ctx, account, false); err != nil {
		as.logger.Error(fmt.Sprintf("[AdminService] ошибка при разбане пользователя %+v, ошибка : %s ", account, err.Error()))

		return service.ErrorCouldNotBan
	}

	go as.mailService.SendMail(ctx, account.Email, "Бан", "Вас разблокировали на платформе kforge.ru по причине")

	return nil
}

func (as *AdminService) CreateModerator(ctx context.Context, account entity.Account) error {
	if as.permissionService.IsAdmin(ctx, account) {
		return service.ErrorAdminChangeRole
	}

	if err := as.accountService.ChangeRole(ctx, account, entity.ModeratorRole); err != nil {
		as.logger.Error(fmt.Sprintf("[AccountService] ошибка при назначении модерки аккаунту %+v, ошибка : %s", account, err.Error()))
		return service.ErrorCouldNotCreateModerator
	}

	go as.mailService.SendMail(ctx, account.Email, "Полномочия", "Вы были назначены модератором на платформе kforge.ru")

	return nil
}

func (as *AdminService) DeleteModerator(ctx context.Context, account entity.Account) error {
	if as.permissionService.IsAdmin(ctx, account) {
		return service.ErrorAdminChangeRole
	}

	if err := as.accountService.ChangeRole(ctx, account, entity.StudentRole); err != nil {
		as.logger.Error(fmt.Sprintf("[AccountService] ошибка при снятии модерки аккаунту %+v, ошибка : %s", account, err.Error()))
		return service.ErrorCouldNotDeleteModerator
	}

	go as.mailService.SendMail(ctx, account.Email, "Полномочия", "Вы были сняты с модерки на платформе kforge.ru")

	return nil
}
