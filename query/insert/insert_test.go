package insert

import (
	"testing"

	"github.com/toyohashi6140/go_sql_builder/column"
	"github.com/toyohashi6140/go_sql_builder/table"
)

func Test_insert_Build(t *testing.T) {
	type args struct {
		or bool
	}
	tests := []struct {
		name    string
		i       *insert
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"case1",
			&insert{table: &table.Table{TName: "member", Alias: "m"}, values: []interface{}{1, "firstname", "lastname", "2000-01-01"}},
			args{false},
			"INSERT INTO member VALUES( ?, ?, ?, ? )",
			false,
		},
		{
			"case2",
			&insert{
				&table.Table{TName: "member", Alias: "m"},
				column.NewColumns(
					&column.Column{CName: "member_id", Alias: "id"},
					&column.Column{CName: "first_name"},
					&column.Column{CName: "last_name"},
					&column.Column{CName: "birthday", Alias: "bd"},
				),
				[]interface{}{1, "firstname", "lastname", "2000-01-01"},
			},
			args{false},
			"INSERT INTO member( member_id, first_name, last_name, birthday ) VALUES( ?, ?, ?, ? )",
			false,
		},
		{
			"error case",
			&insert{
				&table.Table{TName: "member", Alias: "m"},
				column.NewColumns(
					&column.Column{CName: "member_id", Alias: "id"},
					&column.Column{CName: "first_name"},
					&column.Column{CName: "last_name"},
				),
				[]interface{}{1, "firstname", "lastname", "2000-01-01"},
			},
			args{false},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.i.Build(tt.args.or)
			if (err != nil) != tt.wantErr {
				t.Errorf("insert.Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("insert.Build() = %v, want %v", got, tt.want)
			}
		})
	}
}
