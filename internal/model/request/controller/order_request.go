package request

import "github.com/guregu/null"

type OrderProducts struct {
	Id         null.String `json:"id"`
	OrderId    null.String `json:"order_id"`
	ProductId  null.String `json:"product_id"`
	Qty        null.Int    `json:"qty"`
	TotalPrice null.Int    `json:"total_price"`
}
type OrderRequest struct {
	TotalPaid null.Int        `json:"total_price"`
	Products  []OrderProducts `json:"products"`
}
