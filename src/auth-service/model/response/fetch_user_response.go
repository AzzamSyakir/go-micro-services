package response

import (
	"go-micro-services/src/auth-service/entity"
)

type UserResponse struct {
	Users []entity.Auth `json:"users"`
}
