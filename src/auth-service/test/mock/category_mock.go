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
	categoryMock := &CategoryMock{
		Data: []*entity.Category{
			{
				Id:        null.NewString(uuid.NewString(), true),
				Name:      null.NewString("name0", true),
				CreatedAt: null.NewTime(currentTimeFromRfc3339.Add(0*time.Second), true),
				UpdatedAt: null.NewTime(currentTimeFromRfc3339.Add(time.Duration(time.Duration.Seconds(0))), true),
				DeletedAt: null.NewTime(time.Time{}, false),
			},
			{
				Id:        null.NewString(uuid.NewString(), true),
				Name:      null.NewString("name1", true),
				CreatedAt: null.NewTime(currentTimeFromRfc3339.Add(0*time.Second), true),
				UpdatedAt: null.NewTime(currentTimeFromRfc3339.Add(time.Duration(time.Duration.Seconds(1))), true),
				DeletedAt: null.NewTime(time.Time{}, false),
			},
		},
	}
	return categoryMock
}
