package sel

import (
	"testing"

	"github.com/toyohashi6140/go_sql_builder/column"
	"github.com/toyohashi6140/go_sql_builder/table"
	"github.com/toyohashi6140/go_sql_builder/where"
)

func Test_sel_Build(t *testing.T) {
	// make where statement
	comp := where.NewComparisons()
	comp = where.NewComparison(&column.Column{CName: "member_id"}, 1, where.EQ).Append(comp)
	comp = where.NewComparison(&column.Column{CName: "member_id", Alias: "id"}, 2, where.GT).Append(comp)
	cond, _ := comp.MakeCondition(false)
	condOr, _ := comp.MakeCondition(true)

	condIn, _ := where.NewIn("member_id", false, 1, 2, 3).MakeCondition(false)
	condInT, _ := where.NewIn("member_id", true, 1, 2, 3).MakeCondition(true)

	condBet, _ := where.NewBetween("price", 500, 1000).MakeCondition(false)
	condBetT, _ := where.NewBetween("price", 500, 1000).MakeCondition(true)

	condLike, _ := where.NewLike("name", "toyohashi", false).MakeCondition(false)
	condLikeT, _ := where.NewLike("name", "toyohashi", true).MakeCondition(true)

	type fields struct {
		columns column.Columns
		table   *table.Table
		filter  where.Conditions
		groupby column.Columns
		orderby *column.Order
	}
	type args struct {
		or bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"Case1: column and table have no alias",
			fields{columns: column.Columns{{CName: "member_id"}}, table: &table.Table{TName: "member"}},
			args{false},
			"SELECT member_id FROM member",
			false,
		},
		{
			"Case2: column has alias, table has no alias",
			fields{columns: column.Columns{{CName: "member_id", Alias: "id"}}, table: &table.Table{TName: "member"}},
			args{false},
			"SELECT member_id as id FROM member",
			false,
		},
		{
			"Case3: column and table have alias",
			fields{columns: column.Columns{{CName: "member_id", Alias: "id"}}, table: &table.Table{TName: "member", Alias: "m"}},
			args{false},
			"SELECT member_id as id FROM member as m",
			false,
		},
		{
			"Case4: with where statement AND",
			fields{columns: column.Columns{{CName: "member_id", Alias: "id"}}, table: &table.Table{TName: "member", Alias: "m"}, filter: where.Conditions{cond}},
			args{false},
			"SELECT member_id as id FROM member as m WHERE member_id = ? AND member_id > ?",
			false,
		},
		{
			"Case5: with where statement OR",
			fields{columns: column.Columns{{CName: "member_id", Alias: "id"}}, table: &table.Table{TName: "member", Alias: "m"}, filter: where.Conditions{condOr}},
			args{false},
			"SELECT member_id as id FROM member as m WHERE (member_id = ? OR member_id > ?)",
			false,
		},
		{
			"Case6: with where statement(IN operator)",
			fields{columns: column.Columns{{CName: "member_id", Alias: "id"}}, table: &table.Table{TName: "member", Alias: "m"}, filter: where.Conditions{condIn}},
			args{false},
			"SELECT member_id as id FROM member as m WHERE member_id IN ( ?, ?, ? )",
			false,
		},
		{
			"Case7: with where statement(IN operator), or flag true",
			fields{columns: column.Columns{{CName: "member_id", Alias: "id"}}, table: &table.Table{TName: "member", Alias: "m"}, filter: where.Conditions{condInT}},
			args{false},
			"SELECT member_id as id FROM member as m WHERE (member_id NOT IN ( ?, ?, ? ))",
			false,
		},
		{
			"Case8: with where statement(BETWEEN operator)",
			fields{columns: column.Columns{{CName: "member_id", Alias: "id"}}, table: &table.Table{TName: "member", Alias: "m"}, filter: where.Conditions{condBet}},
			args{false},
			"SELECT member_id as id FROM member as m WHERE price BETWEEN ? AND ?",
			false,
		},
		{
			"Case9: with where statement(BETWEEN operator), or flag true",
			fields{columns: column.Columns{{CName: "member_id", Alias: "id"}}, table: &table.Table{TName: "member", Alias: "m"}, filter: where.Conditions{condBetT}},
			args{false},
			"SELECT member_id as id FROM member as m WHERE (price BETWEEN ? AND ?)",
			false,
		},
		{
			"Case10: with where statement(LIKE operator)",
			fields{columns: column.Columns{{CName: "member_id", Alias: "id"}}, table: &table.Table{TName: "member", Alias: "m"}, filter: where.Conditions{condLike}},
			args{false},
			"SELECT member_id as id FROM member as m WHERE name LIKE ?",
			false,
		},
		{
			"Case11: with where statement(LIKE operator), or flag true",
			fields{columns: column.Columns{{CName: "member_id", Alias: "id"}}, table: &table.Table{TName: "member", Alias: "m"}, filter: where.Conditions{condLikeT}},
			args{false},
			"SELECT member_id as id FROM member as m WHERE (name NOT LIKE ?)",
			false,
		},
		{
			"Case12: with multi where statement",
			fields{columns: column.Columns{{CName: "member_id", Alias: "id"}}, table: &table.Table{TName: "member", Alias: "m"}, filter: where.Conditions{cond, condIn, condBet, condLike}},
			args{false},
			"SELECT member_id as id FROM member as m WHERE member_id = ? AND member_id > ? AND member_id IN ( ?, ?, ? ) AND price BETWEEN ? AND ? AND name LIKE ?",
			false,
		},
		{
			"Case13: with multi where statement",
			fields{columns: column.Columns{{CName: "member_id", Alias: "id"}}, table: &table.Table{TName: "member", Alias: "m"}, filter: where.Conditions{cond, condIn, condBet, condLike}},
			args{true},
			"SELECT member_id as id FROM member as m WHERE (member_id = ? AND member_id > ?) OR (member_id IN ( ?, ?, ? )) OR (price BETWEEN ? AND ?) OR (name LIKE ?)",
			false,
		},
		{
			"Case14: aggregate with group by statement",
			fields{
				columns: column.Columns{{CName: "count(*)", Alias: "count"}, {CName: "member_id", Alias: "id"}},
				table:   &table.Table{TName: "member", Alias: "m"},
				filter:  where.Conditions{cond},
				groupby: column.Columns{{CName: "member_id", Alias: "m"}},
			},
			args{false},
			"SELECT count(*) as count, member_id as id FROM member as m WHERE member_id = ? AND member_id > ? GROUP BY member_id",
			false,
		},
		{
			"Case15: aggregate with group by statement and sort by ORDER BY statement",
			fields{
				columns: column.Columns{{CName: "count(*)", Alias: "count"}, {CName: "member_id", Alias: "id"}},
				table:   &table.Table{TName: "member", Alias: "m"},
				filter:  where.Conditions{cond},
				groupby: column.Columns{{CName: "member_id", Alias: "m"}},
				orderby: column.NewOrder(false, &column.Column{CName: "member_id", Alias: "id"}),
			},
			args{false},
			"SELECT count(*) as count, member_id as id FROM member as m WHERE member_id = ? AND member_id > ? GROUP BY member_id ORDER BY id",
			false,
		},
		{
			"Case16: aggregate with group by statement and sort by ORDER BY statement",
			fields{
				columns: column.Columns{{CName: "count(*)", Alias: "count"}, {CName: "member_id", Alias: "id"}},
				table:   &table.Table{TName: "member", Alias: "m"},
				filter:  where.Conditions{cond},
				groupby: column.Columns{{CName: "member_id", Alias: "m"}},
				orderby: column.NewOrder(true, &column.Column{CName: "member_id", Alias: "id"}),
			},
			args{false},
			"SELECT count(*) as count, member_id as id FROM member as m WHERE member_id = ? AND member_id > ? GROUP BY member_id ORDER BY id DESC",
			false,
		},
		{
			"Error Case1: No columns",
			fields{columns: column.Columns{}, table: &table.Table{TName: "member"}},
			args{false},
			"",
			true,
		},
		{
			"Error Case2: undefined table name ",
			fields{columns: column.Columns{{CName: "member_id", Alias: "id"}}, table: &table.Table{}},
			args{false},
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
			}
			got, err := s.Build(tt.args.or)
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