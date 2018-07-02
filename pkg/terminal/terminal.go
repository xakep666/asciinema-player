package terminal

import (
	"errors"
	"io"
	"os"
	"time"

	"github.com/nsf/termbox-go"
	"golang.org/x/crypto/ssh/terminal"
)

// Terminal is an interface for terminal.
type Terminal interface {
	// Size returns terminal width and height.
	Size() (width, height int, err error)

	// Write puts provided bytes to terminal.
	Write([]byte) (int, error)

	// ToRaw puts terminal to "raw" state. Previous state should be saved.
	ToRaw() error

	// Reset resets terminal to saved state.
	Reset() error

	// TimeoutEvent waits for terminal event. It blocks until event caught or timeout exceeded.
	TimeoutEvent(timeout time.Duration) (termbox.Event, error)

	io.Closer
}

// ErrEventTimeout is error returned by TimeoutEvent if timeout exceeded.
var ErrEventTimeout = errors.New("event timeout")

// Pty is a PTY representation.
type Pty struct {
	Stdin     *os.File
	Stdout    *os.File
	prevState *terminal.State
}

// NewPty attaches to current terminal and performs some initializations.
func NewPty() (*Pty, error) {
	if err := termbox.Init(); err != nil {
		return nil, err
	}
	return &Pty{Stdin: os.Stdin, Stdout: os.Stdout}, nil
}

// Size returns terminal size
func (p *Pty) Size() (int, int, error) {
	w, h := termbox.Size()
	return w, h, nil
}

// Write puts data to terminal
func (p *Pty) Write(data []byte) (int, error) {
	n, err := p.Stdout.Write(data)
	if err != nil {
		return 0, err
	}

	// sync on stdout returns "sync error" which we can`t properly handle
	p.Stdout.Sync()
	return n, nil
}

// ToRaw saves terminal state and tries to put it to raw mode.
func (p *Pty) ToRaw() error {
	fd := p.Stdin.Fd()
	var err error
	if terminal.IsTerminal(int(fd)) {
		p.prevState, err = terminal.MakeRaw(int(fd))
		if err != nil {
			return err
		}
	}
	return nil
}

// Reset restores terminal to saved state.
func (p *Pty) Reset() error {
	if p.prevState == nil {
		return nil
	}
	return terminal.Restore(int(p.Stdin.Fd()), p.prevState)
}

// TimeoutEvent polls terminal event but not greater than provided timeout.
// If timeout exceeded ErrEventTimeout will be returned.
func (p *Pty) TimeoutEvent(timeout time.Duration) (termbox.Event, error) {
	ev := make(chan termbox.Event)

	go func() {
		ev <- termbox.PollEvent()
	}()

	after := time.After(timeout)

	select {
	case <-after:
		return termbox.Event{}, ErrEventTimeout
	case event := <-ev:
		return event, nil
	}

}

// Close resets terminal and performs de-initializations.
func (p *Pty) Close() error {
	p.Reset()
	termbox.Close()
	return nil
}
