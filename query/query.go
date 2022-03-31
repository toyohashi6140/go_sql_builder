package query

import "database/sql"

type Bind []interface{}

type Query interface {
	// Build return SQL as string
	// where value etc. are replaced with "?"
	Build(bool) (string, error)

	// Bind returns the original value of what was replaced with "?" As a slice
	// type "Bind" is slice of interface{}
	Bind() Bind

	// Execute processing is different between select and other SQL
	// If select, return the fetched row and column
	// Otherwise, execute SQL and return Result interface
	Execute(*sql.DB) (interface{}, error)
}
