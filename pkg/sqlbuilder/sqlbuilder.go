package sqlbuilder

import (
	"fmt"

	"github.com/ikotiki/sqlbuilder/builder"
	"github.com/ikotiki/sqlbuilder/builder/postgres"
	"github.com/ikotiki/sqlbuilder/builder/sqlite"
)

var drivers = []string{"sqlite3", "postgres"}

// var ErrNoSuchDriver = errors.New("no such driver")

func Drivers() []string {
	return drivers
}

func NewSQLBuilder(driver string) (*builder.SQLBuilder, error) {
	switch driver {
	case "sqlite3":
		return &builder.SQLBuilder{Driver: driver, Builder: sqlite.NewSQLiteBuilder()}, nil
	case "postgres":
		return &builder.SQLBuilder{Driver: driver, Builder: postgres.NewPostgresSQLBuilder()}, nil
	default:
		return nil, fmt.Errorf("sqlbuilder: unsupported driver %s", driver)
	}
}
