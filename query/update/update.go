package update

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/toyohashi6140/go_sql_builder/query"
	"github.com/toyohashi6140/go_sql_builder/table"
	"github.com/toyohashi6140/go_sql_builder/where"
)

type (
	update struct {
		table     *table.Table
		keyvalues []*keyvalue
		filter    *where.Conditions
	}
	keyvalue struct {
		key   string
		value interface{}
	}
)

func NewKeyValue(key string, value interface{}) *keyvalue {
	return &keyvalue{key: key, value: value}
}

func NewKeyValueForMap(m map[string]interface{}) []*keyvalue {
	var kvs []*keyvalue
	for k, v := range m {
		kv := &keyvalue{key: k, value: v}
		kvs = append(kvs, kv)
	}
	return kvs
}

func New(t *table.Table, f *where.Conditions, kvs ...*keyvalue) query.Query {
	return &update{t, kvs, f}
}

func (u *update) Build() (string, error) {
	if u.table == nil || u.table.TName == "" {
		return "", errors.New("no tables selected")
	}
	if u.keyvalues == nil || len(u.keyvalues) == 0 {
		return "", errors.New("no key value to set")
	}
	if u.filter == nil || len(u.filter.Conditions()) == 0 {
		return "", errors.New("you can't build \"UPDATE\" query without where clause")
	}
	format := "UPDATE %s SET %s WHERE %s"
	var kvStr []string
	for _, kv := range u.keyvalues {
		kvStr = append(kvStr, fmt.Sprintf("%s = ?", kv.key))
	}
	return fmt.Sprintf(format, u.table.TName, strings.Join(kvStr, ", "), u.filter.Join()), nil
}

func (u *update) Bind() query.Bind {
	return u.filter.Bind()
}

func (u *update) Execute(db *sql.DB, sql string, bind ...interface{}) (interface{}, error) {
	result, err := db.Exec(sql, bind...)
	if err != nil {
		return nil, err
	}
	return result, nil
}
