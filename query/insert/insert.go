package insert

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/toyohashi6140/go_sql_builder/column"
	"github.com/toyohashi6140/go_sql_builder/query"
	"github.com/toyohashi6140/go_sql_builder/table"
)

type (
	insert struct {
		table   *table.Table
		columns column.Columns
		values  query.Bind
	}

	InsertSetter interface {
		Columns(cols ...*column.Column) query.Query
		ToQuery() query.Query
	}
)

func New(t *table.Table, vals ...interface{}) InsertSetter {
	return &insert{table: t, values: vals}
}

func (i *insert) Columns(cols ...*column.Column) query.Query {
	i.columns = cols
	return nil
}

func (i *insert) ToQuery() query.Query {
	return i
}

func (i *insert) Build() (string, error) {
	if i.table == nil || i.table.TName == "" {
		return "", errors.New("no tables selected")
	}
	var sql string
	questions := query.StrQuestion(len(i.values))
	if len(i.columns) > 0 {
		if len(i.columns) != len(i.values) {
			return "", errors.New("the number of column, and value is different. if you execute this query, might be returned error by DBMS")
		}
		sql = fmt.Sprintf("INSERT INTO %s( %s ) VALUES( %s )", i.table.TName, i.columns.NoAliasName(), questions)
	} else {
		sql = fmt.Sprintf("INSERT INTO %s VALUES( %s )", i.table.TName, questions)
	}
	return sql, nil
}

func (i *insert) Bind() query.Bind {
	return i.values
}

func (i *insert) Execute(db *sql.DB, sql string, values ...interface{}) (interface{}, error) {
	result, err := db.Exec(sql, i.values...)
	if err != nil {
		return nil, err
	}
	return result, nil
}
