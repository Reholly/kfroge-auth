package model

import "github.com/google/uuid"

type Code struct {
	Id        uuid.UUID `db:"id"`
	AccountId uuid.UUID `db:"account_id"`
	Value     string    `db:"value"`
}
