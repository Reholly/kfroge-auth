package dto

import "github.com/pkg/errors"

type LogIn struct {
	UsernameOrEmail string `json:"login"`
	Password        string `json:"password"`
}

func (model *LogIn) Validate() error {
	if model.UsernameOrEmail == "" {
		return errors.New("логин не должен быть пустым")
	}

	if model.Password == "" {
		return errors.New("пароль не должен быть пустым")
	}

	return nil
}
