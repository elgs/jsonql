// exparser_test
package jsonql

import (
	"testing"
)

func TestTokenize(t *testing.T) {

	var pass = []struct {
		in string
		ex []string
	}{
		{"-1 + 2", []string{"-1", "+", "2"}},
		{"-(1+2)", []string{"-", "(", "1", "+", "2", ")"}},
		{"+1+2", []string{"+1", "+", "2"}},
		{"1+2+(3*4)", []string{"1", "+", "2", "+", "(", "3", "*", "4", ")"}},
		{"1+2+(3*4)^34", []string{"1", "+", "2", "+", "(", "3", "*", "4", ")", "^", "34"}},
		{"'123  456' 789", []string{"123  456", "789"}},
		{`123 "456  789"`, []string{"123", "456  789"}},
		{`123 "456  '''789"`, []string{"123", "456  '''789"}},
	}
	var fail = []struct {
		in string
		ex []string
	}{}
	parser := &Parser{
		Operators: SqlOperators,
	}
	for _, v := range pass {
		tokens := parser.Tokenize(v.in)
		if !CompareSlices(tokens, v.ex) {
			t.Error("Expected:", v.ex, len(v.ex), "actual:", tokens, len(tokens))
		}
	}
	for _, v := range fail {
		tokens := parser.Tokenize(v.in)
		if CompareSlices(tokens, v.ex) {
			t.Error("Expected:", v.ex, "actual:", tokens)
		}
	}
}

func TestTokenizeSql(t *testing.T) {
	var pass = []struct {
		in string
		ex []string
	}{
		{"true AND false", []string{"true", "AND", "false"}},
	}
	var fail = []struct {
		in string
		ex []string
	}{}
	parser := &Parser{
		Operators: SqlOperators,
	}
	for _, v := range pass {
		tokens := parser.Tokenize(v.in)
		if !CompareSlices(tokens, v.ex) {
			t.Error("Expected:", v.ex, "actual:", tokens)
		}
	}
	for _, _ = range fail {

	}
}

func TestRPN(t *testing.T) {
	var pass = []struct {
		in string
		ex string
	}{
		{"true AND false", "false"},
		{"true and true", "true"},
		{"false and false", "false"},
		{"true OR true", "true"},
		{"true OR (true and false)", "true"},
		{"true and (false and true)", "false"},
		{"(true and false) or (false or true)", "true"},
	}
	var fail = []struct {
		in string
		ex []string
	}{}
	parser := &Parser{
		Operators: SqlOperators,
	}

	for _, v := range pass {
		result, err := parser.Calculate(v.in)
		if err != nil {
			t.Error(err)
		}
		if result != v.ex {
			t.Error("Expected:", v.ex, "actual:", result)
		}
	}
	for _, _ = range fail {

	}
}
