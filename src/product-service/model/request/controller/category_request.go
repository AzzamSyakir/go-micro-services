package request

import "github.com/guregu/null"

type CategoryRequest struct {
	Name null.String `json:"name"`
}
