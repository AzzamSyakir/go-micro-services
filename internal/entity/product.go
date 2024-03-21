package entity

import (
	"github.com/guregu/null"
)

type Product struct {
	Id        null.String `json:"id" bson:"id"`
	Name      null.String `json:"name" bson:"name"`
	Price     null.String `json:"Price" bson:"price"`
	Stock     null.Int    `json:"stock" bson:"stock"`
	CreatedAt null.Time   `json:"created_at" bson:"created_at"`
	UpdatedAt null.Time   `json:"updated_at" bson:"updated_at"`
	DeletedAt null.Time   `json:"deleted_at" bson:"deleted_at"`
}
