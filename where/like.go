package where

import (
	"fmt"
)

type like struct {
	key     string
	keyword interface{}
	not     bool
}

func NewLike(k string, kw interface{}, not bool) Where {
	return like{key: k, keyword: kw, not: not}
}

func (l like) MakeCondition(or bool) (*Condition, error) {
	var operator string
	if l.not {
		operator = "NOT LIKE"
	} else {
		operator = "LIKE"
	}
	ls := &Condition{
		condition: []string{fmt.Sprintf("%s %s ?", l.key, operator)},
		bind:      []interface{}{fmt.Sprintf("%%%s%%", l.keyword)},
		or:        or,
	}
	return ls, nil
}
