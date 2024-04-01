package request

import "github.com/guregu/null"

type Products struct {
	ProductId null.String `json:"product_id"`
	Qty       null.Int    `json:"qty"`
}
type OrderRequest struct {
	TotalPaid null.String `json:"name"`
	Products  Products    `json:"products"`
}
