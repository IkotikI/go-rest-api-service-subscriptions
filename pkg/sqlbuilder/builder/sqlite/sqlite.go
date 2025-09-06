package sqlite

import (
	"fmt"
	"strings"

	"github.com/ikotiki/sqlbuilder/builder"
)

var spf = fmt.Sprintf

type SQLiteBuilder struct {
}

func NewSQLiteBuilder() *SQLiteBuilder {
	return &SQLiteBuilder{}
}

func (b *SQLiteBuilder) BuildSelect(s []builder.Column) (string, []interface{}) {
	if len(s) == 0 {
		return "SELECT *", nil
	}
	str, args := joinWithCommas(s)
	return spf("SELECT %s", str), args
}

func (b *SQLiteBuilder) BuildFrom(f builder.Table) (string, []interface{}) {
	if f == "" {
		return "", nil
	}
	return spf("FROM `%s`", f), []interface{}{}
}

func (b *SQLiteBuilder) BuildWhere(w []builder.Where) (string, []interface{}) {
	if len(w) == 0 {
		return "", nil
	}
	size := len(w)
	str := make([]string, size)
	args := make([]interface{}, size*2)
	for i, where := range w {
		if where.Column == "" || where.Operator == "" || where.Value == nil || where.Value == "" {
			size--
			continue
		}
		str[i] = spf("`%s` %s ?", where.Column, where.Operator)
		args[i] = where.Value
	}
	// Reduce slice sizes
	str = str[:size]
	args = args[:size*2]
	return spf("WHERE %s", strings.Join(str, " AND ")), args
}

func (b *SQLiteBuilder) BuildGroupBy(g builder.GroupBy) (string, []interface{}) {
	if g == "" {
		return "", nil
	}
	return "GROUP BY ?", []interface{}{g}
}

func (b *SQLiteBuilder) BuildOrderBy(o []builder.OrderBy) (string, []interface{}) {
	if len(o) == 0 {
		return "", nil
	}
	size := len(o)
	str := make([]string, size)
	for i, order := range o {
		if order.Column == "" || order.Order == "" {
			size--
			continue
		}
		str[i] = spf("`%s` %s", order.Column, order.Order)
	}
	// Reduce slice sizes
	str = str[:size]
	return spf("ORDER BY %s", strings.Join(str, ", ")), []interface{}{}
}

func (b *SQLiteBuilder) BuildLimit(l builder.Limit) (string, []interface{}) {
	if l.Limit <= 0 {
		return "", nil
	}
	limit := "LIMIT ?"
	if l.Offset <= 0 {
		return limit, []interface{}{l.Limit}
	}
	offset := "OFFSET ?"
	return spf("%s %s", limit, offset), []interface{}{l.Limit, l.Offset}
}

func (b *SQLiteBuilder) BuildInsertInto(t builder.Table) (string, []interface{}) {
	if len(t) == 0 {
		return "", nil
	}
	return "INSERT INTO ?", []interface{}{t}
}

func (b *SQLiteBuilder) BuildInsertColumns(c []builder.Column) (string, []interface{}) {
	if len(c) == 0 {
		return "", nil
	}
	str, args := joinWithCommas(c)
	return fmt.Sprintf("(%s)", str), args
}

func (b *SQLiteBuilder) BuildValues(t []builder.Value) (string, []interface{}) {
	if len(t) == 0 {
		return "", nil
	}
	str, args := joinWithCommas(t)
	return fmt.Sprintf("VALUES (%s)", str), args
}

func joinWithCommas[T ~string](s []T) (string, []interface{}) {
	str, _ := strings.CutPrefix(strings.Repeat(",?", len(s)), ",")
	args := make([]interface{}, len(s))
	for i, v := range s {
		args[i] = v
	}
	return str, args
}
