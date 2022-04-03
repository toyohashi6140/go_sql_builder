package column

import (
	"fmt"
	"strings"
)

type Column struct {
	CName string
	Alias string
	Value interface{} // for "select" result
}

type Order struct {
	columns Columns
	desc    bool
}

func NewOrder(desc bool, cols ...*Column) *Order {
	return &Order{cols, desc}
}

func (o *Order) Columns() Columns {
	return o.columns
}

func (o *Order) Desc() bool {
	return o.desc
}

func (c *Column) line() string {
	if c.Alias != "" {
		return fmt.Sprintf("%s as %s", c.CName, c.Alias)
	} else {
		return c.CName
	}
}

// Name if you set some alias to column, return alias.
func (c *Column) Name() string {
	if c.Alias != "" {
		return c.Alias
	} else {
		return c.CName
	}
}

type Columns []*Column

func NewColumns(cols ...*Column) Columns {
	return cols
}

// Line if column has alias, return " "name" as "alias" " as text . otherwise, return "name" as text
func (cs Columns) Line() []string {
	columns := []string{}
	for _, c := range cs {
		columns = append(columns, c.line())
	}
	return columns
}

// Name each alias or name of columns join with comma
func (cs Columns) Name() string {
	names := []string{}
	for _, c := range cs {
		names = append(names, c.Name())
	}
	return strings.Join(names, ", ")
}

// Name each name of columns join with comma(for GROUP BY)
func (cs Columns) NoAliasName() string {
	names := []string{}
	for _, c := range cs {
		names = append(names, c.CName)
	}
	return strings.Join(names, ", ")
}
