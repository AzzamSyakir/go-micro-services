package response

import (
	"go-micro-services/src/user-service/entity"
)

type UserResponse struct {
	Users []entity.User `json:"users"`
}
