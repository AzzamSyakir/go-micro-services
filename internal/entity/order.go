package entity

import "github.com/guregu/null"

type Order struct {
	Id          null.String `json:"id"`
	UserId      null.String `json:"user_id"`
	Name        null.String `json:"name"`
	ReceiptCode null.String `json:"receipt_code"`
	TotalPrice  null.Int    `json:"total_price"`
	TotalPaid   null.Int    `json:"total_paid"`
	TotalReturn null.Int    `json:"total_return"`
	CreatedAt   null.Time   `json:"created_at"`
	UpdatedAt   null.Time   `json:"updated_at"`
	DeletedAt   null.Time   `json:"deleted_at"`
}
type OrderProducts struct {
	Id         null.String `json:"id"`
	OrderId    null.String `json:"user_id"`
	ProductId  null.String `json:"product_id"`
	TotalPrice null.Int    `json:"total_price"`
	Qty        null.Int    `json:"qty"`
	CreatedAt  null.Time   `json:"created_at"`
	UpdatedAt  null.Time   `json:"updated_at"`
	DeletedAt  null.Time   `json:"deleted_at"`
}
