// jsonql
package jsonql

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type JSONQL struct {
	Data interface{}
}

func NewStringQuery(jsonString string) (*JSONQL, error) {
	var data = new(interface{})
	err := json.Unmarshal([]byte(jsonString), data)
	if err != nil {
		return nil, err
	}
	return &JSONQL{*data}, nil
}

func NewQuery(jsonObject interface{}) *JSONQL {
	return &JSONQL{jsonObject}
}

func (this *JSONQL) Query(where string) (interface{}, error) {
	parser := &Parser{
		Operators: SqlOperators,
	}
	tokens := parser.Tokenize(where)
	rpn, err := parser.ParseRPN(tokens)
	if err != nil {
		return nil, err
	}

	switch v := this.Data.(type) {
	case []interface{}:
		ret := []interface{}{}
		for _, obj := range v {
			parser.SymbolTable = obj
			r, err := this.processObj(parser, *rpn)
			if err != nil {
				return nil, err
			}
			//			fmt.Println("***")
			//			fmt.Println(obj, where, r)
			if r {
				ret = append(ret, obj)
			}
		}
		return ret, nil
	case map[string]interface{}:
		parser.SymbolTable = v
		return this.processObj(parser, *rpn)
	default:
		return nil, errors.New(fmt.Sprintf("Failed to parse input data."))
	}
}

func (this *JSONQL) processObj(parser *Parser, rpn Lifo) (bool, error) {
	result, err := parser.Evaluate(&rpn, true)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(result)
}
