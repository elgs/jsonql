// jsonql
package jsonql

import (
	"fmt"
)

func main() {
	fmt.Println("Hello World!")
}

type JSONQL struct {
}

func NewStringQuery(jsonString string) *JSONQL {
	return &JSONQL{}
}

func NewObjectQuery(jsonObject map[string]interface{}) *JSONQL {
	return &JSONQL{}
}

func NewArrayQuery(jsonArray []interface{}) *JSONQL {
	return &JSONQL{}
}
