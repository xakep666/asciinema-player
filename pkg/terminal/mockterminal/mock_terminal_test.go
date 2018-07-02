package mockterminal

import (
	"reflect"
	"testing"
	"time"

	"github.com/nsf/termbox-go"
	"github.com/xakep666/asciinema-player/pkg/terminal"
)

func TestSize(t *testing.T) {
	mt := NewMockTerminal(1)
	width, height, err := mt.Size()
	if err != nil {
		t.Errorf("unexpected error %v", err)
		return
	}
	if width != 80 || height != 24 {
		t.Errorf("unexpected size: width=%d, height=%d", width, height)
		return
	}
}

func TestRawMode(t *testing.T) {
	mt := NewMockTerminal(1)
	if err := mt.ToRaw(); err != nil {
		t.Errorf("unexpected error %v", err)
		return
	}
	if !mt.IsRaw() {
		t.Errorf("expected raw mode")
		return
	}
	if err := mt.Reset(); err != nil {
		t.Errorf("unexpected error %v", err)
		return
	}
	if mt.IsRaw() {
		t.Errorf("expected non-raw mode")
		return
	}
}

func TestEvents(t *testing.T) {
	mt := NewMockTerminal(2)
	testEvent := termbox.Event{Type: termbox.EventKey, Key: termbox.KeyCtrlC}
	mt.PutEvent(testEvent)
	mt.PutEvent(termbox.Event{Type: termbox.EventKey, Key: termbox.KeyTab})
	event, err := mt.TimeoutEvent(time.Second)
	if err != nil {
		t.Errorf("unexpected error %v", err)
		return
	}
	if !reflect.DeepEqual(testEvent, event) {
		t.Errorf("got unexpected event %#v", event)
		return
	}
	if err := mt.Reset(); err != nil {
		t.Errorf("unexpected error %v", err)
		return
	}
	_, err = mt.TimeoutEvent(time.Second)
	if err != terminal.ErrEventTimeout {
		t.Errorf("unexpected error %v", err)
		return
	}
}

func TestData(t *testing.T) {
	mt := NewMockTerminal(1)
	defer mt.Close()

	testData := "hello\nworld\n"
	n, err := mt.Write([]byte(testData))
	if err != nil {
		t.Errorf("unexpected error %v", err)
		return
	}
	if n != len(testData) {
		t.Errorf("wanted to write %d bytes, wrote %d bytes", len(testData), n)
		return
	}
	recorded := mt.RecordedData()
	if string(recorded) != testData {
		t.Errorf("unexpected recorded data %v", recorded)
		return
	}
}
