// jsonql
package jsonql

import (
	"errors"
	"fmt"
	"github.com/elgs/gojq"
	"strconv"
)

type JSONQL struct {
	JQ *gojq.JQ
}

func NewStringQuery(jsonString string) (*JSONQL, error) {
	jq, err := gojq.NewStringQuery(jsonString)
	if err != nil {
		return nil, err
	}
	return &JSONQL{jq}, nil
}

func NewQuery(jsonObject interface{}) *JSONQL {
	return &JSONQL{gojq.NewQuery(jsonObject)}
}

func (this *JSONQL) Parse(where string) (interface{}, error) {
	parser := &Parser{
		Operators: SqlOperators,
	}
	tokens := parser.Tokenize(where)
	rpn, err := parser.ParseRPN(tokens)
	if err != nil {
		return nil, err
	}

	switch v := this.JQ.Data.(type) {
	case []interface{}:
		ret := []interface{}{}
		for _, obj := range v {
			parser.SymbolTable = obj
			r, err := this.processObj(parser, rpn)
			if err != nil {
				return nil, err
			}
			if r {
				ret = append(ret, obj)
			}
		}
		return ret, nil
	case map[string]interface{}:
		parser.SymbolTable = v
		return this.processObj(parser, rpn)
	default:
		return nil, errors.New(fmt.Sprintf("Failed to parse input data."))
	}
	return nil, nil
}

func (this *JSONQL) processObj(parser *Parser, rpn *Lifo) (bool, error) {
	result, err := parser.Evaluate(rpn, true)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(result)
}
