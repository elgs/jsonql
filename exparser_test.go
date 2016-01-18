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
		{"'123  456' 789", []string{"'123  456'", "789"}},
		{`123 "456  789"`, []string{"123", "\"456  789\""}},
		{`123 "456  '''789"`, []string{"123", "\"456  '''789\""}},
	}
	var fail = []struct {
		in string
		ex []string
	}{}
	parser := &Parser{
		Operators: sqlOperators,
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
		Operators: sqlOperators,
	}
	for _, v := range pass {
		tokens := parser.Tokenize(v.in)
		if !CompareSlices(tokens, v.ex) {
			t.Error("Expected:", v.ex, "actual:", tokens)
		}
	}
	for range fail {

	}
}

func TestRPN(t *testing.T) {
	var pass = []struct {
		in string
		ex string
	}{
		{"true && false", "false"},
		{"true && true", "true"},
		{"false && false", "false"},
		{"true || true", "true"},
		{"true || (true && false)", "true"},
		{"true && (false && true)", "false"},
		{"(true && false) || (false || true)", "true"},
	}
	var fail = []struct {
		in string
		ex []string
	}{}
	parser := &Parser{
		Operators: sqlOperators,
	}

	for _, v := range pass {
		result, err := parser.Calculate(v.in)
		if err != nil {
			t.Error(err)
		}
		if result != v.ex {
			t.Error(v.in, "Expected:", v.ex, "actual:", result)
		}
	}
	for range fail {

	}
}
