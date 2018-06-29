package asciicast

import (
	"encoding/json"
	"fmt"
)

type FrameType string

const (
	InputFrame  FrameType = "i"
	OutputFrame FrameType = "o"
)

type Frame struct {
	Time float64
	Type FrameType
	Data []byte
}

// jsonFrame consists from relative timestamp (float64), type ("i" or "o") and terminal text
type jsonFrame [3]interface{}

func (f *Frame) UnmarshalJSON(b []byte) error {
	var rawFrame jsonFrame
	if err := json.Unmarshal(b, &rawFrame); err != nil {
		return err
	}

	switch t := rawFrame[0].(type) {
	case float64:
		f.Time = t
	default:
		return fmt.Errorf("invalid frame[0] type %T", t)
	}

	switch frameTypeRaw := rawFrame[1].(type) {
	case string:
		switch FrameType(frameTypeRaw) {
		case InputFrame, OutputFrame:
			f.Type = FrameType(frameTypeRaw)
		default:
			return fmt.Errorf("invalid frame type %s", frameTypeRaw)
		}
	default:
		return fmt.Errorf("invalid frame[1] type %T", frameTypeRaw)
	}

	switch text := rawFrame[2].(type) {
	case string:
		f.Data = []byte(text)
	default:
		return fmt.Errorf("invalid frame[2] type %T", text)
	}

	return nil
}
