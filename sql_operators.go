// mysql_operators
package jsonql

import (
	"errors"
	"fmt"
	"github.com/elgs/gojq"
	"math"
	"regexp"
	"strconv"
	"strings"
)

func evalToken(symbolTable interface{}, token string) (interface{}, error) {
	//	fmt.Println("******************")
	//	fmt.Println(symbolTable)
	//	fmt.Println(token)
	//	fmt.Println("******************")
	v, found := symbolTable.(map[string]interface{})
	if !found {
		return nil, errors.New(fmt.Sprint("Failed to parse token: ", token))
	}
	if token == "true" || token == "false" {
		return token, nil
	}
	if (strings.HasPrefix(token, "'") && strings.HasSuffix(token, "'")) ||
		(strings.HasPrefix(token, "\"") && strings.HasSuffix(token, "\"")) {
		// string
		return token[1 : len(token)-1], nil
	}
	intToken, err := strconv.ParseInt(token, 10, 0)
	if err == nil {
		return intToken, nil
	}
	floatToken, err := strconv.ParseFloat(token, 64)
	if err == nil {
		return floatToken, nil
	}
	jq := gojq.NewQuery(v)
	return jq.Parse(token)
}

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
			return strconv.FormatBool(l || r), nil
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
			return strconv.FormatBool(l && r), nil
		},
	},
	"=": &Operator{
		Precedence: 5,
		Eval: func(symbolTable interface{}, left string, right string) (string, error) {
			l, err := evalToken(symbolTable, left)
			if err != nil {
				return "false", err
			}
			r, err := evalToken(symbolTable, right)
			if err != nil {
				return "false", err
			}
			sl, foundl := l.(string)
			sr, foundr := r.(string)
			if foundl && foundr {
				return strconv.FormatBool(sl == sr), nil
			}
			fl, foundl := l.(float64)
			if foundl {
				fr, foundr := r.(float64)
				if foundr {
					return strconv.FormatBool(fl == fr), nil
				}
				ir, foundr := r.(int64)
				if foundr {
					return strconv.FormatBool(fl == float64(ir)), nil
				}
			}
			return "false", errors.New(fmt.Sprint("Failed to compare: ", left, right))
		},
	},
	"!=": &Operator{
		Precedence: 5,
		Eval: func(symbolTable interface{}, left string, right string) (string, error) {
			l, err := evalToken(symbolTable, left)
			if err != nil {
				return "false", err
			}
			r, err := evalToken(symbolTable, right)
			if err != nil {
				return "false", err
			}
			sl, foundl := l.(string)
			sr, foundr := r.(string)
			if foundl && foundr {
				return strconv.FormatBool(sl != sr), nil
			}
			fl, foundl := l.(float64)
			if foundl {
				fr, foundr := r.(float64)
				if foundr {
					return strconv.FormatBool(fl != fr), nil
				}
				ir, foundr := r.(int64)
				if foundr {
					return strconv.FormatBool(fl != float64(ir)), nil
				}
			}
			return "false", errors.New(fmt.Sprint("Failed to compare: ", left, right))
		},
	},
	">": &Operator{
		Precedence: 5,
		Eval: func(symbolTable interface{}, left string, right string) (string, error) {
			l, err := evalToken(symbolTable, left)
			if err != nil {
				return "false", err
			}
			r, err := evalToken(symbolTable, right)
			if err != nil {
				return "false", err
			}
			sl, foundl := l.(string)
			sr, foundr := r.(string)
			if foundl && foundr {
				return strconv.FormatBool(sl > sr), nil
			}
			fl, foundl := l.(float64)
			if foundl {
				fr, foundr := r.(float64)
				if foundr {
					return strconv.FormatBool(fl > fr), nil
				}
				ir, foundr := r.(int64)
				if foundr {
					return strconv.FormatBool(fl > float64(ir)), nil
				}
			}
			return "false", errors.New(fmt.Sprint("Failed to compare: ", left, right))
		},
	},
	"<": &Operator{
		Precedence: 5,
		Eval: func(symbolTable interface{}, left string, right string) (string, error) {
			l, err := evalToken(symbolTable, left)
			if err != nil {
				return "false", err
			}
			r, err := evalToken(symbolTable, right)
			if err != nil {
				return "false", err
			}
			sl, foundl := l.(string)
			sr, foundr := r.(string)
			if foundl && foundr {
				return strconv.FormatBool(sl < sr), nil
			}
			fl, foundl := l.(float64)
			if foundl {
				fr, foundr := r.(float64)
				if foundr {
					return strconv.FormatBool(fl < fr), nil
				}
				ir, foundr := r.(int64)
				if foundr {
					return strconv.FormatBool(fl < float64(ir)), nil
				}
			}
			return "false", errors.New(fmt.Sprint("Failed to compare: ", left, right))
		},
	},
	">=": &Operator{
		Precedence: 5,
		Eval: func(symbolTable interface{}, left string, right string) (string, error) {
			l, err := evalToken(symbolTable, left)
			if err != nil {
				return "false", err
			}
			r, err := evalToken(symbolTable, right)
			if err != nil {
				return "false", err
			}
			sl, foundl := l.(string)
			sr, foundr := r.(string)
			if foundl && foundr {
				return strconv.FormatBool(sl >= sr), nil
			}
			fl, foundl := l.(float64)
			if foundl {
				fr, foundr := r.(float64)
				if foundr {
					return strconv.FormatBool(fl >= fr), nil
				}
				ir, foundr := r.(int64)
				if foundr {
					return strconv.FormatBool(fl >= float64(ir)), nil
				}
			}
			return "false", errors.New(fmt.Sprint("Failed to compare: ", left, right))
		},
	},
	"<=": &Operator{
		Precedence: 5,
		Eval: func(symbolTable interface{}, left string, right string) (string, error) {
			l, err := evalToken(symbolTable, left)
			if err != nil {
				return "false", err
			}
			r, err := evalToken(symbolTable, right)
			if err != nil {
				return "false", err
			}
			sl, foundl := l.(string)
			sr, foundr := r.(string)
			if foundl && foundr {
				return strconv.FormatBool(sl <= sr), nil
			}
			fl, foundl := l.(float64)
			if foundl {
				fr, foundr := r.(float64)
				if foundr {
					return strconv.FormatBool(fl <= fr), nil
				}
				ir, foundr := r.(int64)
				if foundr {
					return strconv.FormatBool(fl <= float64(ir)), nil
				}
			}
			return "false", errors.New(fmt.Sprint("Failed to compare: ", left, right))
		},
	},
	"RLIKE": &Operator{
		Precedence: 5,
		Eval: func(symbolTable interface{}, left string, right string) (string, error) {
			l, err := evalToken(symbolTable, left)
			if err != nil {
				return "false", err
			}
			r, err := evalToken(symbolTable, right)
			if err != nil {
				return "false", err
			}
			sl, foundl := l.(string)
			sr, foundr := r.(string)
			if foundl && foundr {
				matches, err := regexp.MatchString(sr, sl)
				if err != nil {
					return "false", err
				}
				return strconv.FormatBool(matches), nil
			}
			return "false", errors.New(fmt.Sprint("Failed to compare: ", left, right))

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
