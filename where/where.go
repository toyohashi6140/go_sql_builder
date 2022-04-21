package where

import (
	"fmt"
	"strings"
)

const (
	// AND: 0 OR: 1
	AND = iota
	OR
	GT      string = "gt"
	GE      string = "ge"
	GTE     string = "gte"
	LT      string = "lt"
	LE      string = "le"
	LTE     string = "lte"
	EQ      string = "eq"
	NE      string = "ne"
	NULL    string = "null"
	NOTNULL string = "not null"
)

type Where interface {
	MakeCondition() (*Condition, error)
}

type wheres []Where

func NewWheres(ws ...Where) wheres {
	return ws
}

func (ws wheres) Join(or bool) (*Conditions, error) {
	cs := &Conditions{}
	for _, w := range ws {
		c, err := w.MakeCondition()
		if err != nil {
			return nil, err
		}
		cs.conditions = append(cs.conditions, c)
	}
	return cs, nil
}

type Condition struct {
	condition       []string      // SQL with placeholder
	bind            []interface{} // bind value
	logicalOperator int           // if it is true, each word(condition) is joined by "OR" operator
}

func (c Condition) Bind() []interface{} {
	return c.bind
}

func (c Condition) join() string {
	if c.logicalOperator == 1 {
		return fmt.Sprintf("(%s)", strings.Join(c.condition, " OR "))
	}
	return strings.Join(c.condition, " AND ")
}

type Conditions struct {
	conditions      []*Condition
	logicalOperator int
}

func NewConditions(lo int, conds ...*Condition) *Conditions {
	return &Conditions{conds, lo}
}

func (cs *Conditions) Conditions() []*Condition {
	return cs.conditions
}

func (cs *Conditions) Join() string {
	// select option "OR" or "AND" of join certain condition and another condition.(for example, A=B and in(C,D))
	var option string
	var format string
	if cs.logicalOperator == 1 {
		format = "(%s)"
		option = " OR "
	} else {
		format = "%s"
		option = " AND "
	}

	joinedCond := []string{}
	for _, c := range cs.conditions {
		joinedCond = append(joinedCond, fmt.Sprintf(format, c.join()))
	}

	return strings.Join(joinedCond, option)
}
