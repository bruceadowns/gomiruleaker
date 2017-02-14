package lib

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestInputValid(t *testing.T) {
	in := &Input{
		Target:      "foo",
		ParserCount: 1,
	}
	if err := in.IsValid(); err != nil {
		t.Errorf("Expected valid input")
	}
}

func TestInputInvalid(t *testing.T) {
	in := &Input{
		DownloadDelayMs: 1,
		AccumBatchSize:  2,
		ParserCount:     3,
		Debug:           false,
	}
	if err := in.IsValid(); err == nil {
		t.Errorf("Expected invalid input")
	}
}

func TestInputString(t *testing.T) {
	in := &Input{
		Target:          "foo",
		MiruURL:         "bar",
		DownloadDelayMs: 1,
		AccumBatchSize:  2,
		ParserCount:     3,
		Debug:           true,
	}
	out := fmt.Sprintf("%s", in)
	if !strings.EqualFold("target: foo miru: bar delay: 1 batch: 2 count: 3 debug: true", out) {
		t.Errorf("Unexpected input %s", out)
	}
}

func TestInitInputSuccess(t *testing.T) {
	test := `{
	  "target": "https://foo.com/get/123",
	  "miruUrl": "https://host:port/add",
	  "downloadDelayMs": 500,
	  "accumBatchSize": 10,
	  "parserCount": 2,
	  "postErrorDelayMs": 1000,
	  "debug": true
	}`

	in, err := InitInput(bytes.NewBufferString(test))
	if err != nil {
		t.Error(err)
	}

	if !strings.EqualFold(in.Target, "https://foo.com/get/123") {
		t.Error("Target invalid")
	}
	if !strings.EqualFold(in.MiruURL, "https://host:port/add") {
		t.Error("MiruURL invalid")
	}
	if in.DownloadDelayMs != 500 {
		t.Error("DownloadDelayMs invalid")
	}
	if in.AccumBatchSize != 10 {
		t.Error("AccumBatchSize invalid")
	}
	if in.ParserCount != 2 {
		t.Error("ParserCount invalid")
	}
	if in.PostErrorDelayMs != 1000 {
		t.Error("PostErrorDelayMs invalid")
	}
	if !in.Debug {
		t.Error("Debug invalid")
	}
}

func TestInitInputMinimum(t *testing.T) {
	test := `{"target":"https://foo.com/get/123"}`
	in, err := InitInput(bytes.NewBufferString(test))
	if err != nil {
		t.Error(err)
		return
	}

	if !strings.EqualFold(in.Target, "https://foo.com/get/123") {
		t.Error("Target invalid")
	}
}

func TestInitInputInvalid(t *testing.T) {
	test := `{foo:bar}`

	_, err := InitInput(bytes.NewBufferString(test))
	if err == nil {
		t.Error("Expected decode error")
	}
}

func TestInitInputDecodeError(t *testing.T) {
	test := `{}`

	_, err := InitInput(bytes.NewBufferString(test))
	if err == nil {
		t.Error("Expected invalid error")
	}
}
