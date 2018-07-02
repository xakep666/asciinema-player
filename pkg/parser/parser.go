package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"

	"github.com/xakep666/asciinema-player/pkg/asciicast"
)

// ParseError represents parsing error
type ParseError struct {
	Text          string
	Line          int
	UnderlyingErr error
}

func (e *ParseError) Error() string {
	if e.UnderlyingErr != nil {
		return fmt.Sprintf("%s [%d]: %v", e.Text, e.Line, e.UnderlyingErr)
	} else {
		return fmt.Sprintf("%s [%d]", e.Text, e.Line)
	}
}

// Parse parses asciinema v2 file. It returns error if version is not "2".
// Parsing is two-stage: first stage is header parsing, second stage is frames parsing.
func Parse(rd io.Reader) (*asciicast.Asciicast, error) {
	scanner := bufio.NewScanner(rd)
	scanner.Split(bufio.ScanLines)

	var header asciicast.Header
	if scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &header); err != nil {
			return nil, &ParseError{Text: "malformed header", Line: 0, UnderlyingErr: err}
		}
	} else {
		return nil, &ParseError{Text: "missing header", Line: 0}
	}

	if header.Version != 2 {
		return nil, &ParseError{Text: fmt.Sprintf("invalid version %d", header.Version), Line: 0}
	}

	frames := make([]asciicast.Frame, 0)
	for i := 0; scanner.Scan(); i++ {
		var frame asciicast.Frame
		if err := json.Unmarshal(scanner.Bytes(), &frame); err != nil {
			return nil, &ParseError{Text: "malformed frame", Line: i + 1, UnderlyingErr: err}
		}
		frames = append(frames, frame)
	}

	return &asciicast.Asciicast{
		Header: header,
		Frames: frames,
	}, nil
}
