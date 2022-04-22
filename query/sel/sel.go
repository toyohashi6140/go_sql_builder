package sel

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/toyohashi6140/go_sql_builder/column"
	"github.com/toyohashi6140/go_sql_builder/query"
	"github.com/toyohashi6140/go_sql_builder/table"
	"github.com/toyohashi6140/go_sql_builder/where"
)

type (
	sel struct {
		columns     column.Columns
		table       *table.Table
		filter      *where.Conditions
		groupby     column.Columns
		orderby     *column.Order
		skip, limit int
	}

	// SelectSetter An interface that collects setters for setting arbitrary parameters such as where filter and group by.
	SelectSetter interface {
		// Filter where clause
		Filter(lo int, wheres ...*where.Condition) SelectSetter

		// GroupBy group by columns(group by clause)
		GroupBy(cols ...*column.Column) SelectSetter

		// OrderBy order by columns(order by clause)
		OrderBy(desc bool, cols ...*column.Column) SelectSetter

		// Skip the number of offset
		Skip(int) SelectSetter

		// Limit the number of getting column(limit clause)
		Limit(int) SelectSetter

		// ToQuery By executing this function at any time, it is possible to build and execute SQL as a Query type.
		ToQuery() query.Query
	}
)

// New Calling this returns a Select Setter type. If you don't need to set other items (such as where), you need to call ToQuery and make it executable.
func New(cols column.Columns, tbl *table.Table) SelectSetter {
	return &sel{columns: cols, table: tbl}
}

// Filter Calling this sets a where filter and returns a Select Setter type. If you don't need to set any more items, you need to call ToQuery and make it executable.
func (s *sel) Filter(lo int, wheres ...*where.Condition) SelectSetter {
	s.filter = where.NewConditions(lo, wheres...)
	return s
}

// GroupBy Calling this sets a group by columns and returns a SelectSetter type. If you don't need to set any more items, you need to call ToQuery and make it executable.
func (s *sel) GroupBy(cols ...*column.Column) SelectSetter {
	s.groupby = cols
	return s
}

// OrderBy Calling this sets a where order by columns and returns a SelectSetter type. If you don't need to set any more items, you need to call ToQuery and make it executable.
func (s *sel) OrderBy(desc bool, cols ...*column.Column) SelectSetter {
	s.orderby = column.NewOrder(desc, cols...)
	return s
}

// Skip Calling this sets a offset number and returns a SelectSetter type. if limit is 0 or not set, this setting is ignored.
func (s *sel) Skip(skip int) SelectSetter {
	s.skip = skip
	return s
}

// Limit Calling this sets a limit number and returns a SelectSetter type. If you don't need to set any more items, you need to call ToQuery and make it executable.
func (s *sel) Limit(l int) SelectSetter {
	s.limit = l
	return s
}

// ToQuery makes structs buildable and executable
func (s *sel) ToQuery() query.Query {
	return s
}

// Build build SQL of SELECT sentence which is structured from SELECT and FROM. if you set "true" to "or", each filter works as OR
func (s *sel) Build() (string, error) {
	if len(s.columns) == 0 {
		return "", errors.New("no select columns. columns must need 1 at least")
	}
	if s.table.TName == "" {
		return "", errors.New("no table selected")
	}
	queries := []string{
		fmt.Sprintf("SELECT %s", strings.Join(s.columns.Line(), ", ")),
		fmt.Sprintf("FROM %s", s.table.Line()),
	}
	if s.filter != nil {
		queries = append(queries, fmt.Sprintf("WHERE %s", s.filter.Join()))
	}
	if s.groupby != nil {
		queries = append(queries, fmt.Sprintf("GROUP BY %s", s.groupby.NoAliasName()))
	}
	if s.orderby != nil {
		if s.orderby.Desc() {
			queries = append(queries, fmt.Sprintf("ORDER BY %s %s", s.orderby.Columns().Name(), "DESC"))
		} else {
			queries = append(queries, fmt.Sprintf("ORDER BY %s", s.orderby.Columns().Name()))
		}
	}
	if s.limit > 0 {
		if s.skip > 0 {
			queries = append(queries, fmt.Sprintf("LIMIT %d, %d", s.skip, s.limit))
		} else {
			queries = append(queries, fmt.Sprintf("LIMIT %d", s.limit))
		}
	}
	return strings.Join(queries, " "), nil
}

// Bind returns bind-values "?" in SQL.
func (s *sel) Bind() query.Bind {
	return s.filter.Bind()
}

// Execute Build and Bind, and execute the query based on the SQL and Value created by these.
func (s *sel) Execute(db *sql.DB, sql string, bind ...interface{}) (interface{}, error) {
	rows, err := db.Query(sql, bind...)
	if err != nil {
		return nil, err
	}
	rowVals := []map[string]interface{}{}
	for rows.Next() {
		colVals := []interface{}{}
		for _, v := range s.columns {
			colVals = append(colVals, &v.Value)
		}
		err := rows.Scan(colVals...)
		if err != nil {
			return nil, err
		}

		m := map[string]interface{}{}
		for _, c := range s.columns {
			m[c.Name()] = c.Value
		}
		rowVals = append(rowVals, m)
	}
	return rowVals, nil
}
