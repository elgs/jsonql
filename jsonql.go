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
	//fmt.Println(expression, tokens)
	rpn, err := parser.ParseRPN(tokens)
	if err != nil {
		return nil, err
	}

	switch v := this.JQ.Data.(type) {
	case []interface{}:
		for _, obj := range v {
			this.processObj(obj, parser, rpn)
		}
	case map[string]interface{}:
		this.processObj(v, parser, rpn)
	default:
		return nil, errors.New(fmt.Sprintf("Failed to parse input data."))
	}
	return nil, nil
}

func (this *JSONQL) processObj(obj interface{}, parser *Parser, rpn *Lifo) (bool, error) {
	fmt.Println(obj)
	result, err := parser.Evaluate(rpn, true)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(result)
}
