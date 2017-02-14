package lib

import (
	"math"
	"testing"
)

func TestSuccessExpandTargetURL(t *testing.T) {
	testTable := []struct {
		in       string
		expected *Targets
	}{
		{
			in:       "https://foo.bar:123/bad-emails/get/12345",
			expected: &Targets{"https://foo.bar:123/bad-emails/get/", 12345, 12345},
		},
		{
			in:       "https://foo.bar:123/bad-emails/get/[12345]",
			expected: &Targets{"https://foo.bar:123/bad-emails/get/", 12345, 12345},
		},
		{
			in:       "https://foo.bar:123/bad-emails/get/[:]",
			expected: &Targets{"https://foo.bar:123/bad-emails/get/", 1, math.MaxInt32},
		},
		{
			in:       "https://foo.bar:123/bad-emails/get/[123:]",
			expected: &Targets{"https://foo.bar:123/bad-emails/get/", 123, math.MaxInt32},
		},
		{
			in:       "https://foo.bar:123/bad-emails/get/[:456]",
			expected: &Targets{"https://foo.bar:123/bad-emails/get/", 1, 456},
		},
		{
			in:       "https://foo.bar:123/bad-emails/get/[123:456]",
			expected: &Targets{"https://foo.bar:123/bad-emails/get/", 123, 456},
		},
		{
			in:       "https://foo.bar:123/bad-emails/get/[456:123]",
			expected: &Targets{"https://foo.bar:123/bad-emails/get/", 123, 456},
		},
	}

	for _, test := range testTable {
		result, err := ExpandTargetURL(test.in)
		t.Logf("result: %s", result)
		if err != nil {
			t.Errorf("Error expanding url: %s", err)
		}
		if !test.expected.Equals(result) {
			t.Errorf("Error expanding url: %s", test.in)
		}
	}
}

func TestFailExpandTargetURL(t *testing.T) {
	testTable := []string{
		"https://foo.bar:123/bad-emails/get/",
		"https://foo.bar:123/bad-emails/get/abc",
		"https://foo.bar:123/bad-emails/get/[]",
		"https://foo.bar:123/bad-emails/get/[abc]",
		"https://foo.bar:123/bad-emails/get/[a:b]",
		"https://foo.bar:123/bad-emails/get/[1:2:3]",
		"https://foo.bar:123/bad-emails/get/[123",
		"https://foo.bar:123/bad-emails/get/123]",
		"https://foo.bar:123/bad-emails/get/[12]3]",
		"https://foo.bar:123/bad-emails/get/[1[2]3]]",
	}

	for _, test := range testTable {
		_, err := ExpandTargetURL(test)
		t.Logf("Error: %s", err)
		if err == nil {
			t.Errorf("Expecting Error: %s", test)
		}
	}
}
