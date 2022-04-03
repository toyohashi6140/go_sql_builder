package where

import (
	"fmt"
	"strings"
)

const (
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
	MakeCondition(bool) (*Condition, error)
}

type wheres []Where

func NewWheres(ws ...Where) wheres {
	return ws
}

func (ws wheres) Join(or bool) (Conditions, error) {
	cs := Conditions{}
	for _, w := range ws {
		c, err := w.MakeCondition(or)
		if err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}
	return cs, nil
}

type Condition struct {
	condition []string      // SQL with placeholder
	bind      []interface{} // bind value
	or        bool          // if it is true, each word(condition) is joined by "OR" operator
}

func (c Condition) Bind() []interface{} {
	return c.bind
}

type Conditions []*Condition

func (c Condition) join() string {
	if c.or {
		return fmt.Sprintf("(%s)", strings.Join(c.condition, " OR "))
	}
	return strings.Join(c.condition, " AND ")
}

func (cs Conditions) Join(or bool) string {
	// select option "OR" or "AND" of join certain condition and another condition.(for example, A=B and in(C,D))
	var option string
	var format string
	if or {
		format = "(%s)"
		option = " OR "
	} else {
		format = "%s"
		option = " AND "
	}

	joinedCond := []string{}
	for _, c := range cs {
		joinedCond = append(joinedCond, fmt.Sprintf(format, c.join()))
	}

	return strings.Join(joinedCond, option)
}
