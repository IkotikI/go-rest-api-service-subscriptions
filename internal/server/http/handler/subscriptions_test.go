package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	microservice "github.com/ikotiki/go-rest-api-service-subscriptions"
	mock_service "github.com/ikotiki/go-rest-api-service-subscriptions/internal/service/mocks"
	"github.com/ikotiki/go-rest-api-service-subscriptions/internal/storage"
	"github.com/ikotiki/go-rest-api-service-subscriptions/logger"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var defaultLogLevel = "trace"
var test_time = storage.NewDate(time.Now().Round(24 * time.Hour).UTC())

func createTestRequest(t *testing.T, method string, path string, body interface{}) *http.Request {
	buffer := &bytes.Buffer{}
	err := json.NewEncoder(buffer).Encode(body)
	if err != nil {
		t.Fatalf("error encoding json: %s", err)
	}
	req, _ := http.NewRequest(method, path, buffer)
	return req
}

func Test_createSubscription(t *testing.T) {
	logger.InitLoggerByFlag(defaultLogLevel, true)
	srv := mock_service.NewMockSubscriptions(t)
	router := gin.New()
	g := router.Group("/")
	h := NewSubscriptionHandler(g, srv)
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	type resp struct {
		Success bool
		Obj     int64
		Msg     string
	}
	tests := []struct {
		name    string
		mock    func()
		input   *map[string]interface{}
		want    *resp
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				srv.EXPECT().Create(mock.Anything, mock.Anything).Return(1, nil)
			},
			input: &map[string]interface{}{
				"user_id":       "3d6e2e6c-0d8a-4c1d-9b6f-3b1f9c2b1f9c",
				"service_name":  "test",
				"monthly_price": 100,
				"start_date":    "2020-01-01",
			},
			want:    &resp{Obj: 1, Success: true, Msg: msgSuccess},
			wantErr: false,
		},
		{
			name: "Error (id)",
			mock: func() {
				srv.EXPECT().Create(mock.Anything, mock.Anything).Return(0, errors.New("error binding json: invalid UUID length"))
			},
			input: &map[string]interface{}{
				"user_id":       "3d6e2e6c-0d8a-4c1d-",
				"service_name":  "test",
				"monthly_price": 100,
				"start_date":    "2020-01-01",
			},
			want:    &resp{Obj: 0, Success: false, Msg: "error binding json: invalid UUID length"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(w)

			tt.mock()

			c.Request = createTestRequest(t, "POST", "/subscription", tt.input)
			log.Debug().Interface("request", c.Request).Interface("body", c.Request.Body).Msg("request")

			h.createSubscription(c)

			got := &resp{}
			json.NewDecoder(w.Body).Decode(got)

			log.Debug().Interface("got", got).Msg("got")
			if tt.wantErr {
				assert.Equal(t, tt.want.Success, got.Success)
				assert.Equal(t, tt.want.Obj, got.Obj)
				if tt.want.Msg != "" {
					assert.True(t, strings.Contains(got.Msg, tt.want.Msg))
				}
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}

}

func Test_getSubscriptionByID(t *testing.T) {
	logger.InitLoggerByFlag(defaultLogLevel, true)
	srv := mock_service.NewMockSubscriptions(t)
	router := gin.New()
	g := router.Group("/")
	h := NewSubscriptionHandler(g, srv)
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	testSub := &microservice.Subscription{
		ID:           1,
		UserID:       uuid.MustParse("3d6e2e6c-0d8a-4c1d-9b6f-3b1f9c2b1f9c"),
		ServiceName:  "test",
		MonthlyPrice: 100,
		StartDate:    microservice.NewDate(time.Now()),
	}
	tmp, _ := json.Marshal(testSub)
	testSubJson := map[string]interface{}{}
	json.Unmarshal(tmp, &testSubJson)
	tests := []struct {
		name string
		mock func()
		// input   nil
		want    *resp
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				srv.EXPECT().GetByID(mock.Anything, mock.Anything).Return(testSub, nil)
			},
			want:    &resp{Obj: testSubJson, Success: true, Msg: msgSuccess},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(w)

			tt.mock()

			c.Request = createTestRequest(t, "GET", "/subscription/1", nil)
			c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
			log.Debug().Interface("request", c.Request).Interface("body", c.Request.Body).Msg("request")

			h.getSubscriptionByID(c)

			got := &resp{}
			json.NewDecoder(w.Body).Decode(got)

			log.Debug().Interface("got", got).Msg("got")
			if tt.wantErr {
				assert.Equal(t, tt.want.Success, got.Success)
				assert.Equal(t, tt.want.Obj, got.Obj)
				if tt.want.Msg != "" {
					assert.True(t, strings.Contains(got.Msg, tt.want.Msg))
				}
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}

}

func Test_updateSubscription(t *testing.T) {
	logger.InitLoggerByFlag(defaultLogLevel, true)
	srv := mock_service.NewMockSubscriptions(t)
	router := gin.New()
	g := router.Group("/")
	h := NewSubscriptionHandler(g, srv)
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	tests := []struct {
		name    string
		mock    func()
		input   *map[string]interface{}
		want    *resp
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				srv.EXPECT().Update(mock.Anything, mock.Anything).Return(nil)
			},
			input: &map[string]interface{}{
				"user_id":       "3d6e2e6c-0d8a-4c1d-9b6f-3b1f9c2b1f9c",
				"service_name":  "test",
				"monthly_price": 100,
				"start_date":    "2020-01-01",
				"end_date":      "2020-01-02",
			},
			want:    &resp{Obj: nil, Success: true, Msg: msgSuccess},
			wantErr: false,
		},
		{
			name: "Error (id)",
			mock: func() {
				srv.EXPECT().Update(mock.Anything, mock.Anything).Return(errors.New("error binding json: invalid UUID length"))
			},
			input: &map[string]interface{}{
				"user_id":       "3d6e2e6c-0d8a-4c1d-",
				"service_name":  "test",
				"monthly_price": 100,
				"start_date":    "2020-01-01",
				"end_date":      "2020-01-02",
			},
			want:    &resp{Obj: nil, Success: false, Msg: "error binding json: invalid UUID length"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(w)

			tt.mock()

			c.Request = createTestRequest(t, "POST", "/subscription", tt.input)
			log.Debug().Interface("request", c.Request).Interface("body", c.Request.Body).Msg("request")

			h.updateSubscription(c)

			got := &resp{}
			json.NewDecoder(w.Body).Decode(got)

			log.Debug().Interface("got", got).Msg("got")
			if tt.wantErr {
				assert.Equal(t, tt.want.Success, got.Success)
				assert.Equal(t, tt.want.Obj, got.Obj)
				if tt.want.Msg != "" {
					assert.True(t, strings.Contains(got.Msg, tt.want.Msg))
				}

			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}

}

func Test_deleteSubscription(t *testing.T) {
	logger.InitLoggerByFlag(defaultLogLevel, true)
	srv := mock_service.NewMockSubscriptions(t)
	router := gin.New()
	g := router.Group("/")
	h := NewSubscriptionHandler(g, srv)
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	tests := []struct {
		name    string
		mock    func()
		input   *map[string]interface{}
		want    int
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				srv.EXPECT().DeleteByID(mock.Anything, mock.Anything).Return(nil)
			},
			want: http.StatusNoContent,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(w)

			tt.mock()

			c.Request = createTestRequest(t, "DELETE", "/subscription/1", nil)
			c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

			log.Debug().Interface("request", c.Request).Interface("body", c.Request.Body).Msg("request")

			h.deleteSubscription(c)

			code := w.Result().StatusCode

			log.Debug().Interface("got", code).Msg("got")
			if tt.wantErr {

			} else {
				assert.Equal(t, tt.want, code)
			}
		})
	}

}

