package player

import (
	"encoding/json"
	"fmt"
)

// Frame represents asciinema-v2 frame.
// This is JSON-array with fixed size of 3 elements:
// [0]:	frame delay in seconds (float64),
// [1]:	frame type,
// [2]: frame data (escaped string).
type Frame struct {
	// Time in seconds since record start.
	Time float64

	// Type of frame.
	Type FrameType

	// Data contains frame data.
	Data []byte
}

// FrameUnmarshalError returned if frame conversion to struct failed.
type FrameUnmarshalError struct {
	Description string
	Index       int
}

func (e *FrameUnmarshalError) Error() string {
	return fmt.Sprintf("frame[%d]: %s", e.Index, e.Description)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *Frame) UnmarshalJSON(b []byte) error {
	var rawFrame [3]interface{}
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
