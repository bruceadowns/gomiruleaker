package lib

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Targets ...
type Targets struct {
	Prefix string
	Start  int
	End    int
}

func (t *Targets) String() string {
	return fmt.Sprintf("p: %s s: %d e: %d", t.Prefix, t.Start, t.End)
}

// Equals ...
func (t *Targets) Equals(that *Targets) bool {
	if !strings.EqualFold(t.Prefix, that.Prefix) {
		return false
	}
	if t.Start != that.Start {
		return false
	}
	if t.End != that.End {
		return false
	}

	return true
}

// ExpandTargetURL expand a url path that ends will sequence via slice syntax
func ExpandTargetURL(u string) (*Targets, error) {
	seqStart := strings.Index(u, "[")
	seqEnd := strings.LastIndex(u, "]")

	if (seqStart == -1 && seqEnd != -1) || (seqStart != -1 && seqEnd == -1) {
		return nil, fmt.Errorf("Mismatched brackets")
	}

	if seqStart == -1 && seqEnd == -1 {
		lastSlash := strings.LastIndex(u, "/")
		if lastSlash == -1 {
			return nil, fmt.Errorf("No slash found")
		}
		prefix := u[:lastSlash+1]
		sStart := u[lastSlash+1:]
		iStart, err := strconv.Atoi(sStart)
		if err != nil {
			return nil, err
		}
		return &Targets{prefix, iStart, iStart}, nil
	}

	prefix := u[:seqStart]
	seq := u[seqStart+1 : seqEnd]

	if seq == ":" {
		return &Targets{prefix, 1, math.MaxInt32}, nil
	}

	colon := strings.Index(seq, ":")
	if colon == -1 {
		iStart, err := strconv.Atoi(seq)
		if err != nil {
			return nil, err
		}

		return &Targets{prefix, iStart, iStart}, nil
	}

	sStart := seq[:colon]
	if sStart == "" {
		sEnd := seq[colon+1:]
		iEnd, err := strconv.Atoi(sEnd)
		if err != nil {
			return nil, err
		}
		if iEnd < 1 {
			return nil, fmt.Errorf("Error expanding url (end is negative): %s", u)
		}

		return &Targets{prefix, 1, iEnd}, nil
	}

	iStart, err := strconv.Atoi(sStart)
	if err != nil {
		return nil, err
	}
	if iStart < 1 {
		return nil, fmt.Errorf("Error expanding url (start is negative): %s", u)
	}

	sEnd := seq[colon+1:]
	if sEnd == "" {
		return &Targets{prefix, iStart, math.MaxInt32}, nil
	}

	iEnd, err := strconv.Atoi(sEnd)
	if err != nil {
		return nil, err
	}

	if iStart > iEnd {
		iStart, iEnd = iEnd, iStart
	}

	return &Targets{prefix, iStart, iEnd}, nil
}
