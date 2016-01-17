package jsonql

import (
	"fmt"
	"testing"
)

func TestLifo(t *testing.T) {
	lifo := &Lifo{}
	lifo.Push("1")
	lifo.Push("2")
	lifo.Push("3")
	fmt.Println(lifo.Len())
	lifo.Print()
}
