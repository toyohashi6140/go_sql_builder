package query

import (
	"errors"
	"fmt"
	"go_batch/internal/database/query/sql/column"
	"go_batch/internal/database/query/sql/table"
	"go_batch/internal/database/query/sql/where"
	"strings"
)

type bind []interface{}

type Query interface {
	Build(bool) (string, error)
	Bind() bind
}

type sel struct {
	columns column.Columns
	table   table.Table
	filter  where.Conditions
	groupby column.Columns
	orderby column.Columns
}

func NewSel(cols column.Columns, tbl table.Table, wheres ...*where.Condition) Query {
	return &sel{columns: cols, table: tbl, filter: wheres}
}

func SetGroupbyColumn(q Query, cols ...*column.Column) error {
	sel, ok := q.(*sel)
	if !ok {
		return errors.New("assertion type query to *sel is faild")
	}
	sel.groupby = cols
	return nil
}

func SetOrderByColumn(q Query, cols ...*column.Column) error {
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

func (s *sel) Bind() bind {
	bs := bind{}
	for _, f := range s.filter {
		b := f.Bind()
		bs = append(bs, b...)
	}
	return bs
}
