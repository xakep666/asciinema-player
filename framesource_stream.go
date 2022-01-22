package player

import (
	"encoding/json"
	"fmt"
	"io"
)

// StreamFrameSource reads frames from io.Reader.
type StreamFrameSource struct {
	dec *json.Decoder

	hdr   Header
	frame Frame
	err   error
}

// NewStreamFrameSource constructs StreamFrameSource. It reads Header from input stream.
func NewStreamFrameSource(reader io.Reader) (*StreamFrameSource, error) {
	dec := json.NewDecoder(reader)

	var hdr Header
	if err := dec.Decode(&hdr); err != nil {
		return nil, fmt.Errorf("read header failed: %w", err)
	}

	return &StreamFrameSource{
		dec: dec,
		hdr: hdr,
	}, nil
}

// Header returns asciinema-v2 header.
func (s *StreamFrameSource) Header() Header { return s.hdr }

// Next advances to next available frame. It must return false if error occurs or there is no more frames.
func (s *StreamFrameSource) Next() bool {
	err := s.dec.Decode(&s.frame)
	switch err {
	case nil:
		return true
	case io.EOF: // all done
		return false
	default:
		s.err = err
		return false
	}
}

// Frame returns current frame. It becomes unusable after Next call.
func (s *StreamFrameSource) Frame() Frame { return s.frame }

// Err returns error if it happens during iteration.
func (s *StreamFrameSource) Err() error { return s.err }
