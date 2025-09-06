package postgresql

import (
	"database/sql"
	"fmt"

	"github.com/ikotiki/go-rest-api-service-subscriptions/internal/storage"

	"github.com/ikotiki/sqlbuilder"
	"github.com/ikotiki/sqlbuilder/builder"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

// Database tables list.
const (
	TableSubscriptions string = "subscriptions"
)

// Mapping for abstract storage.QueryArgs to a table name.
var fromToTable = map[storage.From]string{
	storage.FromSubscriptions: TableSubscriptions,
}

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type SQLStorage struct {
	db      *sqlx.DB
	builder *builder.SQLBuilder
}

func NewSQLStorage(cfg Config) (*SQLStorage, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	builder, err := sqlbuilder.NewSQLBuilder(db.DriverName())
	if err != nil {
		return nil, err
	}

	return &SQLStorage{db: db, builder: builder}, nil
}

func (s *SQLStorage) SQLInstance() *sql.DB {
	return s.db.DB
}

func (s *SQLStorage) Close() error {
	return s.db.Close()
}

// func NewPostgresStore(db *SQLStorage) *storage.Storage {
// 	return &storage.Storage{
// 		Subscriptions: NewSubscriptionsStore(db),
// 	}
// }

func sprintf(q string, args ...interface{}) string {
	return fmt.Sprintf(q, args...)
}
