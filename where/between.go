package where

import "fmt"

type between struct {
	key        string
	valA, valB interface{}
}

func NewBetween(k string, a, b interface{}) *between {
	return &between{key: k, valA: a, valB: b}
}

type betweens struct {
	betweens        []*between
	logicalOperator int
}

func NewBetweens(lo int, bs ...*between) Where {
	return &betweens{bs, lo}
}

func (bs *betweens) MakeCondition() (*Condition, error) {
	c := &Condition{logicalOperator: bs.logicalOperator}
	for _, b := range bs.betweens {
		c.condition = append(c.condition, fmt.Sprintf("%s BETWEEN %s AND %s", b.key, "?", "?"))
		c.bind = append(c.bind, b.valA, b.valB)
	}
	return c, nil
}
