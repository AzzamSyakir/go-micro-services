package request

import "github.com/guregu/null"

type ProductPatchOneByIdRequest struct {
	Name       null.String `json:"name"`
	Stock      null.Int    `json:"stock"`
	Price      null.Int    `json:"price"`
	CategoryId null.String `json:"category_id"`
}
type CreateProduct struct {
	CategoryId null.String `json:"category_id"`
	Name       null.String `json:"name"`
	Stock      null.Int    `json:"stock"`
	Price      null.Int    `json:"price"`
}
