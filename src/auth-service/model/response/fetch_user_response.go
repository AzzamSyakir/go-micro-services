package response

import (
	"go-micro-services/src/auth-service/entity"
)

type SessionResponse struct {
	Session []entity.Session `json:"session"`
}
