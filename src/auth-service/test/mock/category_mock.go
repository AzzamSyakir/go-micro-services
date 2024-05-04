package mock

import (
	"go-micro-services/src/auth-service/entity"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
)

type CategoryMock struct {
	Data []*entity.Category
}

func NewCategoryMock() *CategoryMock {
	currentTime := time.Now()
	currentTimeRfc3339 := currentTime.Format(time.RFC3339)
	currentTimeFromRfc3339, parseErr := time.Parse(time.RFC3339, currentTimeRfc3339)
	if parseErr != nil {
		panic(parseErr)
	}
	CategoryMock := &CategoryMock{
		Data: []*entity.Category{
			{
				Id:        null.NewString(uuid.NewString(), true),
				Name:      null.NewString(uuid.NewString(), true),
				CreatedAt: null.NewTime(currentTimeFromRfc3339.Add(0*time.Second), true),
				UpdatedAt: null.NewTime(currentTimeFromRfc3339.Add(0*time.Second), true),
				DeletedAt: null.NewTime(time.Time{}, false),
			},
			{
				Id:        null.NewString(uuid.NewString(), true),
				Name:      null.NewString(uuid.NewString(), true),
				CreatedAt: null.NewTime(currentTimeFromRfc3339.Add(1*time.Second), true),
				UpdatedAt: null.NewTime(currentTimeFromRfc3339.Add(1*time.Second), true),
				DeletedAt: null.NewTime(time.Time{}, false),
			},
		},
	}
	return CategoryMock
}
