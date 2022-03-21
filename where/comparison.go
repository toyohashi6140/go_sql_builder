package where

import (
	"errors"
	"fmt"

	"github.com/toyohashi6140/go_sql_builder/column"
)

type comparison struct {
	col      *column.Column
	val      interface{}
	operator string // eq(=), gt(>), lt(<), gte/ge(>=), lte/le(<=), ne
}

// NewComparison construct "where" filter and return as slice. you must also set an valid operator("gt","ge","gte","lt","le","lte","eq" and "ne").
func NewComparison(col *column.Column, val interface{}, o string) *comparison {
	return &comparison{col: col, val: val, operator: o}
}

// Append
func (c *comparison) Append(w Where) Where {
	comps, ok := w.(*comparisons)
	if ok {
		*comps = append(*comps, c)
	}
	return comps
}

type comparisons []*comparison

// NewComparisons
func NewComparisons() Where {
	return &comparisons{}
}

// NewComparisons make comparisons slice as Where interface. each comparison struct sets uniform operator(EQ)
func NewComparisonsEq(cols column.Columns, vals ...interface{}) (Where, error) {
	if len(cols) != len(vals) {
		return nil, errors.New("the number of column and value is different these must be the same number")
	}
	comps := comparisons{}
	for i, col := range cols {
		comps = append(comps, &comparison{col: col, val: vals[i], operator: EQ})
	}
	return &comps, nil
}

// MakeCondition this function is make where expressoion if you set column name(or alias), table name(or alias), and operator.
func (c *comparisons) MakeCondition(or bool) (*Condition, error) {
	if len(*c) == 0 {
		return nil, errors.New("no keys or values are setting. can't make an expression")
	}
	condStrs := []string{}
	bind := []interface{}{}
	for _, c := range *c {
		if c.col != nil && c.val != nil {
			symbol := operandToSymbol(c.operator)
			if symbol == "" {
				return nil, errors.New("operator is undefined. can't evaluate your where query")
			}
			condStrs = append(condStrs, fmt.Sprintf("%s %s %s", c.col.Name(), symbol, "?"))
			bind = append(bind, c.val)
		} else {
			return nil, errors.New("no keys or values are setting. can't make an expression")
		}
	}
	cond := &Condition{condition: condStrs, bind: bind, or: or}
	return cond, nil
}

func operandToSymbol(operand string) string {
	switch operand {
	case "gt":
		return ">"
	case "ge", "gte":
		return ">="
	case "lt":
		return "<"
	case "le", "lte":
		return "<="
	case "eq":
		return "="
	case "ne":
		return "!="
	default:
		return ""
	}
}
