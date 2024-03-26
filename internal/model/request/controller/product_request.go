package request

import "github.com/guregu/null"

type ProductPatchOneByIdRequest struct {
	Name  null.String `json:"name"`
	Stock null.Int    `json:"stock"`
	Price null.String `json:"price"`
}
