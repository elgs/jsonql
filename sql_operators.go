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
	v, found := symbolTable.(map[string]interface{})
	if !found {
		v = make(map[string]interface{})
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
	return jq.Query(token)
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
			sr, oksr := r.(string)
			if oksr {
				sl, oksl := l.(string)
				if oksl {
					return strconv.FormatBool(sl == sr), nil
				}
			}
			ir, okir := r.(int64)
			if okir {
				switch vl := l.(type) {
				case string:
					il, err := strconv.ParseInt(vl, 10, 0)
					if err != nil {
						return "false", nil
					}
					return strconv.FormatBool(il == ir), nil
				case int64:
					return strconv.FormatBool(vl == ir), nil
				case float64:
					return strconv.FormatBool(vl == float64(ir)), nil
				default:
					return "false", nil
				}
			}
			fr, okfr := r.(float64)
			if okfr {
				switch vl := l.(type) {
				case string:
					fl, err := strconv.ParseFloat(vl, 64)
					if err != nil {
						return "false", nil
					}
					return strconv.FormatBool(fl == fr), nil
				case int64:
					return strconv.FormatBool(float64(vl) == fr), nil
				case float64:
					return strconv.FormatBool(vl == fr), nil
				default:
					return "false", nil
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
			sr, oksr := r.(string)
			if oksr {
				sl, oksl := l.(string)
				if oksl {
					return strconv.FormatBool(sl != sr), nil
				}
			}
			ir, okir := r.(int64)
			if okir {
				switch vl := l.(type) {
				case string:
					il, err := strconv.ParseInt(vl, 10, 0)
					if err != nil {
						return "false", nil
					}
					return strconv.FormatBool(il != ir), nil
				case int64:
					return strconv.FormatBool(vl != ir), nil
				case float64:
					return strconv.FormatBool(vl != float64(ir)), nil
				default:
					return "false", nil
				}
			}
			fr, okfr := r.(float64)
			if okfr {
				switch vl := l.(type) {
				case string:
					fl, err := strconv.ParseFloat(vl, 64)
					if err != nil {
						return "false", nil
					}
					return strconv.FormatBool(fl != fr), nil
				case int64:
					return strconv.FormatBool(float64(vl) != fr), nil
				case float64:
					return strconv.FormatBool(vl != fr), nil
				default:
					return "false", nil
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
			sr, oksr := r.(string)
			if oksr {
				sl, oksl := l.(string)
				if oksl {
					return strconv.FormatBool(sl > sr), nil
				}
			}
			ir, okir := r.(int64)
			if okir {
				switch vl := l.(type) {
				case string:
					il, err := strconv.ParseInt(vl, 10, 0)
					if err != nil {
						return "false", nil
					}
					return strconv.FormatBool(il > ir), nil
				case int64:
					return strconv.FormatBool(vl > ir), nil
				case float64:
					return strconv.FormatBool(vl > float64(ir)), nil
				default:
					return "false", nil
				}
			}
			fr, okfr := r.(float64)
			if okfr {
				switch vl := l.(type) {
				case string:
					fl, err := strconv.ParseFloat(vl, 64)
					if err != nil {
						return "false", nil
					}
					return strconv.FormatBool(fl > fr), nil
				case int64:
					return strconv.FormatBool(float64(vl) > fr), nil
				case float64:
					return strconv.FormatBool(vl > fr), nil
				default:
					return "false", nil
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
			sr, oksr := r.(string)
			if oksr {
				sl, oksl := l.(string)
				if oksl {
					return strconv.FormatBool(sl < sr), nil
				}
			}
			ir, okir := r.(int64)
			if okir {
				switch vl := l.(type) {
				case string:
					il, err := strconv.ParseInt(vl, 10, 0)
					if err != nil {
						return "false", nil
					}
					return strconv.FormatBool(il < ir), nil
				case int64:
					return strconv.FormatBool(vl < ir), nil
				case float64:
					return strconv.FormatBool(vl < float64(ir)), nil
				default:
					return "false", nil
				}
			}
			fr, okfr := r.(float64)
			if okfr {
				switch vl := l.(type) {
				case string:
					fl, err := strconv.ParseFloat(vl, 64)
					if err != nil {
						return "false", nil
					}
					return strconv.FormatBool(fl < fr), nil
				case int64:
					return strconv.FormatBool(float64(vl) < fr), nil
				case float64:
					return strconv.FormatBool(vl < fr), nil
				default:
					return "false", nil
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
			sr, oksr := r.(string)
			if oksr {
				sl, oksl := l.(string)
				if oksl {
					return strconv.FormatBool(sl >= sr), nil
				}
			}
			ir, okir := r.(int64)
			if okir {
				switch vl := l.(type) {
				case string:
					il, err := strconv.ParseInt(vl, 10, 0)
					if err != nil {
						return "false", nil
					}
					return strconv.FormatBool(il >= ir), nil
				case int64:
					return strconv.FormatBool(vl >= ir), nil
				case float64:
					return strconv.FormatBool(vl >= float64(ir)), nil
				default:
					return "false", nil
				}
			}
			fr, okfr := r.(float64)
			if okfr {
				switch vl := l.(type) {
				case string:
					fl, err := strconv.ParseFloat(vl, 64)
					if err != nil {
						return "false", nil
					}
					return strconv.FormatBool(fl <= fr), nil
				case int64:
					return strconv.FormatBool(float64(vl) <= fr), nil
				case float64:
					return strconv.FormatBool(vl <= fr), nil
				default:
					return "false", nil
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
			sr, oksr := r.(string)
			if oksr {
				sl, oksl := l.(string)
				if oksl {
					return strconv.FormatBool(sl <= sr), nil
				}
			}
			ir, okir := r.(int64)
			if okir {
				switch vl := l.(type) {
				case string:
					il, err := strconv.ParseInt(vl, 10, 0)
					if err != nil {
						return "false", nil
					}
					return strconv.FormatBool(il <= ir), nil
				case int64:
					return strconv.FormatBool(vl <= ir), nil
				case float64:
					return strconv.FormatBool(vl <= float64(ir)), nil
				default:
					return "false", nil
				}
			}
			fr, okfr := r.(float64)
			if okfr {
				switch vl := l.(type) {
				case string:
					fl, err := strconv.ParseFloat(vl, 64)
					if err != nil {
						return "false", nil
					}
					return strconv.FormatBool(fl <= fr), nil
				case int64:
					return strconv.FormatBool(float64(vl) <= fr), nil
				case float64:
					return strconv.FormatBool(vl <= fr), nil
				default:
					return "false", nil
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
		Precedence: 7,
		Eval: func(symbolTable interface{}, left string, right string) (string, error) {
			l, err := evalToken(symbolTable, left)
			if err != nil {
				return "false", err
			}
			r, err := evalToken(symbolTable, right)
			if err != nil {
				return "false", err
			}
			il, okil := l.(int64)
			ir, okir := r.(int64)
			fl, okfl := l.(float64)
			fr, okfr := r.(float64)
			if okil && okir { //ii
				return fmt.Sprint(il + ir), nil
			} else if okfl && okfr { //ff
				return fmt.Sprint(fl + fr), nil
			} else if okil && okfr { //if
				return fmt.Sprint(float64(il) + fr), nil
			} else if okfl && okir { //fi
				return fmt.Sprint(fl + float64(ir)), nil
			} else { //else
				return fmt.Sprint("'", l, r, "'"), nil
			}
		},
	},
	"-": &Operator{
		Precedence: 7,
		Eval: func(symbolTable interface{}, left string, right string) (string, error) {
			l, err := evalToken(symbolTable, left)
			if err != nil {
				return "false", err
			}
			r, err := evalToken(symbolTable, right)
			if err != nil {
				return "false", err
			}
			il, okil := l.(int64)
			ir, okir := r.(int64)
			fl, okfl := l.(float64)
			fr, okfr := r.(float64)
			if okil && okir { //ii
				return fmt.Sprint(il - ir), nil
			} else if okfl && okfr { //ff
				return fmt.Sprint(fl - fr), nil
			} else if okil && okfr { //if
				return fmt.Sprint(float64(il) - fr), nil
			} else if okfl && okir { //fi
				return fmt.Sprint(fl - float64(ir)), nil
			} else { //else
				return "", errors.New(fmt.Sprint("Failed to evaluate: ", left, right))
			}
		},
	},
	"*": &Operator{
		Precedence: 9,
		Eval: func(symbolTable interface{}, left string, right string) (string, error) {
			l, err := evalToken(symbolTable, left)
			if err != nil {
				return "false", err
			}
			r, err := evalToken(symbolTable, right)
			if err != nil {
				return "false", err
			}
			il, okil := l.(int64)
			ir, okir := r.(int64)
			fl, okfl := l.(float64)
			fr, okfr := r.(float64)
			if okil && okir { //ii
				return fmt.Sprint(il * ir), nil
			} else if okfl && okfr { //ff
				return fmt.Sprint(fl * fr), nil
			} else if okil && okfr { //if
				return fmt.Sprint(float64(il) * fr), nil
			} else if okfl && okir { //fi
				return fmt.Sprint(fl * float64(ir)), nil
			} else { //else
				return "", errors.New(fmt.Sprint("Failed to evaluate: ", left, right))
			}
		},
	},
	"/": &Operator{
		Precedence: 9,
		Eval: func(symbolTable interface{}, left string, right string) (string, error) {
			l, err := evalToken(symbolTable, left)
			if err != nil {
				return "false", err
			}
			r, err := evalToken(symbolTable, right)
			if err != nil {
				return "false", err
			}
			il, okil := l.(int64)
			ir, okir := r.(int64)
			fl, okfl := l.(float64)
			fr, okfr := r.(float64)
			if ir == 0 || fr == 0 {
				return "", errors.New(fmt.Sprint("Failed to evaluate: ", left, right))
			}
			if okil && okir { //ii
				return fmt.Sprint(il / ir), nil
			} else if okfl && okfr { //ff
				return fmt.Sprint(fl / fr), nil
			} else if okil && okfr { //if
				return fmt.Sprint(float64(il) / fr), nil
			} else if okfl && okir { //fi
				return fmt.Sprint(fl / float64(ir)), nil
			} else { //else
				return "", errors.New(fmt.Sprint("Failed to evaluate: ", left, right))
			}
		},
	},
	"%": &Operator{
		Precedence: 9,
		Eval: func(symbolTable interface{}, left string, right string) (string, error) {
			l, err := evalToken(symbolTable, left)
			if err != nil {
				return "false", err
			}
			r, err := evalToken(symbolTable, right)
			if err != nil {
				return "false", err
			}
			il, okil := l.(int64)
			ir, okir := r.(int64)
			if ir == 0 {
				return "", errors.New(fmt.Sprint("Failed to evaluate: ", left, right))
			}
			if okil && okir { //ii
				return fmt.Sprint(il % ir), nil
			} else { //else
				return "", errors.New(fmt.Sprint("Failed to evaluate: ", left, right))
			}
		},
	},
	"^": &Operator{
		Precedence: 10,
		Eval: func(symbolTable interface{}, left string, right string) (string, error) {
			l, err := evalToken(symbolTable, left)
			if err != nil {
				return "false", err
			}
			r, err := evalToken(symbolTable, right)
			if err != nil {
				return "false", err
			}
			il, okil := l.(int64)
			ir, okir := r.(int64)
			fl, okfl := l.(float64)
			fr, okfr := r.(float64)
			if okil && okir { //ii
				return fmt.Sprint(math.Pow(float64(il), float64(ir))), nil
			} else if okfl && okfr { //ff
				return fmt.Sprint(math.Pow(fl, fr)), nil
			} else if okil && okfr { //if
				return fmt.Sprint(math.Pow(float64(il), fr)), nil
			} else if okfl && okir { //fi
				return fmt.Sprint(math.Pow(fl, float64(ir))), nil
			} else { //else
				return "", errors.New(fmt.Sprint("Failed to evaluate: ", left, right))
			}
		},
	},
}
