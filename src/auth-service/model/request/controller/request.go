package request

import "github.com/guregu/null"

type LoginRequest struct {
	Email    null.String `json:"email"`
	Password null.String `json:"password"`
}
