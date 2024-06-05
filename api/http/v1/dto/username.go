package dto

import "github.com/pkg/errors"

type Username struct {
	Value string `json:"username"`
}

func (model *Username) Validate() error {
	if model.Value == "" {
		return errors.New("поле логин должно быть заполнено")
	}

	return nil
}
