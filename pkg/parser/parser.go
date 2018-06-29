package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"

	"github.com/xakep666/asciinema-player/pkg/asciicast"
)

func Parse(rd io.Reader) (*asciicast.Asciicast, error) {
	scanner := bufio.NewScanner(rd)
	scanner.Split(bufio.ScanLines)

	var header asciicast.Header
	if scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &header); err != nil {
			return nil, fmt.Errorf("malformed header: %v", err)
		}
	} else {
		return nil, fmt.Errorf("missing header")
	}

	if header.Version != 2 {
		return nil, fmt.Errorf("invalid version %d", header.Version)
	}

	frames := make([]asciicast.Frame, 0)
	for i := 0; scanner.Scan(); i++ {
		var frame asciicast.Frame
		if err := json.Unmarshal(scanner.Bytes(), &frame); err != nil {
			return nil, fmt.Errorf("malformed frame %d: %v", i, err)
		}
		frames = append(frames, frame)
	}

	return &asciicast.Asciicast{
		Header: header,
		Frames: frames,
	}, nil
}
