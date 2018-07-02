package asciicast

import (
	"math"
	"testing"
)

var frames = Frames{
	{Time: 0.1, Type: InputFrame, Data: []byte("a")},
	{Time: 0.2, Type: OutputFrame, Data: []byte("b")},
	{Time: 0.3, Type: InputFrame, Data: []byte("c")},
	{Time: 0.4, Type: OutputFrame, Data: []byte("d")},
}

func TestFiltering(t *testing.T) {
	newFrames := frames.Filter(IsOutputFrame)
	for i, frame := range newFrames {
		if frame.Type != OutputFrame {
			t.Errorf("frame %d: unexpected type after filtering %v", i, frame.Type)
		}
	}
}

func TestToRelativeTime(t *testing.T) {
	newFrames := frames.ToRelativeTime()
	for i, frame := range newFrames {
		if math.Abs(frame.Time-0.1) > 0.1 {
			t.Errorf("frame %d: unexpected time %v", i, frame.Time)
		}
	}
}

func TestToAbsoluteTime(t *testing.T) {
	newFrames := frames.ToRelativeTime().ToAbsoluteTime()
	for i, frame := range newFrames {
		if math.Abs(frame.Time-frames[i].Time) > 0.1 {
			t.Errorf("frame %d: unexpected time %v", i, frame.Time)
		}
	}
}

func TestCapRelativeTime(t *testing.T) {
	newFrames := frames.CapRelativeTime(0.3)
	for i, frame := range newFrames {
		if frame.Time > 0.3 {
			t.Errorf("frame %d: unexpected time %v", i, frame.Time)
		}
	}

	// if zero limit given, should not change time
	newFramesUnchanged := frames.CapRelativeTime(0)
	for i, frame := range newFramesUnchanged {
		if frame.Time != frames[i].Time {
			t.Errorf("frame %d: time changed but shouldn`t", i)
		}
	}
}

func TestSpeedAdjustment(t *testing.T) {
	newFrames := frames.AdjustSpeed(2)
	for i, frame := range newFrames {
		if math.Abs(frame.Time-frames[i].Time/2) > 0.01 {
			t.Errorf("frame %d: unexpected time %v", i, frame.Time)
		}
	}
}
