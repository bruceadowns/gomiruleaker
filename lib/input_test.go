package lib

import (
	"bytes"
	"strings"
	"testing"
)

func TestInitInputSuccess(t *testing.T) {
	test := `
miruUrl: https://host:port/add
outputDir: ./dump
downloadDelayMs: 500
accumBatchSize: 10
parserCount: 2
postErrorDelayMs: 1000
postCompress: true
debug: true

targets:
  - type: one
    subType: subone
    start: 123
    end: 124

  - type: two
    subType: subtwo
    start: 125
    end: 126
    limit: 10
`

	in, err := InitInput(bytes.NewBufferString(test))
	if err != nil {
		t.Fatal(err)
	}

	if len(in.Targets) != 2 {
		t.Error("Invalid target count")
		return
	}

	if !strings.EqualFold(in.Targets[0].Type, "one") {
		t.Error("Target type invalid")
	}
	if !strings.EqualFold(in.Targets[0].SubType, "subone") {
		t.Error("Target subtype invalid")
	}
	if in.Targets[0].Start != 123 {
		t.Error("Target start invalid")
	}
	if in.Targets[0].End != 124 {
		t.Error("Target end invalid")
	}
	if in.Targets[0].Limit != 100 {
		t.Errorf("Target limit %d invalid", in.Targets[0].Limit)
	}

	if !strings.EqualFold(in.Targets[1].Type, "two") {
		t.Error("Target type invalid")
	}
	if !strings.EqualFold(in.Targets[1].SubType, "subtwo") {
		t.Error("Target subtype invalid")
	}
	if in.Targets[1].Start != 125 {
		t.Error("Target start invalid")
	}
	if in.Targets[1].End != 126 {
		t.Error("Target end invalid")
	}
	if in.Targets[1].Limit != 10 {
		t.Error("Target limit invalid")
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
	if !in.PostCompress {
		t.Error("PostCompress invalid")
	}
	if !in.Debug {
		t.Error("Debug invalid")
	}
}
