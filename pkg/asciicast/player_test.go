package asciicast

import (
	"testing"
	"time"

	"github.com/nsf/termbox-go"
	"github.com/xakep666/asciinema-player/pkg/terminal/mockterminal"
)

func TestPlay(t *testing.T) {
	testCast := &Asciicast{
		Header: Header{
			Width:         80,
			Height:        24,
			Version:       2,
			IdleTimeLimit: 0.1,
			Timestamp:     uint64(time.Now().Unix()),
		},
		Frames: Frames{
			{Time: 0.1, Type: OutputFrame, Data: []byte("hello")},
		},
	}

	mt := mockterminal.NewMockTerminal(10)
	tp := &TerminalPlayer{
		Terminal: mt,
	}
	err := tp.Play(testCast, time.Millisecond, 1)
	if err != nil {
		t.Errorf("unexpected error %v", err)
		return
	}
	data := mt.RecordedData()
	if string(data) != string(testCast.Frames[0].Data) {
		t.Errorf("unexpected result %s", data)
		return
	}
}

func TestBreak(t *testing.T) {
	testCast := &Asciicast{
		Header: Header{
			Width:         80,
			Height:        24,
			Version:       2,
			IdleTimeLimit: 0.5,
			Timestamp:     uint64(time.Now().Unix()),
		},
		Frames: Frames{
			{Time: 0.5, Type: OutputFrame, Data: []byte("hello")},
			{Time: 1, Type: OutputFrame, Data: []byte("world")},
		},
	}

	mt := mockterminal.NewMockTerminal(10)
	tp := &TerminalPlayer{
		Terminal: mt,
	}

	go func() {
		time.Sleep(501 * time.Millisecond)
		mt.PutEvent(termbox.Event{Type: termbox.EventKey, Key: termbox.KeyCtrlC})
	}()

	err := tp.Play(testCast, 500*time.Millisecond, 1)
	if err != nil {
		t.Errorf("unexpected error %v", err)
		return
	}

	data := mt.RecordedData()
	if string(data) != string(testCast.Frames[0].Data) {
		t.Errorf("unexpected result %s", data)
		return
	}
}

func TestPauseContinue(t *testing.T) {
	testCast := &Asciicast{
		Header: Header{
			Width:         80,
			Height:        24,
			Version:       2,
			IdleTimeLimit: 0.5,
			Timestamp:     uint64(time.Now().Unix()),
		},
		Frames: Frames{
			{Time: 0.5, Type: OutputFrame, Data: []byte("hello")},
			{Time: 1, Type: OutputFrame, Data: []byte("world")},
		},
	}

	mt := mockterminal.NewMockTerminal(10)
	tp := &TerminalPlayer{
		Terminal: mt,
	}

	go func() {
		time.Sleep(501 * time.Millisecond)
		mt.PutEvent(termbox.Event{Type: termbox.EventKey, Key: termbox.KeySpace}) // pause
		data := mt.RecordedData()
		if string(data) != string(testCast.Frames[0].Data) {
			t.Errorf("unexpected result %s", data)
			return
		}
		mt.PutEvent(termbox.Event{Type: termbox.EventKey, Key: termbox.KeySpace}) // continue
	}()

	err := tp.Play(testCast, 500*time.Millisecond, 1)
	if err != nil {
		t.Errorf("unexpected error %v", err)
		return
	}

	if t.Failed() {
		return
	}

	data := mt.RecordedData()
	if string(data) != string(testCast.Frames[0].Data)+string(testCast.Frames[1].Data) {
		t.Errorf("unexpected result %s", data)
		return
	}
}

func TestPauseBreak(t *testing.T) {
	testCast := &Asciicast{
		Header: Header{
			Width:         80,
			Height:        24,
			Version:       2,
			IdleTimeLimit: 0.5,
			Timestamp:     uint64(time.Now().Unix()),
		},
		Frames: Frames{
			{Time: 0.5, Type: OutputFrame, Data: []byte("hello")},
			{Time: 1, Type: OutputFrame, Data: []byte("world")},
		},
	}

	mt := mockterminal.NewMockTerminal(10)
	tp := &TerminalPlayer{
		Terminal: mt,
	}

	go func() {
		time.Sleep(501 * time.Millisecond)
		mt.PutEvent(termbox.Event{Type: termbox.EventKey, Key: termbox.KeySpace}) // pause
		data := mt.RecordedData()
		if string(data) != string(testCast.Frames[0].Data) {
			t.Errorf("unexpected result %s", data)
			return
		}
		mt.PutEvent(termbox.Event{Type: termbox.EventKey, Key: termbox.KeyCtrlC}) // break
	}()

	err := tp.Play(testCast, 500*time.Millisecond, 1)
	if err != nil {
		t.Errorf("unexpected error %v", err)
		return
	}

	data := mt.RecordedData()
	if string(data) != string(testCast.Frames[0].Data) {
		t.Errorf("unexpected result %s", data)
		return
	}
}

func TestPauseNextFrame(t *testing.T) {
	testCast := &Asciicast{
		Header: Header{
			Width:         80,
			Height:        24,
			Version:       2,
			IdleTimeLimit: 0.5,
			Timestamp:     uint64(time.Now().Unix()),
		},
		Frames: Frames{
			{Time: 0.5, Type: OutputFrame, Data: []byte("hello")},
			{Time: 1, Type: OutputFrame, Data: []byte("world")},
			{Time: 1.5, Type: OutputFrame, Data: []byte("nothing")},
		},
	}

	mt := mockterminal.NewMockTerminal(10)
	tp := &TerminalPlayer{
		Terminal: mt,
	}

	go func() {
		time.Sleep(501 * time.Millisecond)
		mt.PutEvent(termbox.Event{Type: termbox.EventKey, Key: termbox.KeySpace}) // pause
		data := mt.RecordedData()
		if string(data) != string(testCast.Frames[0].Data) {
			t.Errorf("unexpected result %s", data)
			return
		}
		mt.PutEvent(termbox.Event{Type: termbox.EventKey, Key: termbox.KeyTab}) // next frame
		time.Sleep(time.Millisecond)
		mt.PutEvent(termbox.Event{Type: termbox.EventKey, Key: termbox.KeyCtrlC}) // break
	}()

	err := tp.Play(testCast, 500*time.Millisecond, 1)
	if err != nil {
		t.Errorf("unexpected error %v", err)
		return
	}

	data := mt.RecordedData()
	if string(data) != string(testCast.Frames[0].Data)+string(testCast.Frames[1].Data) {
		t.Errorf("unexpected result %s", data)
		return
	}
}
