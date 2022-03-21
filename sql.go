package sql

import (
	"fmt"
	"go_batch/internal/database/query/sql/query"
	"go_batch/pkg/controller/database"
)

type sql struct {
	query query.Query
}

func New(q query.Query) sql {
	return sql{query: q}
}

func (s sql) Execute(or bool) error {
	pre, err := s.query.Build(or)
	if err != nil {
		return err
	}
	bind := s.query.Bind()
	transaction := database.GetTransaction(pre, bind...)
	fmt.Println(transaction)
	return nil
}
