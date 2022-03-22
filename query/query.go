package query

import (
	"database/sql"
)

type Bind []interface{}

type Query interface {
	Build(bool) (string, error)
	Bind() Bind
	Execute(db *sql.DB) (interface{}, error)
}
