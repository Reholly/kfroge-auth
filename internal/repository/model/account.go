package model

import (
	"github.com/google/uuid"
	"sso-service/internal/domain/entity"
)

type AccountModel struct {
	ID         uuid.UUID `db:"id"`
	Username   string    `db:"username"`
	Email      string    `db:"email"`
	Name       string    `db:"name"`
	Lastname   string    `db:"lastname"`
	Patronymic string    `db:"patronymic"`

	IsEmailConfirmed bool `db:"is_email_confirmed"`
	IsBanned         bool `db:"is_banned"`

	HashedPassword string `db:"hashed_password"`
	Role           string `db:"role"`
}

func (model *AccountModel) MapToEntity() entity.Account {
	return entity.Account{
		ID:       model.ID,
		Username: model.Username,
		Email:    model.Email,
		Info: entity.AccountInfo{
			Name:       model.Name,
			Lastname:   model.Lastname,
			Patronymic: model.Patronymic,
		},
		IsEmailConfirmed: model.IsEmailConfirmed,
		IsBanned:         model.IsBanned,
		HashedPassword:   model.HashedPassword,
		Role:             model.Role,
	}
}
