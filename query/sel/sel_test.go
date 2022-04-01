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
	cond, _ := comp.MakeCondition(false)

	in := where.NewIn("member_id", false, 1, 2, 3)
	condIn, _ := in.MakeCondition(false)

	between := where.NewBetween("price", 500, 1000)
	condBet, _ := between.MakeCondition(false)

	type fields struct {
		columns column.Columns
		table   table.Table
		filter  where.Conditions
		groupby column.Columns
		orderby column.Columns
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
			fields{columns: column.Columns{{CName: "member_id"}}, table: table.Table{TName: "member"}},
			args{false},
			"SELECT member_id FROM member",
			false,
		},
		{
			"Case2: column has alias, table has no alias",
			fields{columns: column.Columns{{CName: "member_id", Alias: "id"}}, table: table.Table{TName: "member"}},
			args{false},
			"SELECT member_id as id FROM member",
			false,
		},
		{
			"Error Case1: No columns",
			fields{columns: column.Columns{}, table: table.Table{TName: "member"}},
			args{false},
			"",
			true,
		},
		{
			"Error Case2: undefined table name ",
			fields{columns: column.Columns{{CName: "member_id", Alias: "id"}}, table: table.Table{}},
			args{false},
			"",
			true,
		},
		{
			"Case3: column and table have alias",
			fields{columns: column.Columns{{CName: "member_id", Alias: "id"}}, table: table.Table{TName: "member", Alias: "m"}},
			args{false},
			"SELECT member_id as id FROM member as m",
			false,
		},
		{
			"Case4: with where statement",
			fields{columns: column.Columns{{CName: "member_id", Alias: "id"}}, table: table.Table{TName: "member", Alias: "m"}, filter: where.Conditions{cond}},
			args{false},
			"SELECT member_id as id FROM member as m WHERE member_id = ?",
			false,
		},
		{
			"Case5: with where statement(IN operator)",
			fields{columns: column.Columns{{CName: "member_id", Alias: "id"}}, table: table.Table{TName: "member", Alias: "m"}, filter: where.Conditions{condIn}},
			args{false},
			"SELECT member_id as id FROM member as m WHERE member_id IN ( ?, ?, ? )",
			false,
		},
		{
			"Case6: with where statement(BETWEEN operator)",
			fields{columns: column.Columns{{CName: "member_id", Alias: "id"}}, table: table.Table{TName: "member", Alias: "m"}, filter: where.Conditions{condBet}},
			args{false},
			"SELECT member_id as id FROM member as m WHERE price BETWEEN ? AND ?",
			false,
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
