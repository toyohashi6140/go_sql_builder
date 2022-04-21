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

func NewIn(k string, not bool, v ...interface{}) *in {
	return &in{key: k, val: v, not: not}
}

type ins struct {
	ins             []*in
	logicalOperator int
}

func NewIns(lo int, is ...*in) Where {
	return &ins{is, lo}
}

func (is *ins) MakeCondition() (*Condition, error) {
	conds := []string{}
	binds := []interface{}{}
	for _, i := range is.ins {
		questions := query.StrQuestion(len(i.val))
		var operator string
		if i.not {
			operator = "NOT IN"
		} else {
			operator = "IN"
		}
		conds = append(conds, fmt.Sprintf("%s %s ( %s )", i.key, operator, questions))
		binds = append(binds, i.val...)
	}
	return &Condition{
		condition:       conds,
		bind:            binds,
		logicalOperator: is.logicalOperator,
	}, nil
}
