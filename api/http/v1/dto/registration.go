package dto

import (
	"github.com/pkg/errors"
)

type Registration struct {
	Username         string `json:"username"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	RepeatedPassword string `json:"repeated_password"`

	Name       string `json:"name"`
	Lastname   string `json:"lastname"`
	Patronymic string `json:"patronymic"`
}

func (model *Registration) Validate() error {
	if model.Username == "" {
		return errors.New("логин не должен быть пустым")
	}

	if model.Email == "" {
		return errors.New("почта не должна быть пустой")
	}

	if model.Password == "" || model.RepeatedPassword == "" {
		return errors.New("поля паролей не должны быть пустыми")
	}

	if model.Password != model.RepeatedPassword {
		return errors.New("поля паролей должны совпадать")
	}

	if model.Name == "" {
		return errors.New("имя обязательное поле")
	}

	if model.Lastname == "" {
		return errors.New("фамилия обязательное поле")
	}

	if model.Patronymic == "" {
		return errors.New("отчество обязательное поле")
	}

	return nil
}
