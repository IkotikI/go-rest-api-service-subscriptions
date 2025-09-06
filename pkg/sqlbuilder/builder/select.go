package builder

type SelectBuilder interface {
	BuildSelect([]Column) (string, []interface{})
	BuildFrom(Table) (string, []interface{})
	BuildWhere([]Where) (string, []interface{})
	BuildGroupBy(GroupBy) (string, []interface{})
	BuildOrderBy([]OrderBy) (string, []interface{})
	BuildLimit(Limit) (string, []interface{})
}

type SelectArguments struct {
	Arguments

	Select  []Column
	From    Table
	Where   []Where
	GroupBy GroupBy
	OrderBy []OrderBy
	Limit   Limit
}

func (a *SelectArguments) PartOrder() []string {
	return []string{"select", "from", "where", "group_by", "order_by", "limit"}
}

func (a *SelectArguments) BuildPartByName(partName string, b Builder) (string, []interface{}) {
	switch partName {
	case "select":
		return b.BuildSelect(a.Select)
	case "from":
		return b.BuildFrom(a.From)
	case "where":
		return b.BuildWhere(a.Where)
	case "group_by":
		return b.BuildGroupBy(a.GroupBy)
	case "order_by":
		return b.BuildOrderBy(a.OrderBy)
	case "limit":
		return b.BuildLimit(a.Limit)
	default:
		return "", nil
	}
}
