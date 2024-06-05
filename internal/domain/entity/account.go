package entity

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/mail"
	"regexp"
)

const (
	AdminRole     = "admin"
	ModeratorRole = "moderator"
	StudentRole   = "student"

	UsernameMinLength = 5
	UsernameMaxLength = 30
	PasswordMinLength = 5
	PasswordMaxLength = 50
	MainInfoMinLen    = 2
	MainInfoMaxLen    = 50
)

var (
	// PasswordRegex пароль должен состоять из символов обоих регистров и цифр
	PasswordRegex = regexp.MustCompile("^[0-9A-Za-z]+$")
	// UsernameRegex имя пользователя не должно содержать спец.символы, а также не начинаться с цифры.
	UsernameRegex = regexp.MustCompile("^[^[:punct:]0-9]\\w*$")

	// MainInfoRegex поля главной информации должны состоять только из латинских и русских букв
	MainInfoRegex = regexp.MustCompile("^[a-zA-Zа-яА-Я]+$")
)

var (
	ErrorUsernameLength       = errors.New("логин должен быть от 5 до 30 символов")
	ErrorEmail                = errors.New("невалидная почта")
	ErrorPasswordLength       = errors.New("пароль должен быть от 5 до 50 символов")
	ErrorPasswordRequirements = errors.New("пароль должен состоять из символов обоих регистров и цифр")
	ErrorUsernameRequirements = errors.New("имя пользователя не должно содержать спец.символы, а также начинаться с цифры")

	ErrorMainInfoLength  = errors.New("имя, фамилия или отчество должны быть от 2 до 50 символов")
	ErrorMainInfoSymbols = errors.New("имя, фамилия или отчество должны содержать только буквы латинского или русского алфавита")
)

type Account struct {
	ID       uuid.UUID
	Username string
	Email    string
	Info     AccountInfo

	IsEmailConfirmed bool
	IsBanned         bool

	HashedPassword string
	Role           string
}

type AccountInfo struct {
	Name       string
	Lastname   string
	Patronymic string
}

func ValidateAccount(username, password, email string, info AccountInfo) error {
	if len(username) < UsernameMinLength || len(username) > UsernameMaxLength {
		return ErrorUsernameLength
	}

	if !UsernameRegex.Match([]byte(username)) {
		return ErrorUsernameRequirements
	}

	if len(password) < PasswordMinLength || len(password) > PasswordMaxLength {
		return ErrorPasswordLength
	}

	if !PasswordRegex.Match([]byte(password)) {
		return ErrorPasswordRequirements
	}

	if err := ValidateAccountInfo(info); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return ErrorEmail
	}

	return nil
}

func ValidateAccountInfo(info AccountInfo) error {
	if len(info.Name) < MainInfoMinLen || len(info.Name) > MainInfoMaxLen ||
		len(info.Lastname) < MainInfoMinLen || len(info.Lastname) > MainInfoMaxLen ||
		len(info.Patronymic) < MainInfoMinLen || len(info.Patronymic) > MainInfoMaxLen {
		return ErrorMainInfoLength
	}

	if !MainInfoRegex.Match([]byte(info.Name)) {
		return ErrorMainInfoSymbols
	}

	if !MainInfoRegex.Match([]byte(info.Lastname)) {
		return ErrorMainInfoSymbols
	}

	if !MainInfoRegex.Match([]byte(info.Patronymic)) {
		return ErrorMainInfoSymbols
	}

	return nil
}
