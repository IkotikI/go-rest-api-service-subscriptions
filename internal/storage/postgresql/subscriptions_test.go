package postgresql

import (
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/ikotiki/sqlbuilder"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"github.com/ikotiki/go-rest-api-service-subscriptions/internal/storage"
	"github.com/ikotiki/go-rest-api-service-subscriptions/logger"
)

func newMockSQLStorage() (store *SQLStorage, db *sqlx.DB, mock sqlmock.Sqlmock, err error) {
	logger.InitLoggerByFlag("trace", true)

	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	// mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		return nil, nil, nil, err
	}
	db = sqlx.NewDb(mockDB, "sqlmock")

	if err := db.Ping(); err != nil {
		return nil, nil, nil, err
	}

	builder, err := sqlbuilder.NewSQLBuilder("postgres")
	if err != nil {
		return nil, nil, nil, err
	}

	return &SQLStorage{db: db, builder: builder}, db, mock, nil
}

var test_time = storage.NewDate(time.Now())

func TestSubscriptions_Create(t *testing.T) {
	dbStore, db, mock, err := newMockSQLStorage()
	if err != nil {
		t.Fatal("can't crate mock storage", err)
	}
	defer db.Close()
	st := NewSubscriptionsStore(dbStore)

	tests := []struct {
		name    string
		mock    func()
		input   *storage.Subscription
		want    storage.SubscriptionID
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO subscriptions (user_id, service_name, monthly_price, start_date, end_date) VALUES ($1, $2, $3, $4, $5) RETURNING id").
					WillReturnRows(rows)
			},
			input: &storage.Subscription{
				UserID:       uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				ServiceName:  "Yandex Taxi",
				MonthlyPrice: 400,
				StartDate:    test_time,
				EndDate:      test_time.Add(time.Hour),
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := st.Create(t.Context(), tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
func TestSubscriptions_GetByID(t *testing.T) {
	dbStore, db, mock, err := newMockSQLStorage()
	if err != nil {
		t.Fatal("can't crate mock storage", err)
	}
	defer db.Close()
	st := NewSubscriptionsStore(dbStore)

	tests := []struct {
		name    string
		mock    func()
		input   storage.SubscriptionID
		want    *storage.Subscription
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "service_name", "monthly_price", "start_date", "end_date"}).AddRow(
					1,
					uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
					"Yandex Taxi",
					400,
					test_time,
					test_time.Add(time.Hour))

				mock.ExpectQuery("SELECT * FROM subscriptions WHERE id = $1").
					WillReturnRows(rows)
			},
			input: 1,
			want: &storage.Subscription{
				ID:           1,
				UserID:       uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				ServiceName:  "Yandex Taxi",
				MonthlyPrice: 400,
				StartDate:    test_time,
				EndDate:      test_time.Add(time.Hour),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := st.GetByID(t.Context(), tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSubscriptions_Update(t *testing.T) {
	dbStore, db, mock, err := newMockSQLStorage()
	if err != nil {
		t.Fatal("can't crate mock storage", err)
	}
	defer db.Close()
	st := NewSubscriptionsStore(dbStore)

	type args struct {
		item *storage.Subscription
	}
	tests := []struct {
		name    string
		mock    func()
		input   *storage.Subscription
		want    storage.SubscriptionID
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				// rows := sqlmock.NewRows([]string{"id", "user_id", "service_name", "monthly_price", "start_date", "end_date"}).AddRow(
				// 	1,
				// 	uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				// 	"Yandex Taxi",
				// 	400,
				// 	test_time,
				// 	test_time.Add(4*time.Hour))

				mock.ExpectExec(`UPDATE subscriptions SET (service_name, monthly_price, start_date, end_date) = ($2, $3, $4, $5) WHERE id = $1`).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			input: &storage.Subscription{
				ID:           1,
				UserID:       uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				ServiceName:  "Yandex Taxi",
				MonthlyPrice: 400,
				StartDate:    test_time.Add(time.Hour),
				EndDate:      test_time.Add(4 * time.Hour),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := st.Update(t.Context(), tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSubscriptions_Query(t *testing.T) {
	dbStore, db, mock, err := newMockSQLStorage()
	if err != nil {
		t.Fatal("can't crate mock storage", err)
	}
	defer db.Close()
	st := NewSubscriptionsStore(dbStore)

	time_day := time.Hour * 24
	subs := []*storage.Subscription{
		{1, uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), "Yandex Taxi", 400, test_time, test_time.Add(2 * time_day)},
		{2, uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"), "Sberbank Shop", 200, test_time.Add(-120 * time_day), test_time.Add(-90 * time_day)},
		{3, uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"), "Ozon Sales", 360, test_time.Add(-60 * time_day), test_time.Add(-30 * time_day)},
	}

	tests := []struct {
		name    string
		mock    func()
		input   *storage.QueryArgs
		want    []*storage.Subscription
		wantErr bool
	}{
		{
			name: "Ok (All)",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "service_name", "monthly_price", "start_date", "end_date"})
				for _, sub := range subs {
					rows.AddRow(sub.ID, sub.UserID, sub.ServiceName, sub.MonthlyPrice, sub.StartDate, sub.EndDate)
				}
				mock.ExpectQuery("SELECT * FROM subscriptions").
					WillReturnRows(rows)
			},
			input: &storage.QueryArgs{},
			want:  subs,
		},
		{
			name: "Ok (Date range)",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "service_name", "monthly_price", "start_date", "end_date"})
				sub := subs[0]
				rows.AddRow(sub.ID, sub.UserID, sub.ServiceName, sub.MonthlyPrice, sub.StartDate, sub.EndDate)

				mock.ExpectQuery("SELECT * FROM subscriptions WHERE (start_date >= $1) AND (end_date IS NULL AND end_date <= $2)").
					WillReturnRows(rows)
			},
			input: &storage.QueryArgs{
				Where: []storage.Where{
					{"start_date", ">=", test_time.Add(-10 * time_day)},
					{"end_date", "<=", test_time.Add(10 * time_day)},
				},
			},
			want: []*storage.Subscription{subs[0]},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := st.Query(t.Context(), tt.input)
			t.Logf("got: %+v", got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSubscriptions_Sum(t *testing.T) {
	dbStore, db, mock, err := newMockSQLStorage()
	if err != nil {
		t.Fatal("can't crate mock storage", err)
	}
	defer db.Close()
	st := NewSubscriptionsStore(dbStore)

	time_day := time.Hour * 24
	// subs := []*storage.Subscription{
	// 	{1, uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), "Yandex Taxi", 400, test_time, test_time.Add(2 * time_day)},
	// 	{2, uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"), "Sberbank Shop", 200, test_time.Add(-120 * time_day), test_time.Add(-90 * time_day)},
	// 	{3, uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"), "Ozon Sales", 300, test_time.Add(-60 * time_day), test_time.Add(-30 * time_day)},
	// }

	tests := []struct {
		name    string
		mock    func()
		input   *storage.QueryArgs
		want    storage.Price
		wantErr bool
	}{
		{
			name: "Ok (All)",
			mock: func() {
				rows := sqlmock.NewRows([]string{"sum"})
				rows.AddRow(900)

				mock.ExpectQuery("SELECT sum(monthly_price) FROM subscriptions AS sum").
					WillReturnRows(rows)
			},
			input: &storage.QueryArgs{},
			want:  900,
		},
		{
			name: "Ok (Date range)",
			mock: func() {
				rows := sqlmock.NewRows([]string{"sum"})
				rows.AddRow(700)

				mock.ExpectQuery("SELECT sum(monthly_price) FROM subscriptions AS sum WHERE (start_date >= $1) AND (end_date IS NULL AND end_date <= $2)").
					WillReturnRows(rows)
			},
			input: &storage.QueryArgs{
				Where: []storage.Where{
					{"start_date", ">=", test_time.Add(-100 * time_day)},
					{"end_date", "<=", test_time.Add(10 * time_day)},
				},
			},
			want: 700,
		},
		{
			name: "Error (Date range)",
			mock: func() {
				rows := sqlmock.NewRows([]string{"sum"})
				rows.AddRow(0)

				mock.ExpectQuery("SELECT sum(monthly_price) FROM subscriptions AS sum WHERE (start_date >= $1) AND (end_date IS NULL AND end_date <= $2)").
					WillReturnRows(rows)

			},
			input: &storage.QueryArgs{
				Where: []storage.Where{
					{"start_date", ">=", test_time.Add(20 * time_day)},
					{"end_date", "<=", test_time.Add(100 * time_day)},
				},
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := st.Sum(t.Context(), tt.input)
			t.Logf("got: %+v", got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
