package sel

import (
	"testing"

	"github.com/toyohashi6140/go_sql_builder/column"
	"github.com/toyohashi6140/go_sql_builder/table"
	"github.com/toyohashi6140/go_sql_builder/where"
)

func Test_sel_Build(t *testing.T) {
	// make where clause
	cond, _ := where.NewComparisons(
		where.AND,
		where.NewComparison(&column.Column{CName: "member_id"}, 1, where.EQ),
		where.NewComparison(&column.Column{CName: "member_id", Alias: "id"}, 2, where.GT),
	).MakeCondition()
	condOr, _ := where.NewComparisons(
		where.OR,
		where.NewComparison(&column.Column{CName: "member_id"}, 1, where.EQ),
		where.NewComparison(&column.Column{CName: "member_id", Alias: "id"}, 2, where.GT),
	).MakeCondition()

	// in
	condIn, _ := where.NewIns(
		where.AND,
		where.NewIn("member_id", false, 1, 2, 3),
		where.NewIn("company_id", false, 1, 2, 3),
	).MakeCondition()
	condInT, _ := where.NewIns(
		where.OR,
		where.NewIn("member_id", true, 1, 2, 3),
		where.NewIn("company_id", true, 1, 2, 3),
	).MakeCondition()

	// between
	condBet, _ := where.NewBetweens(
		where.AND,
		where.NewBetween("price", 500, 1000),
		where.NewBetween("date", "2022-01-31", "2022-04-30"),
	).MakeCondition()
	condBetT, _ := where.NewBetweens(
		where.OR,
		where.NewBetween("price", 500, 1000),
		where.NewBetween("date", "2022-01-31", "2022-04-30"),
	).MakeCondition()

	//like
	condLike, _ := where.NewLikes(
		where.AND,
		where.NewLike("lastname", "toyohashi", false),
		where.NewLike("firstname", 6140, false),
	).MakeCondition()
	condLikeT, _ := where.NewLikes(
		where.OR,
		where.NewLike("lastname", "toyohashi", true),
		where.NewLike("firstname", 6140, true),
	).MakeCondition()

	type fields struct {
		columns     column.Columns
		table       *table.Table
		filter      *where.Conditions
		groupby     column.Columns
		orderby     *column.Order
		skip, limit int
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"Case1: column and table have no alias",
			fields{
				columns: column.Columns{{CName: "member_id"}},
				table:   &table.Table{TName: "member"},
			},
			"SELECT member_id FROM member",
			false,
		},
		{
			"Case2: column has alias, table has no alias",
			fields{
				columns: column.Columns{{CName: "member_id", Alias: "id"}},
				table:   &table.Table{TName: "member"},
			},
			"SELECT member_id as id FROM member",
			false,
		},
		{
			"Case3: column and table have alias",
			fields{
				columns: column.Columns{{CName: "member_id", Alias: "id"}},
				table:   &table.Table{TName: "member", Alias: "m"},
			},
			"SELECT member_id as id FROM member as m",
			false,
		},
		{
			"Case4: with where clause AND",
			fields{
				columns: column.Columns{{CName: "member_id", Alias: "id"}},
				table:   &table.Table{TName: "member", Alias: "m"},
				filter:  where.NewConditions(where.AND, cond),
			},
			"SELECT member_id as id FROM member as m WHERE member_id = ? AND member_id > ?",
			false,
		},
		{
			"Case5: with where clause OR",
			fields{
				columns: column.Columns{{CName: "member_id", Alias: "id"}},
				table:   &table.Table{TName: "member", Alias: "m"},
				filter:  where.NewConditions(where.OR, condOr),
			},
			"SELECT member_id as id FROM member as m WHERE ((member_id = ? OR member_id > ?))",
			false,
		},
		{
			"Case6: with where clause(IN operator)",
			fields{
				columns: column.Columns{{CName: "member_id", Alias: "id"}},
				table:   &table.Table{TName: "member", Alias: "m"},
				filter:  where.NewConditions(where.AND, condIn),
			},
			"SELECT member_id as id FROM member as m WHERE member_id IN ( ?, ?, ? ) AND company_id IN ( ?, ?, ? )",
			false,
		},
		{
			"Case7: with where clause(IN operator), or flag true",
			fields{
				columns: column.Columns{{CName: "member_id", Alias: "id"}},
				table:   &table.Table{TName: "member", Alias: "m"},
				filter:  where.NewConditions(where.OR, condInT),
			},
			"SELECT member_id as id FROM member as m WHERE ((member_id NOT IN ( ?, ?, ? ) OR company_id NOT IN ( ?, ?, ? )))",
			false,
		},
		{
			"Case8: with where clause(BETWEEN operator)",
			fields{
				columns: column.Columns{{CName: "member_id", Alias: "id"}},
				table:   &table.Table{TName: "member", Alias: "m"},
				filter:  where.NewConditions(where.AND, condBet),
			},
			"SELECT member_id as id FROM member as m WHERE price BETWEEN ? AND ? AND date BETWEEN ? AND ?",
			false,
		},
		{
			"Case9: with where clause(BETWEEN operator), or flag true",
			fields{
				columns: column.Columns{{CName: "member_id", Alias: "id"}},
				table:   &table.Table{TName: "member", Alias: "m"},
				filter:  where.NewConditions(where.OR, condBetT),
			},
			"SELECT member_id as id FROM member as m WHERE ((price BETWEEN ? AND ? OR date BETWEEN ? AND ?))",
			false,
		},
		{
			"Case10: with where clause(LIKE operator)",
			fields{
				columns: column.Columns{{CName: "member_id", Alias: "id"}},
				table:   &table.Table{TName: "member", Alias: "m"},
				filter:  where.NewConditions(where.AND, condLike),
			},
			"SELECT member_id as id FROM member as m WHERE lastname LIKE ? AND firstname LIKE ?",
			false,
		},
		{
			"Case11: with where clause(LIKE operator), or flag true",
			fields{
				columns: column.Columns{{CName: "member_id", Alias: "id"}},
				table:   &table.Table{TName: "member", Alias: "m"},
				filter:  where.NewConditions(where.OR, condLikeT),
			},
			"SELECT member_id as id FROM member as m WHERE ((lastname NOT LIKE ? OR firstname NOT LIKE ?))",
			false,
		},
		{
			"Case12: with multi where clause",
			fields{
				columns: column.Columns{{CName: "member_id", Alias: "id"}},
				table:   &table.Table{TName: "member", Alias: "m"},
				filter:  where.NewConditions(where.AND, cond, condIn, condBet, condLike),
			},
			"SELECT member_id as id FROM member as m WHERE member_id = ? AND member_id > ? AND member_id IN ( ?, ?, ? ) AND company_id IN ( ?, ?, ? ) AND price BETWEEN ? AND ? AND date BETWEEN ? AND ? AND lastname LIKE ? AND firstname LIKE ?",
			false,
		},
		{
			"Case13: with multi where clause",
			fields{
				columns: column.Columns{{CName: "member_id", Alias: "id"}},
				table:   &table.Table{TName: "member", Alias: "m"},
				filter:  where.NewConditions(where.OR, cond, condIn, condBet, condLike),
			},
			"SELECT member_id as id FROM member as m WHERE (member_id = ? AND member_id > ?) OR (member_id IN ( ?, ?, ? ) AND company_id IN ( ?, ?, ? )) OR (price BETWEEN ? AND ? AND date BETWEEN ? AND ?) OR (lastname LIKE ? AND firstname LIKE ?)",
			false,
		},
		{
			"Case14: aggregate with group by clause",
			fields{
				columns: column.Columns{{CName: "count(*)", Alias: "count"}, {CName: "member_id", Alias: "id"}},
				table:   &table.Table{TName: "member", Alias: "m"},
				filter:  where.NewConditions(where.AND, cond),
				groupby: column.Columns{{CName: "member_id", Alias: "m"}},
			},
			"SELECT count(*) as count, member_id as id FROM member as m WHERE member_id = ? AND member_id > ? GROUP BY member_id",
			false,
		},
		{
			"Case15: aggregate with group by clause and sort by ORDER BY clause",
			fields{
				columns: column.Columns{{CName: "count(*)", Alias: "count"}, {CName: "member_id", Alias: "id"}},
				table:   &table.Table{TName: "member", Alias: "m"},
				filter:  where.NewConditions(where.AND, cond),
				groupby: column.Columns{{CName: "member_id", Alias: "m"}},
				orderby: column.NewOrder(false, &column.Column{CName: "member_id", Alias: "id"}),
			},
			"SELECT count(*) as count, member_id as id FROM member as m WHERE member_id = ? AND member_id > ? GROUP BY member_id ORDER BY id",
			false,
		},
		{
			"Case16: aggregate with group by clause and sort by ORDER BY clause",
			fields{
				columns: column.Columns{{CName: "count(*)", Alias: "count"}, {CName: "member_id", Alias: "id"}},
				table:   &table.Table{TName: "member", Alias: "m"},
				filter:  where.NewConditions(where.AND, cond),
				groupby: column.Columns{{CName: "member_id", Alias: "m"}},
				orderby: column.NewOrder(true, &column.Column{CName: "member_id", Alias: "id"}),
			},
			"SELECT count(*) as count, member_id as id FROM member as m WHERE member_id = ? AND member_id > ? GROUP BY member_id ORDER BY id DESC",
			false,
		},
		{
			"Case17: set limit",
			fields{
				columns: column.Columns{{CName: "count(*)", Alias: "count"}, {CName: "member_id", Alias: "id"}},
				table:   &table.Table{TName: "member", Alias: "m"},
				limit:   10,
			},
			"SELECT count(*) as count, member_id as id FROM member as m LIMIT 10",
			false,
		},
		{
			"Case18: set skip and limit",
			fields{
				columns: column.Columns{{CName: "count(*)", Alias: "count"}, {CName: "member_id", Alias: "id"}},
				table:   &table.Table{TName: "member", Alias: "m"},
				skip:    10,
				limit:   10,
			},
			"SELECT count(*) as count, member_id as id FROM member as m LIMIT 10, 10",
			false,
		},
		{
			"Case19: set skip and no set limit(pattarn of limit clause ignored)",
			fields{
				columns: column.Columns{{CName: "count(*)", Alias: "count"}, {CName: "member_id", Alias: "id"}},
				table:   &table.Table{TName: "member", Alias: "m"},
				skip:    10,
			},
			"SELECT count(*) as count, member_id as id FROM member as m",
			false,
		},
		{
			"Error Case1: No columns",
			fields{columns: column.Columns{}, table: &table.Table{TName: "member"}},
			"",
			true,
		},
		{
			"Error Case2: undefined table name",
			fields{columns: column.Columns{{CName: "member_id", Alias: "id"}}, table: &table.Table{}},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sel{
				columns: tt.fields.columns,
				table:   tt.fields.table,
				filter:  tt.fields.filter,
				groupby: tt.fields.groupby,
				orderby: tt.fields.orderby,
				skip:    tt.fields.skip,
				limit:   tt.fields.limit,
			}
			got, err := s.Build()
			if (err != nil) != tt.wantErr {
				t.Errorf("sel.Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("sel.Build() = %v, want %v", got, tt.want)
			}
		})
	}
}
