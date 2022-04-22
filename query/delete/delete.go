package delete

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/toyohashi6140/go_sql_builder/query"
	"github.com/toyohashi6140/go_sql_builder/table"
	"github.com/toyohashi6140/go_sql_builder/where"
)

type delete struct {
	table  *table.Table
	filter *where.Conditions
}

func New(t *table.Table, f *where.Conditions) query.Query {
	return &delete{t, f}
}

func (d *delete) Build() (string, error) {
	if d.table.TName == "" || d.table == nil {
		return "", errors.New("no tables selected")
	}
	if d.filter == nil || len(d.filter.Conditions()) == 0 {
		return "", errors.New("you can't build \"DELETE\" query without where clause")
	}
	return fmt.Sprintf("DELETE FROM %s WHERE %s", d.table.TName, d.filter.Join()), nil
}

func (d *delete) Bind() query.Bind {
	return d.filter.Bind()
}

func (d *delete) Execute(db *sql.DB, sql string, bind ...interface{}) (interface{}, error) {
	result, err := db.Exec(sql, bind...)
	if err != nil {
		return nil, err
	}
	return result, nil
}
