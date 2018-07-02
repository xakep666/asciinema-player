package asciicast

import (
	"encoding/json"
	"fmt"
	"time"
)

// FrameType is a type of frame.
// Currently it is only input ("i") frames and output ("o") frames.
type FrameType string

const (
	InputFrame  FrameType = "i"
	OutputFrame FrameType = "o"
)

// Frame represents asciicast v2 frame.
// This is JSON-array with fixed size of 3 elements:
// [0]:	frame delay in seconds (float64),
// [1]:	frame type,
// [2]: frame data (escaped string).
type Frame struct {
	Time float64
	Type FrameType
	Data []byte
}

type jsonFrame [3]interface{}

// FrameUnmarshalError returned if frame conversion to struct failed
type FrameUnmarshalError struct {
	Description string
	Index       int
}

func (e *FrameUnmarshalError) Error() string {
	return fmt.Sprintf("frame[%d]: %s", e.Index, e.Description)
}

// UnmarshalJSON implements json.Unmarshaler
func (f *Frame) UnmarshalJSON(b []byte) error {
	var rawFrame jsonFrame
	if err := json.Unmarshal(b, &rawFrame); err != nil {
		return err
	}

	switch t := rawFrame[0].(type) {
	case float64:
		f.Time = t
	default:
		return &FrameUnmarshalError{Description: fmt.Sprintf("invalid type %T", t), Index: 0}
	}

	switch frameTypeRaw := rawFrame[1].(type) {
	case string:
		switch FrameType(frameTypeRaw) {
		case InputFrame, OutputFrame:
			f.Type = FrameType(frameTypeRaw)
		default:
			return &FrameUnmarshalError{Description: fmt.Sprintf("invalid value %v", frameTypeRaw), Index: 1}
		}
	default:
		return &FrameUnmarshalError{Description: fmt.Sprintf("invalid type %T", frameTypeRaw), Index: 1}
	}

	switch text := rawFrame[2].(type) {
	case string:
		f.Data = []byte(text)
	default:
		return &FrameUnmarshalError{Description: fmt.Sprintf("invalid type %T", text), Index: 2}
	}

	return nil
}

// Duration returns frame duration converted to time.Duration
func (f *Frame) Duration() time.Duration {
	return time.Duration(float64(time.Second) * f.Time)
}
