package where

import (
	"fmt"
)

type like struct {
	key     string
	keyword interface{}
	not     bool
}

func NewLike(k string, kw interface{}, not bool) *like {
	return &like{key: k, keyword: kw, not: not}
}

type likes struct {
	likes           []*like
	logicalOperator int
}

func NewLikes(lo int, ls ...*like) Where {
	return &likes{ls, lo}
}

func (ls *likes) MakeCondition() (*Condition, error) {
	conds := []string{}
	binds := []interface{}{}
	for _, l := range ls.likes {
		var operator string
		if l.not {
			operator = "NOT LIKE"
		} else {
			operator = "LIKE"
		}
		conds = append(conds, fmt.Sprintf("%s %s ?", l.key, operator))
		binds = append(binds, fmt.Sprintf("%%%s%%", l.keyword))
	}
	return &Condition{
		condition:       conds,
		bind:            binds,
		logicalOperator: ls.logicalOperator,
	}, nil
}
