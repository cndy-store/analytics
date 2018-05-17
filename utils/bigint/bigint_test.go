package bigint

import (
	"testing"
)

func TestParse(t *testing.T) {
	var tests = []struct {
		String string
		Int    int64
	}{
		{"", 0},
		{"5.0000000", 50000000},
		{"5.0", 50000000},
		{"5", 50000000},
		{"50000000", 500000000000000},
	}

	for _, test := range tests {
		res, err := Parse(test.String)
		if err != nil {
			t.Error(err)
		}

		if test.Int != *res {
			t.Errorf("Expected %d got %d", test.Int, *res)
		}
	}
}

func TestToString(t *testing.T) {
	var tests = []struct {
		Int    int64
		String string
	}{
		{0, "0.0000000"},
		{5, "0.0000005"},
		{50000000, "5.0000000"},
	}

	for _, test := range tests {
		res := ToString(test.Int)
		if test.String != res {
			t.Errorf("Expected %s got %s", test.String, res)
		}
	}
}
