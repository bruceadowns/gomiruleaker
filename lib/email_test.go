package lib

import (
	"testing"
	"time"
)

func TestDateFormatterSuccess(t *testing.T) {
	tests := []string{
		"Wed, 04 Apr 2012 11:41:06 -0400",
		"Sun, 10 May 2015 09:13:58 -0400",
	}

	for _, test := range tests {
		if d, err := time.Parse(time.RFC1123Z, test); err == nil {
			t.Logf("Success: %v", d)
		} else {
			t.Errorf("Error parsing email date: %s [%s]", test, err)
		}
	}
}

func TestDateFormatterFailure(t *testing.T) {
	tests := []string{
		"Wed, 4 Apr 2012 11:41:06 -0400",
	}

	for _, test := range tests {
		if _, err := time.Parse(time.RFC1123Z, test); err == nil {
			t.Errorf("Error expecting failure: %s", test)
		} else {
			t.Logf("Failure as expected")
		}
	}
}
