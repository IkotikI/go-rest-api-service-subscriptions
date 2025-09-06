package builder

import (
	"fmt"
	"strings"
)

// var drivers = []string{"sqlite3"}

// var ErrNoSuchDriver = errors.New("no such driver")
// func Drivers() []string {
// 	dest := make([]string, len(drivers))
// 	copy(dest, drivers)
// 	return dest
// }

type Builder interface {
	// Build(Arguments) (string, []interface{})
	// BuildParts([]string, Arguments) (string, []interface{})
	SelectBuilder
	InsertBuilder
}

type SQLBuilder struct {
	Driver  string
	Builder Builder
}

// func NewSQLBuilder(driver string, b Builder) (*SQLBuilder, error) {
// 	return &SQLBuilder{Driver: driver, Builder: b}, nil
// }

type Arguments interface {
	BuildPartByName(partName string, b Builder) (q string, args []interface{})
	PartOrder() []string
}

func (b *SQLBuilder) BuildParts(parts []string, a Arguments) (string, []interface{}) {
	if a == nil {
		return "", nil
	}
	queryParts := make([]string, 0, len(parts))
	sqlArgs := make([]interface{}, 0)
	for _, partPositionedName := range a.PartOrder() {
		for _, partGivenName := range parts {
			if partGivenName == partPositionedName || partGivenName == partPositionedName+"s" {
				queryString, tempArgs := a.BuildPartByName(partGivenName, b.Builder)
				if queryString == "" {
					continue
				}
				queryParts = append(queryParts, queryString)
				for _, tempArg := range tempArgs {
					if tempArg != nil {
						sqlArgs = append(sqlArgs, tempArg)
					}
				}
			}
		}
	}

	return strings.Join(queryParts, " "), sqlArgs
}

func (b *SQLBuilder) Build(a Arguments) (string, []interface{}) {
	if a == nil {
		return "", nil
	}
	queryParts := make([]string, len(a.PartOrder()))
	sqlArgs := make([]interface{}, 0)
	for _, partName := range a.PartOrder() {
		queryString, tempArgs := a.BuildPartByName(partName, b.Builder)
		if queryString == "" {
			continue
		}
		queryParts = append(queryParts, queryString)
		for _, tempArg := range tempArgs {
			if tempArg != nil {
				sqlArgs = append(sqlArgs, tempArg)
			}
		}
	}
	return strings.Join(queryParts, " "), sqlArgs
}

// Support types *SelectArguments, *InsertArguments
func (b *SQLBuilder) ValidateArgs(args interface{}) error {
	switch args.(type) {
	case *SelectArguments, *InsertArguments:
		return nil
	default:
		return fmt.Errorf("unsupported arguments type: %T", args)
	}
}

/* ---- Basic types ---- */
type Column string

type Table string

type Value string

type Where struct {
	Column   string
	Operator string
	Value    interface{}
}

type Limit struct {
	Offset int64
	Limit  int64
}

type GroupBy string

type OrderBy struct {
	Column string
	Order  string
}

/* ---- Parse ---- */

// func (c *Column) ParseURLValues(args []string, default []string) {

// 	if queryValue, ok := queryArgs[key]; ok {
// 		v, err := strconv.Atoi(queryValue[0])
// 			if err == nil {
// 				args.Limit.Limit = v
// 				continue
// 			}
// 		}
// 	}

// 	args.Limit.Limit = defaultArgs[key][0]

// }
