package request

import "github.com/guregu/null"

type UserPatchOneByIdRequest struct {
	Name    null.String `json:"name"`
	Balance null.Int    `json:"balance"`
}
