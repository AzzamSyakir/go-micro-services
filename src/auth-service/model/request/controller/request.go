package request

import "github.com/guregu/null"

type LoginRequest struct {
	Email    null.String `json:"email"`
	Password null.String `json:"password"`
}

type UserPatchOneByIdRequest struct {
	Name     null.String `json:"name"`
	Email    null.String `json:"email"`
	Password null.String `json:"password"`
	Balance  null.Int    `json:"balance"`
}
type Register struct {
	Name     null.String `json:"name"`
	Balance  null.Int    `json:"balance"`
	Email    null.String `json:"email"`
	Password null.String `json:"password"`
}
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
type CategoryRequest struct {
	Name null.String `json:"name"`
}
type OrderProducts struct {
	Id         null.String `json:"id"`
	OrderId    null.String `json:"order_id"`
	ProductId  null.String `json:"product_id"`
	Qty        null.Int    `json:"qty"`
	TotalPrice null.Int    `json:"total_price"`
}
type OrderRequest struct {
	TotalPaid null.Int        `json:"total_paid"`
	Products  []OrderProducts `json:"products"`
}
