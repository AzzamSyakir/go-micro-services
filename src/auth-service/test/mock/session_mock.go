package mock

import (
	"go-micro-services/src/auth-service/entity"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
)

type SessionMock struct {
	Data     []*entity.Session
	UserMock *UserMock
}

func NewSessionMock(
	userMock *UserMock,
) *SessionMock {
	currentTime := time.Now()
	currentTimeRfc3339 := currentTime.Format(time.RFC3339)
	currentTimeFromRfc3339, parseErr := time.Parse(time.RFC3339, currentTimeRfc3339)
	if parseErr != nil {
		panic(parseErr)
	}
	sessionMock := &SessionMock{
		UserMock: userMock,
		Data: []*entity.Session{
			{
				Id:                    null.NewString(uuid.NewString(), true),
				UserId:                null.NewString(userMock.Data[0].Id.String, true),
				AccessToken:           null.NewString(uuid.NewString(), true),
				RefreshToken:          null.NewString(uuid.NewString(), true),
				AccessTokenExpiredAt:  null.NewTime(currentTimeFromRfc3339.Add(time.Minute*10), true),
				RefreshTokenExpiredAt: null.NewTime(currentTimeFromRfc3339.Add(time.Hour*24*2), true),
				CreatedAt:             null.NewTime(currentTimeFromRfc3339.Add(time.Second*0), true),
				UpdatedAt:             null.NewTime(currentTimeFromRfc3339.Add(time.Second*0), true),
				DeletedAt:             null.NewTime(time.Time{}, false),
			},
			{
				Id:                    null.NewString(uuid.NewString(), true),
				UserId:                null.NewString(userMock.Data[1].Id.String, true),
				AccessToken:           null.NewString(uuid.NewString(), true),
				RefreshToken:          null.NewString(uuid.NewString(), true),
				AccessTokenExpiredAt:  null.NewTime(currentTimeFromRfc3339.Add(time.Minute*10), true),
				RefreshTokenExpiredAt: null.NewTime(currentTimeFromRfc3339.Add(time.Hour*24*2), true),
				CreatedAt:             null.NewTime(currentTimeFromRfc3339.Add(time.Second*0), true),
				UpdatedAt:             null.NewTime(currentTimeFromRfc3339.Add(time.Second*0), true),
				DeletedAt:             null.NewTime(time.Time{}, false),
			},
		},
	}
	return sessionMock
}
