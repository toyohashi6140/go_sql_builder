package update

import (
	"testing"

	"github.com/toyohashi6140/go_sql_builder/column"
	"github.com/toyohashi6140/go_sql_builder/table"
	"github.com/toyohashi6140/go_sql_builder/where"
)

func Test_update_Build(t *testing.T) {
	cond, _ := where.NewComparisons(
		where.AND,
		where.NewComparison(
			&column.Column{CName: "member_id"},
			1,
			where.EQ,
		),
	).MakeCondition()
	type fields struct {
		table     *table.Table
		keyvalues []*keyvalue
		filter    *where.Conditions
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"Case1: regular case and simple value set",
			fields{
				&table.Table{TName: "member"},
				[]*keyvalue{{"status", 1}},
				where.NewConditions(where.AND, cond),
			},
			"UPDATE member SET status = ? WHERE member_id = ?",
			false,
		},
		{
			"Case2: regular case and multi value set",
			fields{
				&table.Table{TName: "member"},
				[]*keyvalue{{"status", 1}, {"email", "aaa@gmail.com"}},
				where.NewConditions(where.AND, cond),
			},
			"UPDATE member SET status = ?, email = ? WHERE member_id = ?",
			false,
		},
		{
			"Error Case 1: table is nil",
			fields{
				nil,
				[]*keyvalue{{"status", 1}},
				where.NewConditions(where.AND, cond),
			},
			"",
			true,
		},
		{
			"Error Case 2: no table name",
			fields{
				&table.Table{},
				[]*keyvalue{{"status", 1}},
				where.NewConditions(where.AND, cond),
			},
			"",
			true,
		},
		{
			"Error Case 3: keyvalue is nil",
			fields{
				&table.Table{TName: "a"},
				nil,
				where.NewConditions(where.AND, cond),
			},
			"",
			true,
		},
		{
			"Error Case 4: keyvalue length is 0",
			fields{
				&table.Table{TName: "a"},
				[]*keyvalue{},
				where.NewConditions(where.AND, cond),
			},
			"",
			true,
		},
		{
			"Error Case 5: filter is nil",
			fields{
				&table.Table{TName: "a"},
				[]*keyvalue{{"key", "value"}},
				nil,
			},
			"",
			true,
		},
		{
			"Error Case 6: filter length is 0",
			fields{
				&table.Table{TName: "a"},
				[]*keyvalue{{"key", "value"}},
				&where.Conditions{},
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &update{
				table:     tt.fields.table,
				keyvalues: tt.fields.keyvalues,
				filter:    tt.fields.filter,
			}
			got, err := u.Build()
			if (err != nil) != tt.wantErr {
				t.Errorf("update.Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("update.Build() = %v, want %v", got, tt.want)
			}
		})
	}
}
