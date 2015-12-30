// jsonql
package jsonql

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elgs/gosplitargs"
	"strconv"
	"strings"
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
	return &JSONQL{Data: jsonObject}
}

func (this *JSONQL) Parse(exp string) (interface{}, error) {
	paths, err := gosplitargs.SplitArgs(exp, "\\.", false)
	if err != nil {
		return nil, err
	}
	var context interface{} = this.Data
	for _, path := range paths {
		if len(path) >= 3 && strings.HasPrefix(path, "[") && strings.HasSuffix(path, "]") {
			// array
			index, err := strconv.Atoi(path[1 : len(path)-1])
			if err != nil {
				return nil, err
			}
			switch v := context.(type) {
			case []interface{}:
				{
					if len(v) <= index {
						return nil, errors.New("Index out of range.")
					}
					context = v[index]
				}
			default:
				return nil, errors.New(fmt.Sprint(path, " is not an array. ", v))
			}

		} else {
			// map
			switch v := context.(type) {
			case map[string]interface{}:
				context = v[path]
			default:
				return nil, errors.New(fmt.Sprint(path, " is not an object. ", v))
			}
		}
	}
	return context, nil
}
