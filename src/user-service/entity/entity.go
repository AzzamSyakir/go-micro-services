package entity

import (
	"github.com/guregu/null"
)

type User struct {
	Id        null.String `json:"id"`
	Name      null.String `json:"name" `
	Email     null.String `json:"email"`
	Password  null.String `json:"password"`
	Balance   null.Int    `json:"balance"`
	CreatedAt null.Time   `json:"created_at"`
	UpdatedAt null.Time   `json:"updated_at"`
	DeletedAt null.Time   `json:"deleted_at"`
}
