package asciicast

import (
	"math"
)

type Frames []Frame

func (f Frames) Filter(p func(frame *Frame) bool) Frames {
	var ret Frames

	for _, frame := range f {
		if p(&frame) {
			ret = append(ret, frame)
		}
	}

	return ret
}

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

func IsOutputFrame(frame *Frame) bool {
	return frame.Type == OutputFrame
}
