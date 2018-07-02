package mockterminal

import (
	"bytes"
	"time"

	"github.com/nsf/termbox-go"
	"github.com/xakep666/asciinema-player/pkg/terminal"
)

// MockTerminal is a mock for Terminal interface with event queue and buffer. Used for testing purposes.
type MockTerminal struct {
	eventQueue chan termbox.Event
	buf        *bytes.Buffer
	isRaw      bool
}

// NewMockTerminal constructs a MockTerminal with given queue size.
func NewMockTerminal(queueSize int) *MockTerminal {
	return &MockTerminal{
		eventQueue: make(chan termbox.Event, queueSize),
		buf:        bytes.NewBuffer(make([]byte, 0)),
	}
}

/********************************
*TERMINAL IMPLEMENTING FUNCTIONS*
*********************************/

// Size returns terminal size. Here it is 80x24.
func (m *MockTerminal) Size() (int, int, error) {
	return 80, 24, nil
}

// Write puts data to internal buffer.
func (m *MockTerminal) Write(data []byte) (int, error) {
	return m.buf.Write(data)
}

// ToRaw sets internal "raw" attribute to true.
func (m *MockTerminal) ToRaw() error {
	m.isRaw = true
	return nil
}

// Reset sets internal "raw" attribute to false, resets internal data buffer and clears event queue.
func (m *MockTerminal) Reset() error {
	m.isRaw = false
	m.buf.Reset() // clear buffer data

	// clear events queue
L:
	for {
		select {
		case <-m.eventQueue:
		default:
			break L
		}
	}
	return nil
}

// TimeoutEvent waits event in event queue or returns ErrEventTimeout if no events after timeout.
func (m *MockTerminal) TimeoutEvent(timeout time.Duration) (termbox.Event, error) {
	select {
	case ev := <-m.eventQueue:
		return ev, nil
	case <-time.After(timeout):
		return termbox.Event{}, terminal.ErrEventTimeout
	}
}

// Close calls Reset() actually.
func (m *MockTerminal) Close() error {
	m.Reset()
	return nil
}

/******************************
*EXTERNAL MANAGEMENT FUNCTIONS*
*******************************/

// IsRaw returns "raw" attribute value.
func (m *MockTerminal) IsRaw() bool {
	return m.isRaw
}

// PutEvent puts event to event queue.
func (m *MockTerminal) PutEvent(event termbox.Event) {
	m.eventQueue <- event
}

// RecordedData returns bytes recorded on Write() call.
func (m *MockTerminal) RecordedData() []byte {
	return m.buf.Bytes()
}
