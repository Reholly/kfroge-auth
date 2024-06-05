package service

import (
	"sso-service/config"
	"sso-service/internal/domain/service"
	"sso-service/internal/repository"
	"sso-service/internal/service/implementation"
)

type ServiceManager struct {
	service.AdminService
	service.AuthService
	service.AccountService
	service.CodeService
	service.Logger
	service.MailService
	service.PermissionService
	service.SecurityService
	service.TokenService
}

func NewServiceManager(
	logger service.Logger,
	config config.Config,
	repository repository.RepositoryManager,
) ServiceManager {
	accountService := implementation.NewAccountService(repository, logger)
	codeService := implementation.NewCodeService(logger, accountService, repository)
	mailService := implementation.NewMailService(config.Smtp, logger)
	permissionService := implementation.NewPermissionService(logger)
	securityService := implementation.NewSecurityService(logger, codeService, accountService, mailService, config.Auth)
	tokenService := implementation.NewTokenService(config.Jwt, permissionService, logger)
	adminService := implementation.NewAdminService(logger, accountService, permissionService, mailService)
	authService := implementation.NewAuthService(logger, config.Auth, accountService, codeService, tokenService, permissionService, mailService)

	return ServiceManager{
		AuthService:       authService,
		AccountService:    accountService,
		CodeService:       codeService,
		Logger:            logger,
		MailService:       mailService,
		PermissionService: permissionService,
		SecurityService:   securityService,
		TokenService:      tokenService,
		AdminService:      adminService,
	}
}
