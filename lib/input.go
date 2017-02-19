package lib

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"runtime"

	"gopkg.in/yaml.v2"
)

// InputTarget ...
type InputTarget struct {
	Type    string `yaml:"type"`
	SubType string `yaml:"subType"`
	Start   int    `yaml:"start"`
	End     int    `yaml:"end"`
	Limit   int    `yaml:"limit"`
}

func (i InputTarget) String() string {
	return fmt.Sprintf("type: %s subtype: %s start: %d end: %d limit: %d",
		i.Type, i.SubType, i.Start, i.End, i.Limit)
}

// Input ...
type Input struct {
	Targets          []*InputTarget
	MiruURL          string `yaml:"miruUrl"`
	OutputDir        string `yaml:"outputDir"`
	DownloadDelayMs  int    `yaml:"downloadDelayMs"`
	AccumBatchSize   int    `yaml:"accumBatchSize"`
	ParserCount      int    `yaml:"parserCount"`
	PostErrorDelayMs int    `yaml:"postErrorDelayMs"`
	PostCompress     bool   `yaml:"postCompress"`
	Debug            bool   `yaml:"debug"`
}

// InitInput ...
func InitInput(r io.Reader) (*Input, error) {
	in, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("Error reading config yaml from stdin")
	}

	res := &Input{
		DownloadDelayMs:  200,
		AccumBatchSize:   100,
		ParserCount:      runtime.NumCPU(),
		PostErrorDelayMs: 500,
	}
	if err := yaml.Unmarshal(in, &res); err != nil {
		return nil, fmt.Errorf("Error decoding config input: %s", err)
	}

	for _, v := range res.Targets {
		if v.Start < 1 {
			v.Start = 1
		}
		if v.End < 1 {
			v.End = math.MaxInt32
		}
		if v.Limit < 1 {
			v.Limit = 100
		}
		if v.Start > v.End {
			v.Start, v.End = v.End, v.Start
		}
	}

	return res, nil
}
