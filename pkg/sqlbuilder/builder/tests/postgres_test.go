package sqlbuilder_test

import (
	"strings"
	"testing"

	"github.com/ikotiki/sqlbuilder"
	"github.com/ikotiki/sqlbuilder/builder"
)

func TestPostgreSQL(t *testing.T) {
	// Test 1
	args := &builder.SelectArguments{
		Select: []builder.Column{"id", "name"},
		From:   "users",
		Where: []builder.Where{
			{Column: "id", Operator: "=", Value: "5"},
			{Column: "name", Operator: ">=", Value: "Juan"},
		},
		GroupBy: "id",
		OrderBy: []builder.OrderBy{
			{Column: "id", Order: "DESC"},
			{Column: "name", Order: "ASC"},
		},
		Limit: builder.Limit{Offset: 10, Limit: 10},
	}

	b, err := sqlbuilder.NewSQLBuilder("postgres")
	if err != nil {
		t.Fatalf("Can't create sqlite builder: %v", err)
	}

	parts := []string{"select", "from", "where", "group_by", "order_by", "limit"}
	q, qargs := b.BuildParts(parts, args)

	t.Logf("q: %s, args: %v", q, qargs)

	if strings.Count(q, "?") != len(qargs) {
		t.Errorf("Query should contain %d placeholders, but got %d", len(qargs), strings.Count(q, "?"))
	}
	// Test 2
	args = &builder.SelectArguments{
		From: "users",
		Where: []builder.Where{
			{Column: "id", Operator: ">", Value: "5"},
			{Column: "name", Operator: "=", Value: "Juan"},
		},
		GroupBy: "id",
		OrderBy: []builder.OrderBy{
			{Column: "id"},
			{Column: "name"},
		},
		Limit: builder.Limit{Limit: 10},
	}

	q, qargs = b.BuildParts(parts, args)

	t.Logf("q: %s, args: %v", q, qargs)

	if strings.Count(q, "?") != len(qargs) {
		t.Errorf("Query should contain %d placeholders, but got %d", len(qargs), strings.Count(q, "?"))
	}

}

func TestPostgreSQLInsert(t *testing.T) {
	b, err := sqlbuilder.NewSQLBuilder("postgres")
	if err != nil {
		t.Fatalf("Can't create sqlite builder: %v", err)
	}

	args := &builder.InsertArguments{
		Into:    "users",
		Columns: []builder.Column{"id", "name"},
		Values:  []builder.Value{"5", "Juan"},
		Where:   []builder.Where{{Column: "id", Operator: "=", Value: "5"}},
	}

	parts := args.PartOrder()

	q, qargs := b.BuildParts(parts, args)

	t.Logf("q: %s, args: %v", q, qargs)

	if strings.Count(q, "?") != len(qargs) {
		t.Errorf("Query should contain %d placeholders, but got %d", len(qargs), strings.Count(q, "?"))
	}

}

// func TestMeme(t *testing.T) {
// 	a := make([]string, 501)
// 	for i := 1; i <= len(a)-1; i++ {
// 		a[i] = strconv.Itoa(i)
// 	}
// 	fmt.Printf("[%v]", strings.Join(a, ", ")[2:])
// }
