package implementation

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"net/url"
	"sso-service/config"
	"sso-service/internal/domain/entity"
	"sso-service/internal/domain/service"
	"sso-service/lib/hash"
)

const (
	CodeParam = "code"
)

type AuthService struct {
	logger            service.Logger
	config            config.AuthConfig
	accountService    service.AccountService
	codeService       service.CodeService
	tokenService      service.TokenService
	mailService       service.MailService
	permissionService service.PermissionService
}

func NewAuthService(
	logger service.Logger,
	authConfig config.AuthConfig,
	accountService service.AccountService,
	codeService service.CodeService,
	tokenService service.TokenService,
	permissionService service.PermissionService,
	mailService service.MailService) service.AuthService {
	return &AuthService{
		logger:            logger,
		config:            authConfig,
		accountService:    accountService,
		codeService:       codeService,
		tokenService:      tokenService,
		mailService:       mailService,
		permissionService: permissionService,
	}
}

func (as *AuthService) SignIn(ctx context.Context, usernameOrEmail, password string) (entity.TokenPair, error) {
	account, err := as.accountService.GetAccountByUsernameOrEmail(ctx, usernameOrEmail)

	if err != nil {
		as.logger.Error(fmt.Sprintf("[AuthService] ошибка входа по кредам %s : %s, ошибка: %s", usernameOrEmail, password, err.Error()))
		return entity.TokenPair{}, err
	}

	if account.HashedPassword != hash.GetHashSHA256(password, as.config.PasswordSalt) {
		as.logger.Error(fmt.Sprintf("[AuthService] ошибка входа по кредам %s : %s, неправильный пароль", usernameOrEmail, password))
		return entity.TokenPair{}, service.ErrorWrongPassword
	}

	if account.IsBanned {
		as.logger.Error(fmt.Sprintf("[AuthService] ошибка входа по кредам : %s, аккаунт заблокирован", usernameOrEmail))
		return entity.TokenPair{}, service.ErrorAccountIsInBan
	}

	if !account.IsEmailConfirmed {
		as.logger.Error(fmt.Sprintf("[AuthService] ошибка входа по кредам : %s, почта не подтверждена", usernameOrEmail))
		return entity.TokenPair{}, service.ErrorEmailNotConfirmed
	}

	tokenPair, err := as.tokenService.CreateTokenPair(account)

	if err != nil {
		as.logger.Error(fmt.Sprintf("[AuthService] ошибка генерации токенов %s", err.Error()))
		return entity.TokenPair{}, err
	}

	return tokenPair, nil
}

func (as *AuthService) SignUp(ctx context.Context, username, email, password, repeatedPassword string, info entity.AccountInfo) error {
	if password != repeatedPassword {
		return service.ErrorPasswordMustMatch
	}

	if err := entity.ValidateAccount(username, password, email, info); err != nil {
		return err
	}

	if err := as.accountService.CreateAccount(ctx, username, email, hash.GetHashSHA256(password, as.config.PasswordSalt), info); err != nil {
		as.logger.Error(fmt.Sprintf("[AuthService] ошибка регистрации аккаунта с логином : %s и почтой : %s, ошибка: %s", username, email, err.Error()))
		return err
	}

	account, err := as.accountService.GetAccountByUsernameOrEmail(ctx, username)

	if err != nil {
		as.logger.Error(fmt.Sprintf("[AuthService] ошибка получения аккаунта с логином : %s", username))
		return err
	}

	code, err := as.codeService.CreateEmailConfirmationCode(ctx, account)
	if err != nil {
		as.logger.Error(fmt.Sprintf("[AuthService] ошибка создания кода для подтверждения почты для аккаунта с логином : %s", username))
		return err
	}

	params := url.Values{}
	params.Set(CodeParam, code)
	confirmationUrl := fmt.Sprintf("%s?%s", as.config.EmailConfirmationUrlBase, params.Encode())
	confirmMessage := fmt.Sprintf("Для подтверждения почты перейдите по ссылке %s", confirmationUrl)
	go as.mailService.SendMail(ctx, email, "Подтверждение почты", confirmMessage)

	return nil
}

func (as *AuthService) Refresh(ctx context.Context, refreshToken string) (entity.TokenPair, error) {
	claims, err := as.tokenService.ValidateAndParseClaims(refreshToken)

	if err != nil {
		as.logger.Error(fmt.Sprintf("[AuthService] ошибка распаковки токена %s", refreshToken))
		return entity.TokenPair{}, err
	}

	idClaim, err := as.permissionService.GetIdClaim(ctx, claims)
	if err != nil {
		return entity.TokenPair{}, err
	}

	if idClaim.Value == "" {
		as.logger.Error("[AuthService] ошибка получения claim id")
		return entity.TokenPair{}, service.ErrorCouldNotExtractIdClaim
	}

	id, err := uuid.Parse(idClaim.Value.(string))
	if err != nil {
		as.logger.Error("[AuthService] ошибка получения claim id : невозможно спарсить строку")
		return entity.TokenPair{}, service.ErrorCouldNotExtractIdClaim
	}

	account, err := as.accountService.GetAccountById(ctx, id)

	if account.IsBanned {
		return entity.TokenPair{}, service.ErrorAccountIsInBan
	}

	pair, err := as.tokenService.CreateTokenPair(account)
	if err != nil {
		as.logger.Error(fmt.Sprintf("[AuthService] ошибка создания пары токенов для аккаунта : %+v", account))
		return entity.TokenPair{}, err
	}

	return pair, nil
}
