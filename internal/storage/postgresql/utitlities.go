package postgresql

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/ikotiki/go-rest-api-service-subscriptions/internal/storage"
	"github.com/rs/zerolog/log"

	"github.com/ikotiki/sqlbuilder/builder"
)

// func (s *SQLStorage) Count(ctx context.Context, args *storage.QueryArgs) (n int64, err error) {
// 	q := `
// 		SELECT count(*) AS count
// 	`

// 	var queryEnd string
// 	var queryArgs []interface{}
// 	queryEnd, queryArgs = s.buildParts([]string{"from", "where"}, args)
// 	q += queryEnd

// 	log.Printf("query: `%s` args: %+v\n", q, queryArgs)
// 	count := &[]int64{}
// 	err = s.db.SelectContext(ctx, count, q, queryArgs...)
// 	if err != nil {
// 		return -1, err
// 	}

// 	fmt.Printf("got count %+v", count)

// 	return (*count)[0], nil
// }

type postgresSQLBuilder struct {
	Builder *builder.SQLBuilder
}

func (s *postgresSQLBuilder) parseQueryArgs(args *storage.QueryArgs) *builder.SelectArguments {
	selectArgs := builder.SelectArguments{
		From: builder.Table(args.From),
		Limit: builder.Limit{
			Offset: args.Offset,
			Limit:  args.Limit,
		},
	}

	// Orders
	orders := make([]builder.OrderBy, 0, len(args.Order))
	for _, o := range args.Order {
		orders = append(orders, builder.OrderBy{
			Column: o.OrderBy,
			Order:  string(o.Order),
		})
	}
	selectArgs.OrderBy = orders

	// Wheres ('AND' joined)
	where := make([]builder.Where, 0, len(args.Where))
	for _, w := range args.Where {
		where = append(where, builder.Where{
			Column:   w.Column,
			Operator: string(w.Operator),
			Value:    w.Value,
		})
	}
	selectArgs.Where = where

	return &selectArgs
}

func (s *postgresSQLBuilder) buildParts(parts []string, args *storage.QueryArgs, startIndex int) (query string, queryArgs []interface{}) {
	log.Trace().Msgf("args: %+v\n", args)
	if args == nil {
		return "", []interface{}{}
	}
	builderArgs := s.parseQueryArgs(args)
	log.Trace().Msgf("parseQueryArgs: args: %+v\n", builderArgs)
	qStr, qArgs := s.Builder.BuildParts(parts, builderArgs)
	log.Trace().Msgf("query: `%s` args: %+v\n", qStr, qArgs)
	qStr = s.replacePlaceholders(qStr, startIndex)
	log.Trace().Msgf("replaced placeholders: `%s` args: %+v\n", qStr, qArgs)
	return qStr, qArgs
}

func (s *postgresSQLBuilder) buildWhere(args *storage.QueryArgs) (whereStr string, whereArgs []interface{}) {
	// Custom handling for where statement
	where := []string{}
	queryArgs := []interface{}{}
	i := 1
	for _, w := range args.Where {
		// log.Debug().Msgf("where: %v", where)
		switch w.Column {
		case "end_date":
			where = append(where, sprintf(`(end_date IS NOT NULL AND end_date <= $%d)`, i))
			log.Debug().Msgf("hasColumn(args.Where, end_date): %v", hasColumn(args.Where, "end_date"))
		default:
			where = append(where, sprintf(`(%s %s $%d)`, w.Column, w.Operator, i))
		}
		queryArgs = append(queryArgs, w.Value)
		i++
	}
	if len(where) > 0 {
		whereStr = sprintf("WHERE %s ", strings.Join(where, " AND "))
	}

	return whereStr, queryArgs
}

func hasColumn(where []storage.Where, column string) bool {
	for _, w := range where {
		if w.Column == column {
			return true
		}
	}
	return false
}

// Replace '? ? ?' with '$1 $2 $3'
func (s *postgresSQLBuilder) replacePlaceholders(query string, startIndex int) string {
	re := regexp.MustCompile(`\?`)
	i := -1 // Because increment is first
	return re.ReplaceAllStringFunc(query, func(s string) string {
		i++
		return "$" + strconv.Itoa(startIndex+i)
	})
}

func ptr[T any](x T) *T {
	return &x
}

// func (s *SQLStorage) MakePagination(ctx context.Context, db storage.Storage, table storage.Table, queryArgs *storage.QueryArguments, args *storage.QueryArgs) (storage.Pagination, error) {
// 	if queryArgs == nil {
// 		queryArgs = storage.DefaultQueryArguments
// 	}
// 	if args == nil {
// 		args = queryArgs.ToQueryArgs()
// 	}
// 	args.From = storage.Table(table)

// 	n, err := db.CountWithBuilder(ctx, args)
// 	if err != nil {
// 		return storage.Pagination{}, err
// 	}
// 	total_pages := n / queryArgs.PerPage
// 	if n-total_pages > 0 {
// 		total_pages += 1
// 	}
// 	return storage.Pagination{
// 		Table:        table,
// 		RecordsCount: n,
// 		TotalPages:   total_pages,
// 		Page:         queryArgs.Page,
// 		PerPage:      queryArgs.PerPage,
// 	}, nil
// }
