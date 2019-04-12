package jsonql

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {

	jsonString := `
[
  {
    "name": "elgs",
    "gender": "m",
    "skills": [
      "Golang",
      "Java",
      "C"
	],
	"comments": "",
	"score": 40
  },
  {
    "name": "enny",
    "gender": "f",
    "skills": [
      "IC",
      "Electric design",
      "Verification"
    ],
	"comments": "2nd record",
	"score": 10

  },
  {
    "name": "sam",
    "gender": "m",
    "skills": [
      "Eating",
      "Sleeping",
      "Crawling"
		],
	"comments": "3rd record",
	"score": 30
  }
]
`
	parser, err := NewStringQuery(jsonString)
	if err != nil {
		t.Error(err)
	}

	var pass = []struct {
		in string
		ex int
	}{
		{"name='elgs'", 1},
		{"gender='f'", 1},
		{"comments='2nd record'", 1},
		{"comments=\"\"", 1},
		{"comments=\"\"", 1},
		{"skills.[1]='Sleeping'", 1},
		{"skills contains 'Sleeping'", 1},
		{"skills contains 'Awake'", 0},
		{"name contains 'lgs'", 1},
		{"name contains 'e'", 2},
		{"'Sleeping' in skills", 1},
		{"'Verification' notin skills", 2},
		{"'Awake' in skills", 0},
		{"'lgs' in name", 1},
		{"'e' in name", 2},
		{"score > 10 && score < 40", 1},
	}
	var fail = []struct {
		in string
		ex interface{}
	}{}
	for _, v := range pass {
		result, err := parser.Query(v.in)
		if err != nil {
			t.Error(v.in, err)
		}
		fmt.Println(v.in, result)
		numResults, ok := result.([]interface{})
		if !ok {
			t.Error("Failed to convert result")
		}
		if len(numResults) != v.ex {
			t.Errorf("Got %v results, expected %v for `%v`", len(numResults), v.ex, v.in)
		}
	}
	for range fail {

	}
}
