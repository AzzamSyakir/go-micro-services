package request

import "github.com/guregu/null"

type UserPatchOneByIdRequest struct {
	Name     null.String `json:"name"`
	Email    null.String `json:"email"`
	Password null.String `json:"password"`
	Balance  null.Int    `json:"balance"`
}
type CreateUser struct {
	Name     null.String `json:"name"`
	Balance  null.Int    `json:"balance"`
	Email    null.String `json:"email"`
	Password null.String `json:"password"`
}
