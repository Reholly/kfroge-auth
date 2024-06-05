package dto

import "github.com/pkg/errors"

type MainInfo struct {
	Name       string `json:"name"`
	Lastname   string `json:"lastname"`
	Patronymic string `json:"patronymic"`
}

func (model *MainInfo) Validate() error {
	if model.Name == "" {
		return errors.New("поле имя должно быть заполнено")
	}

	if model.Lastname == "" {
		return errors.New("поле фамилия должно быть заполнено")
	}

	if model.Patronymic == "" {
		return errors.New("поле отчество должно быть заполнено")
	}

	return nil
}
