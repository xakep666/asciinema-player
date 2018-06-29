package asciicast

import (
	"math"
)

// Frames is a list of frame with additional functions.
type Frames []Frame

// Filter returns frame list with frames which satisfies given predicate.
func (f Frames) Filter(p func(frame *Frame) bool) Frames {
	var ret Frames

	for _, frame := range f {
		if p(&frame) {
			ret = append(ret, frame)
		}
	}

	return ret
}

// ToRelativeTime converts frame time to relative values.
func (f Frames) ToRelativeTime() Frames {
	var ret Frames

	prevTime := 0.

	for _, frame := range f {
		ret = append(ret, Frame{
			Time: frame.Time - prevTime,
			Data: append([]byte{}, frame.Data...),
			Type: frame.Type,
		})
		prevTime = frame.Time
	}

	return ret
}

// ToAbsoluteTime converts frame time to absolute values.
func (f Frames) ToAbsoluteTime() Frames {
	var ret Frames

	t := 0.

	for _, frame := range f {
		t += frame.Time
		ret = append(ret, Frame{
			Time: t,
			Type: frame.Type,
			Data: append([]byte{}, frame.Data...),
		})
	}

	return ret
}

// CapRelativeTime sets minimal delay for frames. Frame time must be relative.
func (f Frames) CapRelativeTime(timeLimit float64) Frames {
	var ret Frames

	if timeLimit > 0 {
		for _, frame := range f {
			ret = append(ret, Frame{
				Time: math.Min(frame.Time, timeLimit),
				Type: frame.Type,
				Data: append([]byte{}, frame.Data...),
			})
		}
	} else {
		copy(ret, f)
	}

	return ret
}

// AdjustSpeed allows to change time delay between.
// Speed value must be positive.
// If speed > 1 delay will be increased.
// If speed < 1 delay will be decreased.
func (f Frames) AdjustSpeed(speed float64) Frames {
	var ret Frames

	for _, frame := range f {
		ret = append(ret, Frame{
			Time: frame.Time / speed,
			Type: frame.Type,
			Data: append([]byte{}, frame.Data...),
		})
	}

	return ret
}

// IsOutputFrame determines that given frame is output type.
func IsOutputFrame(frame *Frame) bool {
	return frame.Type == OutputFrame
}