func Test_querySubscriptions(t *testing.T) {
	logger.InitLoggerByFlag(defaultLogLevel, true)
	srv := mock_service.NewMockSubscriptions(t)
	router := gin.New()
	g := router.Group("/")
	h := NewSubscriptionHandler(g, srv)
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	test_time_time, _ := time.Parse("2006-01-02", "2020-01-01")
	test_time = microservice.NewDate(test_time_time.UTC())

	time_day := 24 * time.Hour
	subs := []*storage.Subscription{
		{1, uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), "Yandex Taxi", 400, test_time, test_time.Add(2 * time_day)},
		{2, uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"), "Sberbank Shop", 200, test_time.Add(-120 * time_day), test_time.Add(-90 * time_day)},
		{3, uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"), "Ozon Sales", 300, test_time.Add(-60 * time_day), test_time.Add(-30 * time_day)},
	}
	bytes, _ := json.Marshal(subs)
	jsonSubs := []map[string]interface{}{}
	json.Unmarshal(bytes, &jsonSubs)
	tests := []struct {
		name    string
		mock    func()
		input   *map[string]interface{}
		want    *resp
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				srv.EXPECT().Query(mock.Anything, mock.Anything).Return(subs, nil)
			},
			input: nil,
			want:  &resp{Obj: jsonSubs, Success: true, Msg: msgSuccess},
		},
		{
			name: "Ok (date range)",
			mock: func() {
				srv.EXPECT().Query(mock.Anything, mock.Anything).Return(subs[:0], nil)
			},
			input: &map[string]interface{}{
				"start_date": test_time.Add(-10 * time_day).Format("2006-01-02"),
				"end_date":   test_time.Add(10 * time_day).Format("2006-01-02"),
			},
			want: &resp{Obj: jsonSubs[:0], Success: true, Msg: msgSuccess},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(w)

			tt.mock()

			c.Request = createTestRequest(t, "GET", "/subscription/query", tt.input)
			log.Debug().Interface("request", c.Request).Interface("body", c.Request.Body).Msg("request")

			h.querySubscriptions(c)

			got := &resp{}
			json.NewDecoder(w.Body).Decode(got)

			log.Debug().Interface("got", got).Msg("got")
			if tt.wantErr {
				assert.Equal(t, tt.want.Success, got.Success)
				assert.Equal(t, tt.want.Obj, got.Obj)
				if tt.want.Msg != "" {
					assert.True(t, strings.Contains(got.Msg, tt.want.Msg))
				}
			} else {
				actualSlice, ok := got.Obj.([]interface{})
				require.True(t, ok)
				for i, v := range actualSlice {
					assert.Equal(t, v, actualSlice[i])
				}
			}
		})
	}
}

