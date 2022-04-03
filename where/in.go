package where

import (
	"fmt"

	"github.com/toyohashi6140/go_sql_builder/query"
)

type in struct {
	key string
	val []interface{}
	not bool
}

func NewIn(k string, not bool, v ...interface{}) Where {
	return in{key: k, val: v, not: not}
}

func (i in) MakeCondition(or bool) (*Condition, error) {
	questions := query.StrQuestion(len(i.val))
	var operator string
	if i.not {
		operator = "NOT IN"
	} else {
		operator = "IN"
	}
	return &Condition{
		condition: []string{fmt.Sprintf("%s %s ( %s )", i.key, operator, questions)},
		bind:      i.val,
		or:        or,
	}, nil
}
