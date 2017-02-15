package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"runtime"
)

// Input ...
type Input struct {
	Target           string
	MiruURL          string
	OutputDir        string
	DownloadDelayMs  int
	AccumBatchSize   int
	ParserCount      int
	PostErrorDelayMs int
	PostCompress     bool
	Debug            bool
}

func (i *Input) String() string {
	return fmt.Sprintf("target: %s miru: %s dir: %s delay: %d batch: %d count: %d compress: %t debug: %t",
		i.Target, i.MiruURL, i.OutputDir, i.DownloadDelayMs, i.AccumBatchSize, i.ParserCount, i.PostCompress, i.Debug)
}

// InitInput ...
func InitInput(r io.Reader) (*Input, error) {
	res := &Input{
		DownloadDelayMs:  200,
		AccumBatchSize:   100,
		ParserCount:      runtime.NumCPU(),
		PostErrorDelayMs: 500,
		Debug:            false,
	}
	if err := json.NewDecoder(r).Decode(&res); err != nil {
		return nil, fmt.Errorf("Error decoding config input from stdin: %s", err)
	}

	if err := res.IsValid(); err != nil {
		return nil, fmt.Errorf("Input is not valid: %s", err)
	}

	return res, nil
}

// IsValid ...
func (i *Input) IsValid() error {
	if i.Target == "" {
		return fmt.Errorf("Target is empty")
	}
	if i.ParserCount < 1 {
		return fmt.Errorf("Parser count is invalid %d", i.ParserCount)
	}

	return nil
}
