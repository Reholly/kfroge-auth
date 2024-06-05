package dto

import "github.com/pkg/errors"

type PasswordResetConfirmation struct {
	Code             string `json:"code"`
	Password         string `json:"password"`
	RepeatedPassword string `json:"repeated_password"`
}

func (model *PasswordResetConfirmation) Validate() error {
	if model.Code == "" {
		return errors.New("поле код обязательно")
	}

	if model.Password == "" || model.RepeatedPassword == "" {
		return errors.New("поля паролей обязательны")
	}

	if model.Password != model.RepeatedPassword {
		return errors.New("поле паролей должны совпадать")
	}

	return nil
}
