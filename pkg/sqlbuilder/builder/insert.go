package builder

// Part of builder interface
type InsertBuilder interface {
	BuildInsertInto(Table) (string, []interface{})
	BuildInsertColumns([]Column) (string, []interface{})
	BuildValues([]Value) (string, []interface{})
	BuildWhere([]Where) (string, []interface{})
}

type InsertArguments struct {
	Arguments

	Into    Table
	Columns []Column
	Values  []Value
	Where   []Where
}

func (a *InsertArguments) PartOrder() []string {
	return []string{"insert", "column", "value", "into", "where"}
}

func (a *InsertArguments) BuildPartByName(partName string, b Builder) (string, []interface{}) {
	switch partName {
	case "insert":
		return b.BuildInsertInto(a.Into)
	case "columns", "column":
		return b.BuildInsertColumns(a.Columns)
	case "values", "value":
		return b.BuildValues(a.Values)
	case "where":
		return b.BuildWhere(a.Where)
	default:
		return "", nil
	}
}
