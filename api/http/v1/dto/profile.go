package dto

import "sso-service/internal/domain/entity"

type Profile struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	Lastname   string `json:"lastname"`
	Patronymic string `json:"patronymic"`
}

func MapToDto(account entity.Account) Profile {
	return Profile{
		Username:   account.Username,
		Email:      account.Email,
		Name:       account.Info.Name,
		Lastname:   account.Info.Lastname,
		Patronymic: account.Info.Patronymic,
	}
}
