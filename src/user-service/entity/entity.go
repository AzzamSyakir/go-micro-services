package entity

import (
	"github.com/guregu/null"
)

type User struct {
	Id        null.String `json:"id" :"id"`
	Name      null.String `json:"name" :"name"`
	Email     null.String `json:"email" :"email"`
	Password  null.String `json:"password" :"password"`
	Balance   null.Int    `json:"balance" :"balance"`
	CreatedAt null.Time   `json:"created_at" :"created_at"`
	UpdatedAt null.Time   `json:"updated_at" :"updated_at"`
	DeletedAt null.Time   `json:"deleted_at" :"deleted_at"`
}
