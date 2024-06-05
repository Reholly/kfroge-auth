package dto

import "github.com/pkg/errors"

type Email struct {
	Value string `json:"value"`
}

func (model *Email) Validate() error {
	if model.Value == "" {
		return errors.New("поле почта не должно быть пустым")
	}

	return nil
}
