package entity

import (
	"github.com/guregu/null"
)

type User struct {
	Id        null.String `json:"id" bson:"id"`
	Name      null.String `json:"name" bson:"name"`
	Email     null.String `json:"email" bson:"email"`
	Password  null.String `json:"password" bson:"password"`
	Balance   null.Int    `json:"balance" bson:"balance"`
	CreatedAt null.Time   `json:"created_at" bson:"created_at"`
	UpdatedAt null.Time   `json:"updated_at" bson:"updated_at"`
	DeletedAt null.Time   `json:"deleted_at" bson:"deleted_at"`
}
