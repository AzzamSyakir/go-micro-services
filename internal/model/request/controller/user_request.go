package request

import "github.com/guregu/null"

type UserPatchOneByIdRequest struct {
	Name  null.String `json:"name"`
	Saldo null.String `json:"saldo"`
}
