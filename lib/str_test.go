package lib

import (
	"strings"
	"testing"
)

func TestTrunc(t *testing.T) {
	testTable := []struct {
		in       string
		expected string
	}{
		{
			in:       "",
			expected: "",
		}, {
			in:       "1234567890",
			expected: "1234567890",
		}, {
			in:       "123456789012345678901234567890123456789012345678901234567890123456789012345678901",
			expected: "12345678901234567890123456789012345678901234567890123456789012345678901234567890...",
		}}

	for _, test := range testTable {
		out := Trunc(test.in)
		if !strings.EqualFold(test.expected, out) {
			t.Errorf("Error truncating: %s", test.in)
		}
	}
}
