package asciicast

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

const (
	invalidTimeTypeFrame  = `["a", "i", "c"]`
	invalidKindTypeFrame  = `[0.1, 0, "c"]`
	invalidKindValueFrame = `[0.1, "v", "c"]`
	invalidDataTypeFrame  = `[0.1, "i", 0]`
	goodFrame             = `[0.1, "i", "c"]`
	badSyntaxFrame        = `[0,1, i, "c"]`
)

func TestFrameUnmarshal(t *testing.T) {
	testCases := []struct {
		TestData      string
		ExpectedError error
	}{
		{TestData: invalidTimeTypeFrame, ExpectedError: &FrameUnmarshalError{Description: "invalid type string", Index: 0}},
		{TestData: invalidKindTypeFrame, ExpectedError: &FrameUnmarshalError{Description: "invalid type float64", Index: 1}},
		{TestData: invalidKindValueFrame, ExpectedError: &FrameUnmarshalError{Description: "invalid value v", Index: 1}},
		{TestData: invalidDataTypeFrame, ExpectedError: &FrameUnmarshalError{Description: "invalid type float64", Index: 2}},
		{TestData: goodFrame, ExpectedError: nil},
	}

	for i, tc := range testCases {
		var frame Frame
		err := json.Unmarshal([]byte(tc.TestData), &frame)
		if tc.ExpectedError == nil {
			if err != nil {
				t.Errorf("case %d: unexpected error: %v", i, err)
			}
			continue
		}
		if !reflect.DeepEqual(err, tc.ExpectedError) {
			t.Errorf("case %d: unexpected error: %#v", i, err)
			continue
		}
	}

	if err := json.Unmarshal([]byte(badSyntaxFrame), &Frame{}); err == nil {
		t.Errorf("case %d: unexpected nil error", len(testCases))
	}
}

func TestFrameDuration(t *testing.T) {
	frame := Frame{
		Time: 10,
		Type: OutputFrame,
		Data: []byte("hello"),
	}
	duration := frame.Duration()
	if duration != time.Duration(float64(time.Second)*frame.Time) {
		t.Errorf("unexpected frame duration: %v", duration)
	}
}

func TestError(t *testing.T) {
	err := &FrameUnmarshalError{Description: "hello", Index: 0}
	if err.Error() != "frame[0]: hello" {
		t.Errorf("unexpected error text: %v", err.Error())
	}
}
