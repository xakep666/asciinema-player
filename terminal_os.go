package player

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

const (
	space = 0x20
	ctrlC = 0x03
)

var ErrNotTerminal = fmt.Errorf("stdin is not terminal")

// OSTerminal represents terminal on operating system.
type OSTerminal struct {
	file *os.File
	stop chan struct{}

	width, height int
	state         *term.State
}

// NewOSTerminal constructs OSTerminal from stdin.
// It returns ErrNotTerminal if stdin is not terminal.
func NewOSTerminal() (*OSTerminal, error) {
	return NewOSTerminalFromFile(os.Stdin)
}

// NewOSTerminalFromFile constructs OSTerminal from file.
// It returns ErrNotTerminal if file is not terminal.
func NewOSTerminalFromFile(file *os.File) (*OSTerminal, error) {
	if !term.IsTerminal(int(file.Fd())) {
		return nil, ErrNotTerminal
	}

	width, height, err := term.GetSize(int(file.Fd()))
	if err != nil {
		return nil, fmt.Errorf("get terminal size failed: %s", err)
	}

	return &OSTerminal{
		file:   file,
		stop:   make(chan struct{}),
		width:  width,
		height: height,
	}, nil
}

func (t *OSTerminal) Write(p []byte) (n int, err error) { return t.file.Write(p) }

// Close closes terminal (stop control loop). It doesn't close underlying file.
func (t *OSTerminal) Close() error {
	close(t.stop)
	return nil
}

func (t *OSTerminal) Dimensions() (width, height int) { return t.width, t.height }

func (t *OSTerminal) ToRaw() error {
	state, err := term.MakeRaw(int(t.file.Fd()))
	if err != nil {
		return fmt.Errorf("%s", err) // decouple from lib error
	}

	// attempt to enable vt100 capabilities to support ansi escape codes
	if err = enableVT100(int(t.file.Fd())); err != nil {
		return fmt.Errorf("%s", err)
	}

	t.state = state

	return nil
}

func (t *OSTerminal) Restore() error {
	if t.state == nil {
		return nil
	}

	// attempt to reset terminal to remove effects possibly set by player
	if _, err := fmt.Fprint(t.file, "\033c"); err != nil {
		return fmt.Errorf("reset terminal failed: %w", err)
	}

	if err := term.Restore(int(t.file.Fd()), t.state); err != nil {
		return fmt.Errorf("%s", err)
	}

	return nil
}

func (t *OSTerminal) Control(control PlaybackControl) {
	var buf [3]byte // for control sequences beginning with "ESC-["
	for {
		select {
		case <-t.stop:
			return
		default:
		}

		n, err := t.file.Read(buf[:])
		if err != nil {
			return
		}

		if n != 1 {
			continue
		}

		switch buf[0] {
		case space:
			control.Pause()
		case ctrlC:
			control.Stop()
		}
	}
}
