package entity

import (
	"github.com/guregu/null"
)

type Product struct {
	Id         null.String `json:"id"`
	Sku        null.String `json:"sku"`
	Name       null.String `json:"name"`
	Stock      null.Int    `json:"stock"`
	Price      null.Int    `json:"price"`
	CategoryId null.String `json:"category_id"`
	CreatedAt  null.Time   `json:"created_at"`
	UpdatedAt  null.Time   `json:"updated_at"`
	DeletedAt  null.Time   `json:"deleted_at"`
}
