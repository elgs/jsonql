package jsonql

import (
	"testing"
)

func TestCompareSlices(t *testing.T) {
	var pass = []struct{ in, ex []string }{
		{[]string{"1", "2", "3"}, []string{"1", "2", "3"}},
		{[]string{}, []string{}},
		{[]string{"a", "b", "你好"}, []string{"a", "b", "你好"}},
	}
	var fail = []struct{ in, ex []string }{
		{[]string{"1"}, []string{}},
		{[]string{}, []string{"1", "2"}},
	}
	for _, v := range pass {
		if !CompareSlices(v.in, v.ex) {
			t.Error("Expected:", v.ex, "actual:", v.in)
		}
	}
	for _, v := range fail {
		if CompareSlices(v.in, v.ex) {
			t.Error("Expected:", v.ex, "actual:", v.in)
		}
	}
}

func TestReverseString(t *testing.T) {
	var pass = []struct{ in, ex string }{
		{"", ""},
		{"Hello, 世界", "界世 ,olleH"},
		{"Backward", "drawkcaB"},
	}
	var fail = []struct{ in, ex string }{
		{"", " "},
		{" ", ""},
		{"111", "112"},
	}
	for _, v := range pass {
		if v.ex != ReverseString(v.in) {
			t.Error("Expected:", v.ex, "actual:", v.in)
		}
	}
	for _, v := range fail {
		if v.ex == ReverseString(v.in) {
			t.Error("Expected:", v.ex, "actual:", v.in)
		}
	}
}
