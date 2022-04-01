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

type sel struct {
	columns column.Columns
	table   table.Table
	filter  where.Conditions
	groupby column.Columns
	orderby column.Columns
}

func New(cols column.Columns, tbl table.Table, wheres ...*where.Condition) query.Query {
	return &sel{columns: cols, table: tbl, filter: wheres}
}

func SetGroupbyColumn(q query.Query, cols ...*column.Column) error {
	sel, ok := q.(*sel)
	if !ok {
		return errors.New("assertion type query to *sel is faild")
	}
	sel.groupby = cols
	return nil
}

func SetOrderByColumn(q query.Query, cols ...*column.Column) error {
	sel, ok := q.(*sel)
	if !ok {
		return errors.New("assertion type query to *sel is faild")
	}
	sel.orderby = cols
	return nil
}

// Build build SQL of SELECT sentence which is structured from SELECT and FROM. if you set "true" to "or", each filter works as OR
func (s *sel) Build(or bool) (string, error) {
	if len(s.columns) == 0 {
		return "", errors.New("no select columns. columns must need 1 at least")
	}
	if s.table.TName == "" {
		return "", errors.New("no table selected")
	}

	queries := []string{
		fmt.Sprintf("SELECT %s", strings.Join(s.columns.Line(), " ,")),
		fmt.Sprintf("FROM %s", s.table.Line()),
	}
	if len(s.filter) > 0 {
		queries = append(queries, fmt.Sprintf("WHERE %s", s.filter.Join(or)))
	}
	if len(s.groupby) > 0 {
		queries = append(queries, fmt.Sprintf("GROUP BY %s", s.groupby.Name()))
	}
	if len(s.orderby) > 0 {
		queries = append(queries, fmt.Sprintf("ORDER BY %s", s.orderby.Name()))
	}
	return strings.Join(queries, " "), nil
}

// Bind returns bind-values "?" in SQL.
func (s *sel) Bind() query.Bind {
	bs := query.Bind{}
	for _, f := range s.filter {
		b := f.Bind()
		bs = append(bs, b...)
	}
	return bs
}

// Execute Build and Bind, and execute the query based on the SQL and Value created by these.
func (s *sel) Execute(db *sql.DB) (interface{}, error) {
	sql, err := s.Build(false)
	if err != nil {
		return nil, err
	}
	bindVals := s.Bind()
	rows, err := db.Query(sql, bindVals...)
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
