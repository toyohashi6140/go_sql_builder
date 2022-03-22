package builder

import (
	"database/sql"
	"fmt"

	"github.com/toyohashi6140/go_sql_builder/query"
)

type builder struct {
	query query.Query
}

func New(q query.Query) builder {
	return builder{query: q}
}

func (b builder) Execute(or bool) error {
	i, err := b.query.Execute(&sql.DB{})
	if err != nil {
		return err
	}
	fmt.Println(i)
	return nil
}
