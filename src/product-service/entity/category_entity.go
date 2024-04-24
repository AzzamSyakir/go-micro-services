package entity

import (
	"github.com/guregu/null"
)

type Category struct {
	Id        null.String `json:"id"`
	Name      null.String `json:"name"`
	CreatedAt null.Time   `json:"created_at"`
	UpdatedAt null.Time   `json:"updated_at"`
	DeletedAt null.Time   `json:"deleted_at"`
}
