// mysql_operators
package jsonql

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var SqlOperators = map[string]*Operator{
	"OR": &Operator{
		Precedence: 1,
		Eval: func(symbolTable interface{}, left string, right string) (string, error) {
			l, err := strconv.ParseBool(left)
			if err != nil {
				return "false", err
			}
			r, err := strconv.ParseBool(right)
			if err != nil {
				return "false", err
			}
			return strconv.FormatBool(l && r), nil
		},
	},
	"AND": &Operator{
		Precedence: 3,
		Eval: func(symbolTable interface{}, left string, right string) (string, error) {
			l, err := strconv.ParseBool(left)
			if err != nil {
				return "false", err
			}
			r, err := strconv.ParseBool(right)
			if err != nil {
				return "false", err
			}
			return strconv.FormatBool(l || r), nil
		},
	},
	"=": &Operator{
		Precedence: 5,
		Eval: func(symbolTable interface{}, left string, right string) (string, error) {
			l, err := strconv.ParseBool(left)
			if err != nil {
				return "false", err
			}
			r, err := strconv.ParseBool(right)
			if err != nil {
				return "false", err
			}
			return strconv.FormatBool(l == r), nil
		},
	},
	"!=": &Operator{
		Precedence: 5,
		Eval: func(symbolTable interface{}, left string, right string) (string, error) {
			return strconv.FormatBool(left != right), nil
		},
	},
	">": &Operator{
		Precedence: 5,
		Eval: func(symbolTable interface{}, left string, right string) (string, error) {
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
		},
	},
	"<": &Operator{
		Precedence: 5,
		Eval: func(symbolTable interface{}, left string, right string) (string, error) {
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
		},
	},
	">=": &Operator{
		Precedence: 5,
		Eval: func(symbolTable interface{}, left string, right string) (string, error) {
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
		},
	},
	"<=": &Operator{
		Precedence: 5,
		Eval: func(symbolTable interface{}, left string, right string) (string, error) {
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
		},
	},
	//	"LIKE": &Operator{
	//		Precedence: 5,
	//	},
	//	"NOT_LIKE": &Operator{
	//		Precedence: 5,
	//	},
	//	"IS_NULL": &Operator{
	//		Precedence: 5,
	//	},
	//	"IS_NOT_NULL": &Operator{
	//		Precedence: 5,
	//	},
	"RLIKE": &Operator{
		Precedence: 5,
		Eval: func(symbolTable interface{}, left string, right string) (string, error) {
			if strings.HasPrefix(right, "'") && strings.HasSuffix(right, "'") {
				r := right[1 : len(right)-1]
				matches, err := regexp.MatchString(r, left)
				if err != nil {
					return "false", err
				}
				return strconv.FormatBool(matches), nil
			} else {
				return "", errors.New(fmt.Sprint("Failed to evaluate:", left, "RLIKE", right))
			}
		},
	},
	"+": &Operator{
		Precedence: 1,
		Eval: func(symbolTable interface{}, left string, right string) (string, error) {
			isDec := strings.Contains(left, ".") || strings.Contains(right, ".")
			if isDec {
				l, err := strconv.ParseFloat(left, 64)
				r, err := strconv.ParseFloat(right, 64)
				return fmt.Sprint(l + r), err
			} else {
				l, err := strconv.ParseInt(left, 10, 64)
				r, err := strconv.ParseInt(right, 10, 64)
				return fmt.Sprint(l + r), err
			}
		},
	},
	"-": &Operator{
		Precedence: 1,
		Eval: func(symbolTable interface{}, left string, right string) (string, error) {
			isDec := strings.Contains(left, ".") || strings.Contains(right, ".")
			if isDec {
				l, err := strconv.ParseFloat(left, 64)
				r, err := strconv.ParseFloat(right, 64)
				return fmt.Sprint(l - r), err
			} else {
				l, err := strconv.ParseInt(left, 10, 64)
				r, err := strconv.ParseInt(right, 10, 64)
				return fmt.Sprint(l - r), err
			}
		},
	},
	"*": &Operator{
		Precedence: 3,
		Eval: func(symbolTable interface{}, left string, right string) (string, error) {
			isDec := strings.Contains(left, ".") || strings.Contains(right, ".")
			if isDec {
				l, err := strconv.ParseFloat(left, 64)
				r, err := strconv.ParseFloat(right, 64)
				return fmt.Sprint(l * r), err
			} else {
				l, err := strconv.ParseInt(left, 10, 64)
				r, err := strconv.ParseInt(right, 10, 64)
				return fmt.Sprint(l * r), err
			}
		},
	},
	"/": &Operator{
		Precedence: 3,
		Eval: func(symbolTable interface{}, left string, right string) (string, error) {
			isDec := strings.Contains(left, ".") || strings.Contains(right, ".")
			if isDec {
				l, err := strconv.ParseFloat(left, 64)
				r, err := strconv.ParseFloat(right, 64)
				if r == 0 {
					return "", errors.New(fmt.Sprint("Failed to evaluate:", left, "/", right))
				}
				return fmt.Sprint(l / r), err
			} else {
				l, err := strconv.ParseInt(left, 10, 64)
				r, err := strconv.ParseInt(right, 10, 64)
				if r == 0 {
					return "", errors.New(fmt.Sprint("Failed to evaluate:", left, "/", right))
				}
				return fmt.Sprint(l / r), err
			}
		},
	},
	"%": &Operator{
		Precedence: 3,
		Eval: func(symbolTable interface{}, left string, right string) (string, error) {
			isDec := strings.Contains(left, ".") || strings.Contains(right, ".")
			if isDec {
				return "", errors.New(fmt.Sprint("Failed to evaluate:", left, "/", right))
			} else {
				l, err := strconv.ParseInt(left, 10, 64)
				r, err := strconv.ParseInt(right, 10, 64)
				if r == 0 {
					return "", errors.New(fmt.Sprint("Failed to evaluate:", left, "/", right))
				}
				return fmt.Sprint(l % r), err
			}
		},
	},
	"^": &Operator{
		Precedence: 4,
		Eval: func(symbolTable interface{}, left string, right string) (string, error) {
			isDec := strings.Contains(left, ".") || strings.Contains(right, ".")
			l, err := strconv.ParseFloat(left, 64)
			r, err := strconv.ParseFloat(right, 64)
			if isDec {
				return fmt.Sprint(math.Pow(l, r)), err
			} else {
				return fmt.Sprint(int64(math.Pow(l, r))), err
			}
		},
	},
}
