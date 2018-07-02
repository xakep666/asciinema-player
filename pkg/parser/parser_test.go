package parser

import (
	"errors"
	"strings"
	"testing"

	"github.com/xakep666/asciinema-player/pkg/asciicast"
)

const malformedHeaderFile = `{"blabla"}
[0.1, "o", "a"]`

const missingHeaderFile = ``

const invalidVersionHeaderFile = `{"version": 4}
[0.1, "o", "a"]`

const malformedFrameFile = `{"version": 2}
[0.1, "v", "c"]`

const goodFile = `{"version": 2, "width": 237, "height": 57, "timestamp": 1530264003, "env": {"SHELL": "/usr/bin/zsh", "TERM": "xterm-256color"}}
[0.1, "o", "frame"]`

func TestParsing(t *testing.T) {
	testCases := []struct {
		TestData          string
		ExpectedErrorText string
	}{
		{TestData: malformedHeaderFile, ExpectedErrorText: "malformed header"},
		{TestData: missingHeaderFile, ExpectedErrorText: "missing header"},
		{TestData: invalidVersionHeaderFile, ExpectedErrorText: "invalid version 4"},
		{TestData: malformedFrameFile, ExpectedErrorText: "malformed frame"},
	}

	for i, testCase := range testCases {
		_, err := Parse(strings.NewReader(testCase.TestData))
		if err == nil {
			t.Errorf("case %d: unexpected nil error", i)
			continue
		}
		perr, ok := err.(*ParseError)
		if !ok {
			t.Errorf("case %d: unexpected error %v", i, err)
			continue
		}
		if perr.Text != testCase.ExpectedErrorText {
			t.Errorf("case %d: unexpected error text: %s", i, perr.Text)
		}
	}

	parsed, err := Parse(strings.NewReader(goodFile))
	if err != nil {
		t.Errorf("case %d: unexpected error %v", len(testCases), err)
		return
	}
	if len(parsed.Frames) == 0 {
		t.Errorf("case %d: no frames parsed", len(testCases))
		return
	}
	frame := parsed.Frames[0]
	if frame.Time != 0.1 || frame.Type != asciicast.OutputFrame || string(frame.Data) != "frame" {
		t.Errorf("case %d: unexpected frame %#v", len(testCases), frame)
		return
	}
}

func TestParseError(t *testing.T) {
	underlyingNilErr := (&ParseError{
		Text:          "text",
		Line:          0,
		UnderlyingErr: nil,
	}).Error()
	if underlyingNilErr != "text [0]" {
		t.Errorf("unexpected error string %v", underlyingNilErr)
	}

	underlyingNonNilErr := (&ParseError{
		Text:          "text",
		Line:          0,
		UnderlyingErr: errors.New("underlying"),
	}).Error()
	if underlyingNonNilErr != "text [0]: underlying" {
		t.Errorf("unexpected error string %v", underlyingNonNilErr)
	}
}
