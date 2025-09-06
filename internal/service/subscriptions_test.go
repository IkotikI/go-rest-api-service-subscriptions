package service

import (
	"testing"
	"time"

	"github.com/ikotiki/go-rest-api-service-subscriptions/internal/storage"
	"github.com/ikotiki/go-rest-api-service-subscriptions/logger"
	"github.com/stretchr/testify/assert"
)

func Test_parseQueryArgs(t *testing.T) {
	logger.InitLoggerByFlag("trace", false)
	srv := &SubscriptionService{}

	testTimeStrStart := "2020-01-01"
	testTimeStrEnd := "2020-02-01"
	testTimeStart, _ := time.Parse("2006-01-02", testTimeStrStart)
	testTimeEnd, _ := time.Parse("2006-01-02", testTimeStrEnd)
	tests := []struct {
		name    string
		mock    func()
		input   *SubscriptionQueryArgs
		want    *storage.QueryArgs
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {},
			input: &SubscriptionQueryArgs{
				UserID:    "123e4567-e89b-12d3-a456-426614174000",
				StartDate: testTimeStrStart,
				EndDate:   testTimeStrEnd,
				Order: []Order{
					{
						OrderBy: "service_name",
						Order:   "DESC",
					},
					{
						OrderBy: "user_id",
						Order:   "ASC",
					},
				},
			},
			want: &storage.QueryArgs{
				Where: []storage.Where{
					{
						Column:   "user_id",
						Operator: storage.OpEqual,
						Value:    "123e4567-e89b-12d3-a456-426614174000",
					},
					{
						Column:   "start_date",
						Operator: storage.OpMoreOrEqual,
						Value:    testTimeStart,
					},
					{
						Column:   "end_date",
						Operator: storage.OpLessOrEqual,
						Value:    testTimeEnd,
					},
				},
				Order: []storage.OrderStruct{
					{
						OrderBy: "service_name",
						Order:   "DESC",
					},
					{
						OrderBy: "user_id",
						Order:   "ASC",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := srv.parseQueryArgs(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			// assert.NoError(t)
		})
	}
}
