package delete

import (
	"testing"

	"github.com/toyohashi6140/go_sql_builder/column"
	"github.com/toyohashi6140/go_sql_builder/table"
	"github.com/toyohashi6140/go_sql_builder/where"
)

func Test_delete_Build(t *testing.T) {
	cond, _ := where.NewComparisons(
		where.AND,
		where.NewComparison(
			&column.Column{CName: "member_id"},
			1,
			where.EQ,
		),
		where.NewComparison(
			&column.Column{CName: "status"},
			0,
			where.EQ,
		),
	).MakeCondition()
	type fields struct {
		table  *table.Table
		filter *where.Conditions
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"Case 1: OK",
			fields{
				&table.Table{TName: "member"},
				where.NewConditions(where.AND, cond),
			},
			"DELETE FROM member WHERE member_id = 1 AND status = 0",
			false,
		},
		{
			"ErrorCase 1: no tables",
			fields{
				nil,
				where.NewConditions(where.AND, cond),
			},
			"",
			true,
		},
		{
			"ErrorCase 2: table name is unset",
			fields{
				&table.Table{},
				where.NewConditions(where.AND, cond),
			},
			"",
			true,
		},
		{
			"ErrorCase 3: no conditions",
			fields{
				&table.Table{TName: "member"},
				nil,
			},
			"",
			true,
		},
		{
			"ErrorCase 4: conditions is enpty",
			fields{
				&table.Table{TName: "member"},
				&where.Conditions{},
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &delete{
				table:  tt.fields.table,
				filter: tt.fields.filter,
			}
			got, err := d.Build()
			if (err != nil) != tt.wantErr {
				t.Errorf("delete.Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("delete.Build() = %v, want %v", got, tt.want)
			}
		})
	}
}