func Test_sumSubscriptions(t *testing.T) {
	logger.InitLoggerByFlag(defaultLogLevel, true)
	srv := mock_service.NewMockSubscriptions(t)
	router := gin.New()
	g := router.Group("/")
	h := NewSubscriptionHandler(g, srv)
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	time_day := time.Hour * 24
	// subs := []*storage.Subscription{
	// 	{1, uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), "Yandex Taxi", 400, test_time, test_time.Add(2 * time_day)},
	// 	{2, uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"), "Sberbank Shop", 200, test_time.Add(-120 * time_day), test_time.Add(-90 * time_day)},
	// 	{3, uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"), "Ozon Sales", 300, test_time.Add(-60 * time_day), test_time.Add(-30 * time_day)},
	// }
	type resp struct {
		Success bool
		Obj     int64
		Msg     string
	}
	tests := []struct {
		name    string
		mock    func()
		input   *map[string]interface{}
		want    *resp
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				srv.EXPECT().Sum(mock.Anything, mock.Anything).Return(900, nil)
			},
			input:   nil,
			want:    &resp{Obj: 900, Success: true, Msg: msgSuccess},
			wantErr: false,
		},
		{
			name: "Error (id)",
			mock: func() {
				srv.EXPECT().Sum(mock.Anything, mock.Anything).Return(400, nil)
			},
			input: &map[string]interface{}{
				"start_date": test_time.Add(-10 * time_day).Format("2006-01-02"),
				"end_date":   test_time.Add(10 * time_day).Format("2006-01-02"),
			},
			want: &resp{Obj: 900, Success: true, Msg: msgSuccess},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(w)

			tt.mock()

			c.Request = createTestRequest(t, "GET", "/subscription/query", tt.input)
			log.Debug().Interface("request", c.Request).Interface("body", c.Request.Body).Msg("request")

			h.sumSubscriptions(c)

			got := &resp{}
			json.NewDecoder(w.Body).Decode(got)

			log.Debug().Interface("got", got).Msg("got")
			if tt.wantErr {
				assert.Equal(t, tt.want.Success, got.Success)
				assert.Equal(t, tt.want.Obj, got.Obj)
				if tt.want.Msg != "" {
					assert.True(t, strings.Contains(got.Msg, tt.want.Msg))
				}
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
