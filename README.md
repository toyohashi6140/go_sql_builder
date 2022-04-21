# go_sql_builder

This package enables automatic generation and execution of SQL.
All you have to do is create a Column type and name the column.

Below is the sample code.
The variable db is your database connection information, here we assume it was received in advance.

``` 
package sample

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/toyohashi6140/go_sql_builder/column"
	"github.com/toyohashi6140/go_sql_builder/query/sel"
	"github.com/toyohashi6140/go_sql_builder/table"
	"github.com/toyohashi6140/go_sql_builder/where"
)

func Sample(db *sql.DB) {
	// select columns
	selCols := column.NewColumns(
		&column.Column{CName: "member_id", Alias: "id"},
		&column.Column{CName: "status"},
	)

	// where columns
	whereCols := column.NewColumns(
		selCols[0],
		&column.Column{CName: "email"},
	)

	// you can set a separate operator for each column by doing the following
	// If you want to use the selected columns in the where clause,
	// (perhaps when you need to write a Where clause with an alias)
	// include the element of selCols in whereCols or use the element of selCols directly in the argument of NewComparison.
	cond, err := where.NewComparisons(
		where.AND,
		where.NewComparison(whereCols[0], 46156, where.GE),
		where.NewComparison(selCols[1], 1, where.GE),
		where.NewComparison(whereCols[1], "sakamoto@bestperson.jp", where.EQ),
	).MakeCondition()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Use the above if you want to use an operator for each expression.
	// If all equations are equalities, we provide a way to set them all at once.
	compEq, err := where.NewComparisonsEq(whereCols, 46156, "sakamoto@bestperson.jp")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	cond2, err := compEq.MakeCondition()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// If you want to express a query that uses the "in" operator, this function is effective.
	// If true is specified for the second argument, it will be a negative operation.
	condI, err := where.NewIns(
		where.AND,
		where.NewIn("member_id", false, 1, 2, 3, 4, 5, 6),
	).MakeCondition()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// If you want to use "between" operator, this function can be used.
	// The arguments are the column name, A, B in order from the left.
	// Now you can create statement " `column` Between A and B "
	condB, err := where.NewBetweens(
		where.AND,
		where.NewBetween("entry_date", time.Date(2022, time.December, 20, 0, 0, 0, 0, time.Local), time.Now()),
	).MakeCondition()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Use this for "like" syntax
	// If true is specified for the third argument, it will be a negative operation.
	condL, err := where.NewLikes(
		where.AND,
		where.NewLike("my_name_sei", "武藤", false),
	).MakeCondition()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	q := sel.New(selCols, &table.Table{TName: "member"})

	err = sel.SetFilterColumn(q, where.AND, cond, cond2, condI, condB, condL)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = sel.SetGroupbyColumn(q, selCols[0])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = sel.SetOrderByColumn(q, false, selCols[0])
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	query, err := q.Build()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(query)
	fmt.Println(q.Bind())

	// db is your Database.
	i, err := q.Execute(db, "", 1)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	maps, ok := i.([]map[string]interface{})
	if !ok {
		fmt.Println("No Select Query")
		return
	}

	for _, m := range maps {
		for k, v := range m {
			fmt.Printf("key: %v, val: %v\n", k, v)
		}
	}
}

```
