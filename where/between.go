package where

import "fmt"

type between struct {
	key        string
	valA, valB interface{}
}

func (b between) MakeCondition(or bool) (*Condition, error) {
	c := &Condition{
		condition: []string{fmt.Sprintf("%s BETWEEN %s AND %s", b.key, "?", "?")},
		bind:      []interface{}{b.valA, b.valB},
		or:        or,
	}
	return c, nil
}

func NewBetween(k string, a, b interface{}) Where {
	return between{key: k, valA: a, valB: b}
}
