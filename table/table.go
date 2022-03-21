package table

import "fmt"

type Table struct {
	TName string
	Alias string
}

func (t Table) Line() string {
	if t.Alias != "" {
		return fmt.Sprintf("%s as %s", t.TName, t.Alias)
	} else {
		return t.TName
	}
}
