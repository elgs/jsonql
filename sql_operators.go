// mysql_operators
package jsonql

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var SqlOperators = map[string]*Operator{
	"OR": &Operator{
		Precedence: 1,
		Eval:       evalSql,
	},
	"AND": &Operator{
		Precedence: 3,
		Eval:       evalSql,
	},
	"=": &Operator{
		Precedence: 5,
		Eval:       evalSql,
	},
	"!=": &Operator{
		Precedence: 5,
		Eval:       evalSql,
	},
	">": &Operator{
		Precedence: 5,
		Eval:       evalSql,
	},
	"<": &Operator{
		Precedence: 5,
		Eval:       evalSql,
	},
	">=": &Operator{
		Precedence: 5,
		Eval:       evalSql,
	},
	"<=": &Operator{
		Precedence: 5,
		Eval:       evalSql,
	},
	//	"LIKE": &Operator{
	//		Precedence: 5,
	//		Eval:       evalSql,
	//	},
	//	"NOT_LIKE": &Operator{
	//		Precedence: 5,
	//		Eval:       evalSql,
	//	},
	//	"IS_NULL": &Operator{
	//		Precedence: 5,
	//		Eval:       evalSql,
	//	},
	//	"IS_NOT_NULL": &Operator{
	//		Precedence: 5,
	//		Eval:       evalSql,
	//	},
	"RLIKE": &Operator{
		Precedence: 5,
		Eval:       evalSql,
	},
}

func evalSql(op string, left string, right string) (string, error) {
	switch op {
	case "AND":
		l, err := strconv.ParseBool(left)
		if err != nil {
			return "false", err
		}
		r, err := strconv.ParseBool(right)
		if err != nil {
			return "false", err
		}
		return strconv.FormatBool(l && r), nil
	case "=":
		l, err := strconv.ParseBool(left)
		if err != nil {
			return "false", err
		}
		r, err := strconv.ParseBool(right)
		if err != nil {
			return "false", err
		}
		return strconv.FormatBool(l || r), nil
	case "!=":
		return strconv.FormatBool(left != right), nil
	case ">":
		if strings.HasPrefix(right, "'") && strings.HasSuffix(right, "'") {
			r := right[1 : len(right)-1]
			return strconv.FormatBool(left > r), nil
		} else {
			l, err := strconv.ParseFloat(left, 64)
			if err != nil {
				return "false", err
			}
			r, err := strconv.ParseFloat(right, 64)
			if err != nil {
				return "false", err
			}
			return strconv.FormatBool(l > r), nil
		}
	case "<":
		if strings.HasPrefix(right, "'") && strings.HasSuffix(right, "'") {
			r := right[1 : len(right)-1]
			return strconv.FormatBool(left < r), nil
		} else {
			l, err := strconv.ParseFloat(left, 64)
			if err != nil {
				return "false", err
			}
			r, err := strconv.ParseFloat(right, 64)
			if err != nil {
				return "false", err
			}
			return strconv.FormatBool(l < r), nil
		}
	case ">=":
		if strings.HasPrefix(right, "'") && strings.HasSuffix(right, "'") {
			r := right[1 : len(right)-1]
			return strconv.FormatBool(left >= r), nil
		} else {
			l, err := strconv.ParseFloat(left, 64)
			if err != nil {
				return "false", err
			}
			r, err := strconv.ParseFloat(right, 64)
			if err != nil {
				return "false", err
			}
			return strconv.FormatBool(l >= r), nil
		}
	case "<=":
		if strings.HasPrefix(right, "'") && strings.HasSuffix(right, "'") {
			r := right[1 : len(right)-1]
			return strconv.FormatBool(left <= r), nil
		} else {
			l, err := strconv.ParseFloat(left, 64)
			if err != nil {
				return "false", err
			}
			r, err := strconv.ParseFloat(right, 64)
			if err != nil {
				return "false", err
			}
			return strconv.FormatBool(l <= r), nil
		}
		//	case "LIKE":
		//		return fmt.Sprint("(", left, " LIKE '", right, "')"), nil
		//	case "NOT_LIKE":
		//		return fmt.Sprint("(", left, " NOT LIKE '", right, "')"), nil
		//	case "IS_NULL":
		//		return fmt.Sprint("(", left, " IS NULL)"), nil
		//	case "IS_NOT_NULL":
		//		return fmt.Sprint("(", left, " IS NOT NULL)"), nil
	case "RLIKE":
		if strings.HasPrefix(right, "'") && strings.HasSuffix(right, "'") {
			r := right[1 : len(right)-1]
			matches, err := regexp.MatchString(r, left)
			if err != nil {
				return "false", err
			}
			return strconv.FormatBool(matches), nil
		}
	}
	return "false", errors.New(fmt.Sprint("Failed to evaluate:", left, op, right))
}
