package query

import (
	"database/sql"
	"strings"
)

type Bind []interface{}

// Query can build SQL in function Build(), set bind value in function Bind(), and can execute SQL in function Execute()
type Query interface {
	// Build return SQL as string
	// where value etc. are replaced with "?"
	Build() (string, error)

	// Bind returns the original value of what was replaced with "?" As a slice
	// type "Bind" is slice of interface{}
	Bind() Bind

	// Execute processing is different between select and other SQL
	// If select, return the fetched row and column
	// Otherwise, execute SQL and return Result interface
	Execute(*sql.DB, string, ...interface{}) (interface{}, error)
}

// StrQuestion returns a question string according to the length of value
func StrQuestion(len int) string {
	questions := []string{}
	for i := 0; i < len; i++ {
		questions = append(questions, "?")
	}
	return strings.Join(questions, ", ")
}
