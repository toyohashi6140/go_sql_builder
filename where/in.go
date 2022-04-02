package where

import (
	"fmt"
	"strings"
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
	questions := []string{}
	for idx := 0; idx < len(i.val); idx++ {
		questions = append(questions, "?")
	}
	var operator string
	if i.not {
		operator = "NOT IN"
	} else {
		operator = "IN"
	}
	return &Condition{
		condition: []string{fmt.Sprintf("%s %s ( %s )", i.key, operator, strings.Join(questions, ", "))},
		bind:      i.val,
		or:        or,
	}, nil
}
