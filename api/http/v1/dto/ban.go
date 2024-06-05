package dto

import "github.com/pkg/errors"

type Ban struct {
	Username string `json:"username"`
	Reason   string `json:"reason"`
}

func (model *Ban) Validate() error {
	if model.Username == "" {
		return errors.New("логин должен быть заполнен")
	}

	if model.Reason == "" {
		return errors.New("причина важное поле и должно быть заполнено")
	}

	return nil
}
