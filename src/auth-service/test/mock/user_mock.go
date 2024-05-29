package mock

import (
	"go-micro-services/src/auth-service/entity"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
)

type UserMock struct {
	Data []*entity.User
}

func NewUserMock() *UserMock {
	currentTime := time.Now()
	currentTimeRfc3339 := currentTime.Format(time.RFC3339)
	currentTimeFromRfc3339, parseErr := time.Parse(time.RFC3339, currentTimeRfc3339)
	if parseErr != nil {
		panic(parseErr)
	}
	userMock := &UserMock{
		Data: []*entity.User{
			{
				Id:        null.NewString(uuid.NewString(), true),
				Name:      null.NewString("name"+uuid.NewString(), true),
				Balance:   null.NewInt(100000, true),
				Email:     null.NewString("email"+uuid.NewString()+"@mail.com", true),
				Password:  null.NewString("password0", true),
				CreatedAt: null.NewTime(currentTimeFromRfc3339.Add(0*time.Second), true),
				UpdatedAt: null.NewTime(currentTimeFromRfc3339.Add(0*time.Second), true),
			},
			{
				Id:        null.NewString(uuid.NewString(), true),
				Name:      null.NewString("name"+uuid.NewString(), true),
				Balance:   null.NewInt(100000, true),
				Email:     null.NewString("email"+uuid.NewString()+"@mail.com", true),
				Password:  null.NewString("password1", true),
				CreatedAt: null.NewTime(currentTimeFromRfc3339.Add(1*time.Second), true),
				UpdatedAt: null.NewTime(currentTimeFromRfc3339.Add(1*time.Second), true),
			},
		},
	}
	return userMock
}
