package sql

import "github.com/toyohashi6140/go_sql_builder/query"

type sql struct {
	query query.Query
}

func New(q query.Query) sql {
	return sql{query: q}
}

// func (s sql) Execute(or bool) error {
// 	pre, err := s.query.Build(or)
// 	if err != nil {
// 		return err
// 	}
// 	bind := s.query.Bind()
// 	transaction := database.GetTransaction(pre, bind...)
// 	fmt.Println(transaction)
// 	return nil
// }
